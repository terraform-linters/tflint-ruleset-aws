package aws

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"golang.org/x/net/idna"
)

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/configs/resource.go#L484-L496
type ProviderConfigRef struct {
	Name       string
	NameRange  hcl.Range
	Alias      string
	AliasRange *hcl.Range // nil if alias not set

	// TODO: this may not be set in some cases, so it is not yet suitable for
	// use outside of this package. We currently only use it for internal
	// validation, but once we verify that this can be set in all cases, we can
	// export this so providers don't need to be re-resolved.
	// This same field is also added to the Provider struct.
	// providerType addrs.Provider
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/configs/resource.go#L498-L569
func DecodeProviderConfigRef(expr hcl.Expression, argName string) (*ProviderConfigRef, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	var shimDiags hcl.Diagnostics
	expr, shimDiags = shimTraversalInString(expr, false)
	diags = append(diags, shimDiags...)

	traversal, travDiags := hcl.AbsTraversalForExpr(expr)

	// AbsTraversalForExpr produces only generic errors, so we'll discard
	// the errors given and produce our own with extra context. If we didn't
	// get any errors then we might still have warnings, though.
	if !travDiags.HasErrors() {
		diags = append(diags, travDiags...)
	}

	if len(traversal) < 1 || len(traversal) > 2 {
		// A provider reference was given as a string literal in the legacy
		// configuration language and there are lots of examples out there
		// showing that usage, so we'll sniff for that situation here and
		// produce a specialized error message for it to help users find
		// the new correct form.
		if exprIsNativeQuotedString(expr) {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider configuration reference",
				Detail:   "A provider configuration reference must not be given in quotes.",
				Subject:  expr.Range().Ptr(),
			})
			return nil, diags
		}

		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider configuration reference",
			Detail:   fmt.Sprintf("The %s argument requires a provider type name, optionally followed by a period and then a configuration alias.", argName),
			Subject:  expr.Range().Ptr(),
		})
		return nil, diags
	}

	// verify that the provider local name is normalized
	name := traversal.RootName()
	nameDiags := checkProviderNameNormalized(name, traversal[0].SourceRange())
	diags = append(diags, nameDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	ret := &ProviderConfigRef{
		Name:      name,
		NameRange: traversal[0].SourceRange(),
	}

	if len(traversal) > 1 {
		aliasStep, ok := traversal[1].(hcl.TraverseAttr)
		if !ok {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid provider configuration reference",
				Detail:   "Provider name must either stand alone or be followed by a period and then a configuration alias.",
				Subject:  traversal[1].SourceRange().Ptr(),
			})
			return ret, diags
		}

		ret.Alias = aliasStep.Name
		ret.AliasRange = aliasStep.SourceRange().Ptr()
	}

	return ret, diags
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/configs/compat_shim.go#L21-L92
// shimTraversalInString takes any arbitrary expression and checks if it is
// a quoted string in the native syntax. If it _is_, then it is parsed as a
// traversal and re-wrapped into a synthetic traversal expression and a
// warning is generated. Otherwise, the given expression is just returned
// verbatim.
//
// This function has no effect on expressions from the JSON syntax, since
// traversals in strings are the required pattern in that syntax.
//
// If wantKeyword is set, the generated warning diagnostic will talk about
// keywords rather than references. The behavior is otherwise unchanged, and
// the caller remains responsible for checking that the result is indeed
// a keyword, e.g. using hcl.ExprAsKeyword.
func shimTraversalInString(expr hcl.Expression, wantKeyword bool) (hcl.Expression, hcl.Diagnostics) {
	// ObjectConsKeyExpr is a special wrapper type used for keys on object
	// constructors to deal with the fact that naked identifiers are normally
	// handled as "bareword" strings rather than as variable references. Since
	// we know we're interpreting as a traversal anyway (and thus it won't
	// matter whether it's a string or an identifier) we can safely just unwrap
	// here and then process whatever we find inside as normal.
	if ocke, ok := expr.(*hclsyntax.ObjectConsKeyExpr); ok {
		expr = ocke.Wrapped
	}

	if !exprIsNativeQuotedString(expr) {
		return expr, nil
	}

	strVal, diags := expr.Value(nil)
	if diags.HasErrors() || strVal.IsNull() || !strVal.IsKnown() {
		// Since we're not even able to attempt a shim here, we'll discard
		// the diagnostics we saw so far and let the caller's own error
		// handling take care of reporting the invalid expression.
		return expr, nil
	}

	// The position handling here isn't _quite_ right because it won't
	// take into account any escape sequences in the literal string, but
	// it should be close enough for any error reporting to make sense.
	srcRange := expr.Range()
	startPos := srcRange.Start // copy
	startPos.Column++          // skip initial quote
	startPos.Byte++            // skip initial quote

	traversal, tDiags := hclsyntax.ParseTraversalAbs(
		[]byte(strVal.AsString()),
		srcRange.Filename,
		startPos,
	)
	diags = append(diags, tDiags...)

	if wantKeyword {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Quoted keywords are deprecated",
			Detail:   "In this context, keywords are expected literally rather than in quotes. Terraform 0.11 and earlier required quotes, but quoted keywords are now deprecated and will be removed in a future version of Terraform. Remove the quotes surrounding this keyword to silence this warning.",
			Subject:  &srcRange,
		})
	} else {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Quoted references are deprecated",
			Detail:   "In this context, references are expected literally rather than in quotes. Terraform 0.11 and earlier required quotes, but quoted references are now deprecated and will be removed in a future version of Terraform. Remove the quotes surrounding this reference to silence this warning.",
			Subject:  &srcRange,
		})
	}

	return &hclsyntax.ScopeTraversalExpr{
		Traversal: traversal,
		SrcRange:  srcRange,
	}, diags
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/configs/util.go#L8-L18
// exprIsNativeQuotedString determines whether the given expression looks like
// it's a quoted string in the HCL native syntax.
//
// This should be used sparingly only for situations where our legacy HCL
// decoding would've expected a keyword or reference in quotes but our new
// decoding expects the keyword or reference to be provided directly as
// an identifier-based expression.
func exprIsNativeQuotedString(expr hcl.Expression) bool {
	_, ok := expr.(*hclsyntax.TemplateExpr)
	return ok
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/configs/provider.go#L256-L282
// checkProviderNameNormalized verifies that the given string is already
// normalized and returns an error if not.
func checkProviderNameNormalized(name string, declrange hcl.Range) hcl.Diagnostics {
	var diags hcl.Diagnostics
	// verify that the provider local name is normalized
	normalized, err := IsProviderPartNormalized(name)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider local name",
			Detail:   fmt.Sprintf("%s is an invalid provider local name: %s", name, err),
			Subject:  &declrange,
		})
		return diags
	}
	if !normalized {
		// we would have returned this error already
		normalizedProvider, _ := ParseProviderPart(name)
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid provider local name",
			Detail:   fmt.Sprintf("Provider names must be normalized. Replace %q with %q to fix this error.", name, normalizedProvider),
			Subject:  &declrange,
		})
	}
	return diags
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/addrs/provider.go#L454-L464
// IsProviderPartNormalized compares a given string to the result of ParseProviderPart(string)
func IsProviderPartNormalized(str string) (bool, error) {
	normalized, err := ParseProviderPart(str)
	if err != nil {
		return false, err
	}
	if str == normalized {
		return true, nil
	}
	return false, nil
}

