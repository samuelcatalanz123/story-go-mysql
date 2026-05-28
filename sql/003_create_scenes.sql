CREATE TABLE scenes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL UNIQUE,
  text TEXT NULL,
  start_timeline INT NOT NULL,
  end_timeline INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE scene_characters (
  scene_id BIGINT UNSIGNED NOT NULL,
  character_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (scene_id, character_id),
  CONSTRAINT fk_scene_characters_scene
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE,
  CONSTRAINT fk_scene_characters_character
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE TABLE scene_locations (
  scene_id BIGINT UNSIGNED NOT NULL,
  location_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (scene_id, location_id),
  CONSTRAINT fk_scene_locations_scene
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE,
  CONSTRAINT fk_scene_locations_location
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);