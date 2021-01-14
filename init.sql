CREATE TABLE "flags" ("id" integer,"image" text NOT NULL UNIQUE ,"flag" text NOT NULL UNIQUE , PRIMARY KEY (id));
CREATE TABLE "owned" ("id" integer,"image" varchar NOT NULL UNIQUE,"owned_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id));

