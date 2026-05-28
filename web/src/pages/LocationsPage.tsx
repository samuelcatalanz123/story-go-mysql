import { ResourcePage } from "../components/ResourcePage";
import {
  listLocations,
  createLocation,
  updateLocation,
  deleteLocation,
} from "../api/resources";

export function LocationsPage() {
  return (
    <ResourcePage
      heading="Lugares"
      list={listLocations}
      create={createLocation}
      update={updateLocation}
      remove={deleteLocation}
    />
  );
}
