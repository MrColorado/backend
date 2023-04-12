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
    id integer NOT NULL,
    title character varying NOT NULL,
    author character varying NOT NULL,
    description character varying NOT NULL,
    nb_chapter integer NOT NULL,
    first_chapter character varying NOT NULL
);
ALTER TABLE public.novel OWNER TO root_user;
--
-- Name: novel_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.novel_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.novel_id_seq OWNER TO root_user;
--
-- Name: novel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.novel_id_seq OWNED BY public.novel.id;
--
-- Name: novel id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ALTER COLUMN id
SET DEFAULT nextval('public.novel_id_seq'::regclass);
--
-- Name: novel novel_pl; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_pl PRIMARY KEY (id);
--
-- Name: novel novel_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_un UNIQUE (title);
--
-- PostgreSQL database dump complete