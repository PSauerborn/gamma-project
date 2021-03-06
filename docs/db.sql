--
-- PostgreSQL database dump
--

-- Dumped from database version 12.4 (Debian 12.4-1.pgdg100+1)
-- Dumped by pg_dump version 12.4 (Debian 12.4-1.pgdg100+1)

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

SET default_table_access_method = heap;

--
-- Name: assigned_jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.assigned_jobs (
    id uuid NOT NULL,
    uid text NOT NULL
);


ALTER TABLE public.assigned_jobs OWNER TO postgres;

--
-- Name: file_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.file_metadata (
    file_id uuid NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    size integer NOT NULL,
    metadata json DEFAULT '{}'::json NOT NULL,
    archived boolean DEFAULT false NOT NULL,
    file_name text NOT NULL
);


ALTER TABLE public.file_metadata OWNER TO postgres;

--
-- Name: jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.jobs (
    id uuid NOT NULL,
    name text NOT NULL,
    due timestamp without time zone,
    meta json DEFAULT '{}'::json NOT NULL,
    state integer NOT NULL,
    created timestamp without time zone NOT NULL,
    assigned boolean DEFAULT false NOT NULL
);


ALTER TABLE public.jobs OWNER TO postgres;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_roles (
    uid text NOT NULL,
    role integer NOT NULL
);


ALTER TABLE public.user_roles OWNER TO postgres;

--
-- Name: assigned_jobs assigned_jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.assigned_jobs
    ADD CONSTRAINT assigned_jobs_pkey PRIMARY KEY (id);


--
-- Name: file_metadata file_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.file_metadata
    ADD CONSTRAINT file_metadata_pkey PRIMARY KEY (file_id);


--
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (uid);


--
-- PostgreSQL database dump complete
--

