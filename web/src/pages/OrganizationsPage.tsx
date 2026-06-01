import { ResourcePage } from "../components/ResourcePage";
import {
  listOrganizations,
  createOrganization,
  updateOrganization,
  deleteOrganization,
} from "../api/resources";

export function OrganizationsPage() {
  return (
    <ResourcePage
      heading="Organizaciones"
      list={listOrganizations}
      create={createOrganization}
      update={updateOrganization}
      remove={deleteOrganization}
    />
  );
}
