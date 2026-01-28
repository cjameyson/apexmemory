--
-- PostgreSQL database dump
--

-- Dumped from database version 18beta2 (Debian 18~beta2-1.pgdg120+1)
-- Dumped by pg_dump version 18beta2 (Debian 18~beta2-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: app; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA app;


--
-- Name: auth_provider; Type: TYPE; Schema: app; Owner: -
--

CREATE TYPE app.auth_provider AS ENUM (
    'password',
    'google',
    'apple'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auth_identities; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.auth_identities (
    id uuid DEFAULT uuidv7() NOT NULL,
    user_id uuid NOT NULL,
    provider app.auth_provider NOT NULL,
    provider_user_id text NOT NULL,
    email app.citext,
    email_verified_at timestamp with time zone,
    password_hash text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT auth_identities_check CHECK (((provider <> 'password'::app.auth_provider) OR ((password_hash IS NOT NULL) AND (email IS NOT NULL))))
);


--
-- Name: notebooks; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.notebooks (
    user_id uuid NOT NULL,
    id uuid DEFAULT uuidv7() NOT NULL,
    name text NOT NULL,
    description text,
    emoji text,
    color text,
    fsrs_settings jsonb DEFAULT '{"params": [0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001, 1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014, 1.8729, 0.5425, 0.0912, 0.0658, 0.1542], "version": "6", "enable_fuzzing": true, "learning_steps": [60, 600], "maximum_interval": 36500, "relearning_steps": [600], "desired_retention": 0.9}'::jsonb NOT NULL,
    "position" integer DEFAULT 0 NOT NULL,
    archived_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT notebooks_name_length CHECK ((length(name) <= 255)),
    CONSTRAINT notebooks_name_not_empty CHECK ((length(TRIM(BOTH FROM name)) > 0))
);


--
-- Name: schema_versions; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.schema_versions (
    version integer NOT NULL
);


--
-- Name: user_sessions; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.user_sessions (
    id uuid DEFAULT uuidv7() NOT NULL,
    user_id uuid NOT NULL,
    token_hash bytea NOT NULL,
    user_agent text,
    ip_address inet,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_used_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.users (
    id uuid DEFAULT uuidv7() NOT NULL,
    email app.citext NOT NULL,
    email_verified_at timestamp with time zone,
    username app.citext NOT NULL,
    display_name text,
    avatar_url text,
    locale text DEFAULT 'en-US'::text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: auth_identities auth_identities_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.auth_identities
    ADD CONSTRAINT auth_identities_pkey PRIMARY KEY (id);


--
-- Name: auth_identities auth_identities_provider_provider_user_id_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.auth_identities
    ADD CONSTRAINT auth_identities_provider_provider_user_id_key UNIQUE (provider, provider_user_id);


--
-- Name: auth_identities auth_identities_user_id_provider_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.auth_identities
    ADD CONSTRAINT auth_identities_user_id_provider_key UNIQUE (user_id, provider);


--
-- Name: notebooks notebooks_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notebooks
    ADD CONSTRAINT notebooks_pkey PRIMARY KEY (user_id, id);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_token_hash_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.user_sessions
    ADD CONSTRAINT user_sessions_token_hash_key UNIQUE (token_hash);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: ix_auth_identities_user; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_auth_identities_user ON app.auth_identities USING btree (user_id);


--
-- Name: ix_notebooks_user_list; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_notebooks_user_list ON app.notebooks USING btree (user_id, archived_at, "position");


--
-- Name: ix_user_sessions_expires_at; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_user_sessions_expires_at ON app.user_sessions USING btree (expires_at);


--
-- Name: ix_user_sessions_user_id; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_user_sessions_user_id ON app.user_sessions USING btree (user_id);


--
-- Name: ux_auth_identities_password_email; Type: INDEX; Schema: app; Owner: -
--

CREATE UNIQUE INDEX ux_auth_identities_password_email ON app.auth_identities USING btree (email) WHERE (provider = 'password'::app.auth_provider);


--
-- Name: auth_identities trg_auth_identities_set_updated_at; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_auth_identities_set_updated_at BEFORE UPDATE ON app.auth_identities FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();


--
-- Name: notebooks trg_notebooks_set_updated_at; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_notebooks_set_updated_at BEFORE UPDATE ON app.notebooks FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();


--
-- Name: users trg_users_set_updated_at; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_users_set_updated_at BEFORE UPDATE ON app.users FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();


--
-- Name: auth_identities auth_identities_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.auth_identities
    ADD CONSTRAINT auth_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- Name: notebooks notebooks_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notebooks
    ADD CONSTRAINT notebooks_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

