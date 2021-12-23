import = "aws-sdk-go/models/apis/appmesh/2019-01-25/api-2.json"

mapping "aws_appmesh_gateway_route" {
  name = ResourceName
  mesh_name = ResourceName
  virtual_gateway_name = ResourceName
  mesh_owner = AccountId
  spec = GatewayRouteSpec
  tags = TagList
}

mapping "aws_appmesh_mesh" {
  name = ResourceName
}

mapping "aws_appmesh_route" {
  name                = ResourceName
  mesh_name           = ResourceName
  virtual_router_name = ResourceName
}

mapping "aws_appmesh_virtual_gateway" {
  name = ResourceName
  mesh_name = ResourceName
  mesh_owner = AccountId
  spec = VirtualGatewaySpec
  tags = TagList
}

mapping "aws_appmesh_virtual_node" {
  name      = ResourceName
  mesh_name = ResourceName
}

mapping "aws_appmesh_virtual_router" {
  name      = ResourceName
  mesh_name = ResourceName
}

mapping "aws_appmesh_virtual_service" {
  name      = ResourceName
  mesh_name = ResourceName
}