// original code: https://github.com/hashicorp/terraform/blob/3fbedf25430ead97eb42575d344427db3c32d524/internal/addrs/provider.go#L385-L442
// ParseProviderPart processes an addrs.Provider namespace or type string
// provided by an end-user, producing a normalized version if possible or
// an error if the string contains invalid characters.
//
// A provider part is processed in the same way as an individual label in a DNS
// domain name: it is transformed to lowercase per the usual DNS case mapping
// and normalization rules and may contain only letters, digits, and dashes.
// Additionally, dashes may not appear at the start or end of the string.
//
// These restrictions are intended to allow these names to appear in fussy
// contexts such as directory/file names on case-insensitive filesystems,
// repository names on GitHub, etc. We're using the DNS rules in particular,
// rather than some similar rules defined locally, because the hostname part
// of an addrs.Provider is already a hostname and it's ideal to use exactly
// the same case folding and normalization rules for all of the parts.
//
// In practice a provider type string conventionally does not contain dashes
// either. Such names are permitted, but providers with such type names will be
// hard to use because their resource type names will not be able to contain
// the provider type name and thus each resource will need an explicit provider
// address specified. (A real-world example of such a provider is the
// "google-beta" variant of the GCP provider, which has resource types that
// start with the "google_" prefix instead.)
//
// It's valid to pass the result of this function as the argument to a
// subsequent call, in which case the result will be identical.
func ParseProviderPart(given string) (string, error) {
	if len(given) == 0 {
		return "", fmt.Errorf("must have at least one character")
	}

	// We're going to process the given name using the same "IDNA" library we
	// use for the hostname portion, since it already implements the case
	// folding rules we want.
	//
	// The idna library doesn't expose individual label parsing directly, but
	// once we've verified it doesn't contain any dots we can just treat it
	// like a top-level domain for this library's purposes.
	if strings.ContainsRune(given, '.') {
		return "", fmt.Errorf("dots are not allowed")
	}

	// We don't allow names containing multiple consecutive dashes, just as
	// a matter of preference: they look weird, confusing, or incorrect.
	// This also, as a side-effect, prevents the use of the "punycode"
	// indicator prefix "xn--" that would cause the IDNA library to interpret
	// the given name as punycode, because that would be weird and unexpected.
	if strings.Contains(given, "--") {
		return "", fmt.Errorf("cannot use multiple consecutive dashes")
	}

	result, err := idna.Lookup.ToUnicode(given)
	if err != nil {
		return "", fmt.Errorf("must contain only letters, digits, and dashes, and may not use leading or trailing dashes")
	}

	return result, nil
}
