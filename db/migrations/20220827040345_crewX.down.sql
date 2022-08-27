ALTER TABLE "scene" DROP CONSTRAINT "scene_story_id";

ALTER TABLE "story" DROP CONSTRAINT "story_customer_id";

ALTER TABLE "scene" DROP CONSTRAINT "scene_background_audio_id";

ALTER TABLE "scene" DROP CONSTRAINT "scene_generated_audio_id";

DROP TABLE "story";

DROP TABLE "scene";

DROP TABLE "background_audio";

DROP TABLE "generated_audio";

DROP TABLE "image";

DROP TABLE "customer";