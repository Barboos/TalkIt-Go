--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: microposts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.microposts (
    id uuid NOT NULL,
    title character varying(255) NOT NULL,
    content character varying(255) NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.microposts OWNER TO postgres;

--
-- Name: relationships; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.relationships (
    id uuid NOT NULL,
    follower_id uuid NOT NULL,
    followed_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.relationships OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    name character varying(255) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: microposts microposts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.microposts
    ADD CONSTRAINT microposts_pkey PRIMARY KEY (id);


--
-- Name: relationships relationships_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.relationships
    ADD CONSTRAINT relationships_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: relationships_followed_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX relationships_followed_id_idx ON public.relationships USING btree (followed_id);


--
-- Name: relationships_follower_id_followed_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX relationships_follower_id_followed_id_idx ON public.relationships USING btree (follower_id, followed_id);


--
-- Name: relationships_follower_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX relationships_follower_id_idx ON public.relationships USING btree (follower_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

