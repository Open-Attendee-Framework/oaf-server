-- +migrate Up
CREATE TABLE "Organizations" (
	"OrganizationID" serial NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	"Picture" bytea,
	CONSTRAINT "Organizations_pk" PRIMARY KEY ("OrganizationID")
);



CREATE TABLE "Users" (
	"UserID" serial NOT NULL,
	"Username" VARCHAR(255) NOT NULL,
	"Password" VARCHAR(255) NOT NULL,
	"Shownname" VARCHAR(255),
	"EMail" VARCHAR(255) NOT NULL,
	"SuperUser" BOOLEAN NOT NULL,
	CONSTRAINT "Users_pk" PRIMARY KEY ("UserID")
);



CREATE TABLE "Members" (
	"SectionID" integer NOT NULL,
	"UserID" integer NOT NULL,
	"Rights" integer NOT NULL,
	CONSTRAINT "Members_pk" PRIMARY KEY ("SectionID", "UserID")
);



CREATE TABLE "Events" (
	"EventID" serial NOT NULL,
	"OrganizationID" integer NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	"Address" VARCHAR(255),
	"Start" TIMESTAMP NOT NULL,
	"End" TIMESTAMP,
	"Creator" integer NOT NULL,
	CONSTRAINT "Events_pk" PRIMARY KEY ("EventID")
);



CREATE TABLE "Attendees" (
	"EventID" integer NOT NULL,
	"UserID" integer NOT NULL,
	"Commitment" integer NOT NULL,
	"Comment" VARCHAR(255),
	CONSTRAINT "Attendees_pk" PRIMARY KEY ("EventID", "UserID")
);



CREATE TABLE "Comments" (
	"CommentID" serial NOT NULL,
	"EventID" integer NOT NULL,
	"UserID" integer NOT NULL,
	"Comment" VARCHAR(255) NOT NULL,
	CONSTRAINT "Comments_pk" PRIMARY KEY ("CommentID")
);


CREATE TABLE "Sections" (
	"SectionID" serial NOT NULL,
	"OrganizationID" integer NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	CONSTRAINT "Sections_pk" PRIMARY KEY ("SectionID")
);


CREATE TABLE "Info" (
	"Key" varchar(255) NOT NULL,
	"Value" varchar(255) NOT NULL,
	CONSTRAINT "Info_pk" PRIMARY KEY ("Key")
);


ALTER TABLE "Members" ADD CONSTRAINT "Members_fk0" FOREIGN KEY ("SectionID") REFERENCES "Sections"("SectionID");
ALTER TABLE "Members" ADD CONSTRAINT "Members_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Events" ADD CONSTRAINT "Events_fk0" FOREIGN KEY ("OrganizationID") REFERENCES "Organizations"("OrganizationID");
ALTER TABLE "Events" ADD CONSTRAINT "Events_fk1" FOREIGN KEY ("Creator") REFERENCES "Users"("UserID");

ALTER TABLE "Attendees" ADD CONSTRAINT "Attendees_fk0" FOREIGN KEY ("EventID") REFERENCES "Events"("EventID");
ALTER TABLE "Attendees" ADD CONSTRAINT "Attendees_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Comments" ADD CONSTRAINT "Comments_fk0" FOREIGN KEY ("EventID") REFERENCES "Events"("EventID");
ALTER TABLE "Comments" ADD CONSTRAINT "Comments_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Sections" ADD CONSTRAINT "Sections_fk0" FOREIGN KEY ("OrganizationID") REFERENCES "Organizations"("OrganizationID");

-- +migrate Down

DROP TABLE "Info";
DROP TABLE "Attendees";
DROP TABLE "Comments";
DROP TABLE "Events";
DROP TABLE "Members";
DROP TABLE "Sections";
DROP TABLE "Organizations";
DROP TABLE "Users";
