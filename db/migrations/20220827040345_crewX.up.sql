CREATE TABLE "story" (
    "id" varchar(255) PRIMARY KEY,
    "scene_id" varchar(255),
    "name" varchar(255),
    "description" varchar(255),
    "customer_id" varchar(255),
    "status" varchar(255),
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "scene" (
    "id" varchar(255) PRIMARY KEY,
    "generated_audio_id" varchar(255),
    "background_audio_id" varchar(255),
    "status" varchar(255),
    "scene_number" int8,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "background_audio" (
    "id" varchar(255) PRIMARY KEY,
    "path" varchar(255),
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "generated_audio" (
    "id" varchar(255) PRIMARY KEY,
    "path" varchar(255),
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "image" (
    "id" varchar(255) PRIMARY KEY,
    "path" varchar(255),
    "encoded_image" varchar(255),
    "scene_id" varchar(255),
    "is_deleted" boolean,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "customer" (
  "id" varchar(255) PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255),
  "mobile" varchar(255) NOT NULL,
    "created_at" timestamp,
  "updated_at" timestamp
);

ALTER TABLE "story" ADD CONSTRAINT "story_scene_id" FOREIGN KEY ("scene_id") REFERENCES "scene" ("id");

ALTER TABLE "story" ADD CONSTRAINT "story_customer_id" FOREIGN KEY ("customer_id") REFERENCES "customer" ("id");

ALTER TABLE "scene" ADD CONSTRAINT "scene_background_audio_id" FOREIGN KEY ("background_audio_id") REFERENCES "background_audio" ("id");

ALTER TABLE "scene" ADD CONSTRAINT "scene_generated_audio_id" FOREIGN KEY ("generated_audio_id") REFERENCES "generated_audio" ("id");
