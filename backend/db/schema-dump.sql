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


--
-- Name: card_state; Type: TYPE; Schema: app; Owner: -
--

CREATE TYPE app.card_state AS ENUM (
    'new',
    'learning',
    'review',
    'relearning'
);


--
-- Name: rating; Type: TYPE; Schema: app; Owner: -
--

CREATE TYPE app.rating AS ENUM (
    'again',
    'hard',
    'good',
    'easy'
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
-- Name: cards; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.cards (
    user_id uuid NOT NULL,
    id uuid DEFAULT uuidv7() NOT NULL,
    notebook_id uuid NOT NULL,
    fact_id uuid NOT NULL,
    element_id text DEFAULT ''::text NOT NULL,
    state app.card_state DEFAULT 'new'::app.card_state NOT NULL,
    stability real,
    difficulty real,
    step smallint,
    due timestamp with time zone,
    last_review timestamp with time zone,
    elapsed_days real DEFAULT 0 NOT NULL,
    scheduled_days real DEFAULT 0 NOT NULL,
    reps integer DEFAULT 0 NOT NULL,
    lapses integer DEFAULT 0 NOT NULL,
    suspended_at timestamp with time zone,
    buried_until date,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT cards_valid_difficulty CHECK (((difficulty IS NULL) OR ((difficulty >= (1.0)::double precision) AND (difficulty <= (10.0)::double precision)))),
    CONSTRAINT cards_valid_element_id CHECK (((element_id = ''::text) OR (element_id ~ '^c[1-9][0-9]{0,2}$'::text) OR (element_id ~ '^m_[a-zA-Z0-9_-]{6,24}$'::text))),
    CONSTRAINT cards_valid_stability CHECK (((stability IS NULL) OR (stability >= (0)::double precision))),
    CONSTRAINT cards_valid_step_state CHECK ((((state = ANY (ARRAY['learning'::app.card_state, 'relearning'::app.card_state])) AND (step IS NOT NULL) AND (step >= 0)) OR ((state = ANY (ARRAY['new'::app.card_state, 'review'::app.card_state])) AND (step IS NULL))))
);


--
-- Name: facts; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.facts (
    user_id uuid NOT NULL,
    id uuid DEFAULT uuidv7() NOT NULL,
    notebook_id uuid NOT NULL,
    fact_type text DEFAULT 'basic'::text NOT NULL,
    content jsonb NOT NULL,
    source_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT facts_valid_content CHECK (((jsonb_typeof((content -> 'version'::text)) = 'number'::text) AND (jsonb_typeof((content -> 'fields'::text)) = 'array'::text))),
    CONSTRAINT facts_valid_type CHECK ((fact_type = ANY (ARRAY['basic'::text, 'cloze'::text, 'image_occlusion'::text])))
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
    "position" integer DEFAULT 0 NOT NULL,
    total_cards integer DEFAULT 0 NOT NULL,
    fsrs_settings jsonb DEFAULT '{}'::jsonb NOT NULL,
    archived_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT notebooks_name_length CHECK ((length(name) <= 255)),
    CONSTRAINT notebooks_name_not_empty CHECK ((length(TRIM(BOTH FROM name)) > 0)),
    CONSTRAINT notebooks_total_cards_non_negative CHECK ((total_cards >= 0))
);


--
-- Name: notes; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.notes (
    user_id uuid NOT NULL,
    id uuid DEFAULT uuidv7() NOT NULL,
    notebook_id uuid NOT NULL,
    note_type text DEFAULT 'basic'::text NOT NULL,
    content jsonb NOT NULL,
    source_id uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT notes_valid_content CHECK (((jsonb_typeof((content -> 'version'::text)) = 'number'::text) AND (jsonb_typeof((content -> 'fields'::text)) = 'array'::text))),
    CONSTRAINT notes_valid_type CHECK ((note_type = ANY (ARRAY['basic'::text, 'cloze'::text, 'image_occlusion'::text])))
);


--
-- Name: reviews; Type: TABLE; Schema: app; Owner: -
--

CREATE TABLE app.reviews (
    user_id uuid NOT NULL,
    id uuid DEFAULT uuidv7() NOT NULL,
    card_id uuid,
    notebook_id uuid NOT NULL,
    fact_id uuid,
    element_id text,
    reviewed_at timestamp with time zone DEFAULT now() NOT NULL,
    rating app.rating NOT NULL,
    review_duration_ms integer,
    state_before app.card_state NOT NULL,
    stability_before real,
    difficulty_before real,
    elapsed_days real NOT NULL,
    scheduled_days real NOT NULL,
    state_after app.card_state NOT NULL,
    stability_after real NOT NULL,
    difficulty_after real NOT NULL,
    interval_days real NOT NULL,
    retrievability real,
    created_at timestamp with time zone DEFAULT now() NOT NULL
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
-- Name: cards cards_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (user_id, id);


--
-- Name: cards cards_user_id_fact_id_element_id_key; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.cards
    ADD CONSTRAINT cards_user_id_fact_id_element_id_key UNIQUE (user_id, fact_id, element_id);


--
-- Name: facts facts_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.facts
    ADD CONSTRAINT facts_pkey PRIMARY KEY (user_id, id);


--
-- Name: notebooks notebooks_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notebooks
    ADD CONSTRAINT notebooks_pkey PRIMARY KEY (user_id, id);


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notes
    ADD CONSTRAINT notes_pkey PRIMARY KEY (user_id, id);


--
-- Name: reviews reviews_pkey; Type: CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (user_id, id);


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
-- Name: ix_cards_due; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_cards_due ON app.cards USING btree (user_id, notebook_id, due) WHERE ((suspended_at IS NULL) AND (buried_until IS NULL));


--
-- Name: ix_cards_fact; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_cards_fact ON app.cards USING btree (user_id, fact_id);


--
-- Name: ix_cards_state; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_cards_state ON app.cards USING btree (user_id, notebook_id, state);


--
-- Name: ix_facts_notebook; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_facts_notebook ON app.facts USING btree (user_id, notebook_id, created_at DESC);


--
-- Name: ix_notebooks_user_list; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_notebooks_user_list ON app.notebooks USING btree (user_id, archived_at, "position");


--
-- Name: ix_notes_notebook; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_notes_notebook ON app.notes USING btree (user_id, notebook_id, created_at DESC);


--
-- Name: ix_reviews_card; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_reviews_card ON app.reviews USING btree (user_id, card_id, reviewed_at DESC) WHERE (card_id IS NOT NULL);


--
-- Name: ix_reviews_notebook_time; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_reviews_notebook_time ON app.reviews USING btree (user_id, notebook_id, reviewed_at DESC);


--
-- Name: ix_reviews_optimizer; Type: INDEX; Schema: app; Owner: -
--

CREATE INDEX ix_reviews_optimizer ON app.reviews USING btree (user_id, notebook_id, state_before, rating) WHERE (scheduled_days >= (1)::double precision);


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
-- Name: cards trg_cards_set_updated_at; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_cards_set_updated_at BEFORE UPDATE ON app.cards FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();


--
-- Name: cards trg_cards_sync_notebook_count; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_cards_sync_notebook_count AFTER INSERT OR DELETE OR UPDATE OF notebook_id ON app.cards FOR EACH ROW EXECUTE FUNCTION app_code.tg_sync_notebook_card_count();


--
-- Name: facts trg_facts_set_updated_at; Type: TRIGGER; Schema: app; Owner: -
--

CREATE TRIGGER trg_facts_set_updated_at BEFORE UPDATE ON app.facts FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();


--
-- Name: auth_identities auth_identities_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.auth_identities
    ADD CONSTRAINT auth_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- Name: cards cards_user_id_fact_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.cards
    ADD CONSTRAINT cards_user_id_fact_id_fkey FOREIGN KEY (user_id, fact_id) REFERENCES app.facts(user_id, id) ON DELETE CASCADE;


--
-- Name: cards cards_user_id_notebook_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.cards
    ADD CONSTRAINT cards_user_id_notebook_id_fkey FOREIGN KEY (user_id, notebook_id) REFERENCES app.notebooks(user_id, id) ON DELETE CASCADE;


--
-- Name: facts facts_user_id_notebook_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.facts
    ADD CONSTRAINT facts_user_id_notebook_id_fkey FOREIGN KEY (user_id, notebook_id) REFERENCES app.notebooks(user_id, id) ON DELETE CASCADE;


--
-- Name: notebooks notebooks_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notebooks
    ADD CONSTRAINT notebooks_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- Name: notes notes_user_id_notebook_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.notes
    ADD CONSTRAINT notes_user_id_notebook_id_fkey FOREIGN KEY (user_id, notebook_id) REFERENCES app.notebooks(user_id, id) ON DELETE CASCADE;


--
-- Name: reviews reviews_user_id_card_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.reviews
    ADD CONSTRAINT reviews_user_id_card_id_fkey FOREIGN KEY (user_id, card_id) REFERENCES app.cards(user_id, id) ON DELETE SET NULL;


--
-- Name: reviews reviews_user_id_fact_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.reviews
    ADD CONSTRAINT reviews_user_id_fact_id_fkey FOREIGN KEY (user_id, fact_id) REFERENCES app.facts(user_id, id) ON DELETE SET NULL;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: app; Owner: -
--

ALTER TABLE ONLY app.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

