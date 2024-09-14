# A data sources to retrieve the existing CNAME records.
data "pihole_cname_records" "records" {}

# When retrieving a CNAME resource via a data source immediately after creating the resource,
# It is recommended to set a depends_on constraint in order to let the resource finish
# creating before performing the fetching of all CNAME records.
data "pihole_cname_records" "records" {
    depends_on = [RESOURCE_IDENTIFIER]
}
