import { ResourcePage } from "../components/ResourcePage";
import {
  listConflicts,
  createConflict,
  updateConflict,
  deleteConflict,
} from "../api/resources";

export function ConflictsPage() {
  return (
    <ResourcePage
      heading="Conflictos"
      list={listConflicts}
      create={createConflict}
      update={updateConflict}
      remove={deleteConflict}
    />
  );
}
