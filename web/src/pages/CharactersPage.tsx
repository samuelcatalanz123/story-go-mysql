import { ResourcePage } from "../components/ResourcePage";
import {
  listCharacters,
  createCharacter,
  updateCharacter,
  deleteCharacter,
} from "../api/resources";

export function CharactersPage() {
  return (
    <ResourcePage
      heading="Personajes"
      list={listCharacters}
      create={createCharacter}
      update={updateCharacter}
      remove={deleteCharacter}
      detailBase="/characters"
    />
  );
}
