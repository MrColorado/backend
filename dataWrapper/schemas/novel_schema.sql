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
-- Name: author; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.author (
    id integer NOT NULL,
    name character varying NOT NULL
);
ALTER TABLE public.author OWNER TO root_user;
--
-- Name: author_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.author_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.author_id_seq OWNER TO root_user;
--
-- Name: author_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.author_id_seq OWNED BY public.author.id;
--
-- Name: book; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.book (
    id integer NOT NULL,
    fk_novel_id integer NOT NULL,
    start integer NOT NULL,
    "end" integer NOT NULL,
    last_update timestamp without time zone NOT NULL
);
ALTER TABLE public.book OWNER TO root_user;
--
-- Name: chapter_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.chapter_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.chapter_id_seq OWNER TO root_user;
--
-- Name: chapter_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.chapter_id_seq OWNED BY public.book.id;
--
-- Name: genre; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.genre (
    id integer NOT NULL,
    name character varying NOT NULL
);
ALTER TABLE public.genre OWNER TO root_user;
--
-- Name: genre_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.genre_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.genre_id_seq OWNER TO root_user;
--
-- Name: genre_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.genre_id_seq OWNED BY public.genre.id;
--
-- Name: novel; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.novel (
    id integer NOT NULL,
    title character varying NOT NULL,
    description character varying NOT NULL,
    nb_chapter integer NOT NULL,
    first_chapter character varying NOT NULL,
    current_chapter integer NOT NULL,
    next_url character varying NOT NULL,
    last_update timestamp without time zone DEFAULT now() NOT NULL,
    cover_path character varying NOT NULL,
    fk_author_id integer NOT NULL
);
ALTER TABLE public.novel OWNER TO root_user;
--
-- Name: novel_genre_map; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.novel_genre_map (
    fk_novel_id integer NOT NULL,
    fk_genre_id integer NOT NULL
);
ALTER TABLE public.novel_genre_map OWNER TO root_user;
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
-- Name: novel_tag_map; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.novel_tag_map (
    fk_novel_id integer NOT NULL,
    fk_tag_id integer NOT NULL
);
ALTER TABLE public.novel_tag_map OWNER TO root_user;
--
-- Name: novel_tag_map_fk_novel_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.novel_tag_map_fk_novel_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.novel_tag_map_fk_novel_id_seq OWNER TO root_user;
--
-- Name: novel_tag_map_fk_novel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.novel_tag_map_fk_novel_id_seq OWNED BY public.novel_tag_map.fk_novel_id;
--
-- Name: tag; Type: TABLE; Schema: public; Owner: root_user
--

CREATE TABLE public.tag (
    id integer NOT NULL,
    name character varying NOT NULL
);
ALTER TABLE public.tag OWNER TO root_user;
--
-- Name: tag_id_seq; Type: SEQUENCE; Schema: public; Owner: root_user
--

CREATE SEQUENCE public.tag_id_seq AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
ALTER TABLE public.tag_id_seq OWNER TO root_user;
--
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root_user
--

ALTER SEQUENCE public.tag_id_seq OWNED BY public.tag.id;
--
-- Name: author id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.author
ALTER COLUMN id
SET DEFAULT nextval('public.author_id_seq'::regclass);
--
-- Name: book id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.book
ALTER COLUMN id
SET DEFAULT nextval('public.chapter_id_seq'::regclass);
--
-- Name: genre id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.genre
ALTER COLUMN id
SET DEFAULT nextval('public.genre_id_seq'::regclass);
--
-- Name: novel id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ALTER COLUMN id
SET DEFAULT nextval('public.novel_id_seq'::regclass);
--
-- Name: tag id; Type: DEFAULT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.tag
ALTER COLUMN id
SET DEFAULT nextval('public.tag_id_seq'::regclass);
--
-- Name: author author_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.author
ADD CONSTRAINT author_pk PRIMARY KEY (id);
--
-- Name: author author_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.author
ADD CONSTRAINT author_un UNIQUE (name);
--
-- Name: book book_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.book
ADD CONSTRAINT book_un UNIQUE (fk_novel_id, start);
--
-- Name: book chapter_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.book
ADD CONSTRAINT chapter_pk PRIMARY KEY (id);
--
-- Name: genre genre_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.genre
ADD CONSTRAINT genre_pk PRIMARY KEY (id);
--
-- Name: genre genre_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.genre
ADD CONSTRAINT genre_un UNIQUE (name);
--
-- Name: novel_genre_map novel_genre_map_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_genre_map
ADD CONSTRAINT novel_genre_map_pk PRIMARY KEY (fk_novel_id, fk_genre_id);
--
-- Name: novel novel_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_pk PRIMARY KEY (id);
--
-- Name: novel_tag_map novel_tag_map_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_tag_map
ADD CONSTRAINT novel_tag_map_pk PRIMARY KEY (fk_novel_id, fk_tag_id);
--
-- Name: novel_tag_map novel_tag_map_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_tag_map
ADD CONSTRAINT novel_tag_map_un UNIQUE (fk_novel_id, fk_tag_id);
--
-- Name: novel novel_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_un UNIQUE (title);
--
-- Name: tag tag_pk; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.tag
ADD CONSTRAINT tag_pk PRIMARY KEY (id);
--
-- Name: tag tag_un; Type: CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.tag
ADD CONSTRAINT tag_un UNIQUE (name);
--
-- Name: book book_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.book
ADD CONSTRAINT book_fk FOREIGN KEY (fk_novel_id) REFERENCES public.novel(id);
--
-- Name: novel novel_author_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel
ADD CONSTRAINT novel_author_fk FOREIGN KEY (fk_author_id) REFERENCES public.author(id);
--
-- Name: novel_genre_map novel_genre_map_genre_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_genre_map
ADD CONSTRAINT novel_genre_map_genre_fk FOREIGN KEY (fk_genre_id) REFERENCES public.genre(id);
--
-- Name: novel_genre_map novel_genre_map_novel_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_genre_map
ADD CONSTRAINT novel_genre_map_novel_fk FOREIGN KEY (fk_novel_id) REFERENCES public.novel(id);
--
-- Name: novel_tag_map novel_tag_map_novel_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_tag_map
ADD CONSTRAINT novel_tag_map_novel_fk FOREIGN KEY (fk_novel_id) REFERENCES public.novel(id);
--
-- Name: novel_tag_map novel_tag_map_tag_fk; Type: FK CONSTRAINT; Schema: public; Owner: root_user
--

ALTER TABLE ONLY public.novel_tag_map
ADD CONSTRAINT novel_tag_map_tag_fk FOREIGN KEY (fk_tag_id) REFERENCES public.tag(id);
--
-- PostgreSQL database dump complete
- -