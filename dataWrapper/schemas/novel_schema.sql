--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Debian 15.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.2 (Debian 15.2-1.pgdg110+1)
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
-- Name: novel; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.novel (
    title character varying NOT NULL,
    author character varying NOT NULL,
    description character varying NOT NULL,
    nb_chapter integer NOT NULL,
    first_chapter character varying NOT NULL,
    current_chapter integer NOT NULL,
    next_url character varying NOT NULL
);
ALTER TABLE public.novel OWNER TO root_user;
--
-- Name: novel novel_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_pk PRIMARY KEY (title);
--
-- PostgreSQL database dump complete