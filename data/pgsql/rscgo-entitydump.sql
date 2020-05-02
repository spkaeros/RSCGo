--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Debian 12.2-1)
-- Dumped by pg_dump version 12.2 (Debian 12.2-1)

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
-- Name: boundarys; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.boundarys (
    id bigint NOT NULL,
    name text,
    description text,
    command_one text,
    command_two text,
    height bigint,
    color1 bigint,
    color2 bigint,
    solid bigint,
    door bigint
);


ALTER TABLE public.boundarys OWNER TO zach;

--
-- Name: game_object_locations; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.game_object_locations (
    id text,
    x text,
    y text,
    direction text,
    boundary text
);


ALTER TABLE public.game_object_locations OWNER TO zach;

--
-- Name: game_objects; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.game_objects (
    id bigint NOT NULL,
    name text,
    description text,
    command_one text,
    command_two text,
    type bigint,
    width bigint,
    height bigint,
    modelheight bigint
);


ALTER TABLE public.game_objects OWNER TO zach;

--
-- Name: item_locations; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.item_locations (
    id bigint,
    x bigint,
    y bigint,
    amount bigint,
    respawn bigint
);


ALTER TABLE public.item_locations OWNER TO zach;

--
-- Name: item_wieldable; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.item_wieldable (
    id bigint NOT NULL,
    sprite bigint,
    type bigint,
    armour_points bigint,
    magic_points bigint,
    prayer_points bigint,
    range_points bigint,
    weapon_aim_points bigint,
    weapon_power_points bigint,
    pos bigint,
    femaleonly boolean
);


ALTER TABLE public.item_wieldable OWNER TO zach;

--
-- Name: item_wieldable_requirements; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.item_wieldable_requirements (
    id bigint,
    skillindex bigint,
    level bigint
);


ALTER TABLE public.item_wieldable_requirements OWNER TO zach;

--
-- Name: items; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.items (
    id bigint NOT NULL,
    name text,
    description text,
    command text,
    base_price bigint,
    stackable boolean,
    special boolean,
    members boolean
);


ALTER TABLE public.items OWNER TO zach;

--
-- Name: npc_drops; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.npc_drops (
    npcid bigint,
    itemid bigint,
    minamount bigint,
    maxamount bigint,
    probability double precision
);


ALTER TABLE public.npc_drops OWNER TO zach;

--
-- Name: npc_locations; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.npc_locations (
    id text,
    startx text,
    minx text,
    maxx text,
    starty text,
    miny text,
    maxy text
);


ALTER TABLE public.npc_locations OWNER TO zach;

--
-- Name: npcs; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.npcs (
    id bigint NOT NULL,
    name text,
    description text,
    command text,
    hits bigint,
    attack bigint,
    strength bigint,
    defense bigint,
    hostility integer DEFAULT 0
);


ALTER TABLE public.npcs OWNER TO zach;

--
-- Name: prayers; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.prayers (
    id bigint NOT NULL,
    name text,
    description text,
    required_level bigint,
    drain_rate bigint
);


ALTER TABLE public.prayers OWNER TO zach;

--
-- Name: shop_items; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.shop_items (
    storeid bigint,
    itemid bigint,
    amount bigint
);


ALTER TABLE public.shop_items OWNER TO zach;

--
-- Name: shops; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.shops (
    id bigint NOT NULL,
    name text,
    general boolean
);


ALTER TABLE public.shops OWNER TO zach;

--
-- Name: spell_aggressive_level; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.spell_aggressive_level (
    id bigint NOT NULL,
    spell bigint
);


ALTER TABLE public.spell_aggressive_level OWNER TO zach;

--
-- Name: spell_runes; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.spell_runes (
    spellid bigint,
    itemid bigint,
    amount bigint
);


ALTER TABLE public.spell_runes OWNER TO zach;

--
-- Name: spells; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.spells (
    id bigint NOT NULL,
    name text,
    description text,
    required_level bigint,
    rune_amount bigint,
    type bigint,
    experience bigint
);


ALTER TABLE public.spells OWNER TO zach;

--
-- Name: tiles; Type: TABLE; Schema: public; Owner: zach
--

CREATE TABLE public.tiles (
    colour bigint,
    unknown bigint,
    objecttype bigint
);


ALTER TABLE public.tiles OWNER TO zach;

--
-- Data for Name: boundarys; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.boundarys (id, name, description, command_one, command_two, height, color1, color2, solid, door) FROM stdin;
144	plankstimber		walkto	examine	192	46	46	1	0
0	Wall		walkto	examine	192	2	2	1	0
1	Doorframe		walkto	close	192	4	4	0	1
2	Door	The door is shut	open	examine	192	0	0	1	1
3	Window		walkto	examine	192	5	5	1	0
4	Fence		walkto	examine	192	10	10	1	0
5	railings		walkto	examine	192	12	12	1	0
6	Stained glass window		walkto	examine	192	18	18	1	0
7	Highwall		walkto	examine	275	2	2	1	0
8	Door	The door is shut	open	examine	275	0	0	1	1
9	Doorframe		walkto	close	275	4	4	0	1
10	battlement		walkto	examine	70	2	2	1	0
11	Doorframe		walkto	examine	192	4	4	1	0
12	snowwall		walkto	examine	192	-31711	-31711	1	0
13	arrowslit		walkto	examine	192	7	7	1	0
14	timberwall		walkto	examine	192	21	21	1	0
15	timberwindow		walkto	examine	192	22	22	1	0
16	blank		walkto	examine	192	12345678	12345678	0	0
17	highblank		walkto	examine	275	12345678	12345678	0	0
18	mossybricks		walkto	examine	192	23	23	1	0
19	Door	The door is shut	open	examine	192	0	0	1	1
20	Door	The door is shut	open	examine	192	0	0	1	1
21	Door	The door is shut	open	examine	192	0	0	1	1
22	Odd looking wall	This wall doesn't look quite right	push	examine	192	2	2	1	1
23	Door	The door is shut	open	examine	192	0	0	1	1
24	web	A spider's web	walkto	examine	192	26	26	1	1
25	Door	The door is shut	open	examine	192	0	0	1	1
26	Door	The door is shut	open	examine	192	0	0	1	1
27	Door	The door is shut	open	examine	192	0	0	1	1
28	Door	The door is shut	open	examine	192	0	0	1	1
29	Door	The door is shut	open	examine	192	0	0	1	1
30	Door	The door is shut	open	examine	192	0	0	1	1
31	Door	The door is shut	open	examine	192	0	0	1	1
32	Door	The door is shut	open	examine	192	0	0	1	1
33	Door	The door is shut	open	examine	192	0	0	1	1
34	Window		walkto	examine	192	27	27	1	0
35	Door	The door is shut	open	examine	192	0	0	1	1
36	Door	The door is shut	open	examine	192	0	0	1	1
37	Door	The door is shut	open	examine	192	0	0	1	1
38	Door	The door is shut	open	examine	192	0	0	1	1
39	Door	The door is shut	open	examine	192	0	0	1	1
40	Door	The door is shut	open	examine	192	0	0	1	1
41	Crumbled		walkto	examine	192	28	28	1	0
42	Cavern		walkto	examine	192	29	29	1	0
43	Door	The door is shut	open	examine	192	0	0	1	1
44	Door	The door is shut	open	examine	192	0	0	1	1
45	Door	The door is shut	open	examine	192	0	0	1	1
46	cavern2		walkto	examine	192	30	30	1	0
47	Door	The door is shut	open	examine	192	0	0	1	1
48	Door	The door is shut	open	examine	192	0	0	1	1
49	Door	The door is shut	open	examine	192	0	0	1	1
50	Door	The door is shut	open	examine	192	0	0	1	1
51	Door	The door is shut	open	examine	192	0	0	1	1
52	Door	The door is shut	open	examine	192	0	0	1	1
53	Door	The door is shut	open	examine	192	0	0	1	1
54	Door	The door is shut	open	examine	192	0	0	1	1
55	Door	The door is shut	open	examine	192	0	0	1	1
56	Wall		walkto	examine	192	3	3	1	0
57	Door	The door is shut	open	examine	192	0	0	1	1
58	Strange looking wall	This wall doesn't look quite right	push	examine	192	29	29	1	1
59	Door	The door is shut	open	examine	192	0	0	1	1
60	Door	The door is shut	open	examine	192	0	0	1	1
61	Door	The door is shut	open	examine	192	0	0	1	1
62	memberrailings		walkto	examine	192	12	12	1	0
63	Door	The door is shut	open	examine	192	0	0	1	1
64	Door	The door is shut	open	examine	192	0	0	1	1
65	Magic Door	The door is shut	open	examine	192	0	0	1	1
66	Door	The door is shut	open	examine	192	0	0	1	1
67	Door	The door is shut	open	examine	192	0	0	1	1
68	Door	The door is shut	open	examine	192	0	0	1	1
69	Door	The door is shut	open	examine	192	0	0	1	1
70	Door	The door is shut	open	examine	192	0	0	1	1
71	Door	The door is shut	open	examine	192	0	0	1	1
72	Door	The door is shut	open	examine	192	0	0	1	1
73	Door	The door is shut	open	examine	192	0	0	1	1
74	Door	The door is shut	open	examine	192	0	0	1	1
75	Door	The door is shut	open	examine	192	0	0	1	1
76	Door	The door is shut	open	examine	192	0	0	1	1
77	Door	The door is shut	open	examine	192	0	0	1	1
78	Door	The door is shut	open	examine	192	0	0	1	1
79	Strange Panel	This wall doesn't look quite right	push	examine	192	21	21	1	1
80	Door	The door is shut	open	examine	192	0	0	1	1
81	Door	The door is shut	open	examine	192	0	0	1	1
82	Door	The door is shut	open	examine	192	0	0	1	1
83	Door	The door is shut	open	examine	192	0	0	1	1
84	Door	The door is shut	open	examine	192	0	0	1	1
85	Door	The door is shut	open	examine	192	0	0	1	1
86	blockblank		walkto	examine	192	12345678	12345678	1	0
87	unusual looking wall	This wall doesn't look quite right	push	examine	192	2	2	1	1
88	Door	The door is shut	open	examine	192	0	0	1	1
89	Door	The door is shut	open	examine	192	0	0	1	1
90	Door	The door is shut	open	examine	192	0	0	1	1
91	Door	The door is shut	open	examine	192	0	0	1	1
92	Door	The door is shut	open	examine	192	0	0	1	1
93	Door	The door is shut	open	pick lock	192	0	0	1	1
94	Door	The door is shut	open	pick lock	192	0	0	1	1
95	Door	The door is shut	open	pick lock	192	0	0	1	1
96	Door	The door is shut	open	pick lock	192	0	0	1	1
97	Door	The door is shut	open	pick lock	192	0	0	1	1
98	Door	The door is shut	open	examine	192	0	0	1	1
99	Door	The door is shut	open	pick lock	192	0	0	1	1
100	Door	The door is shut	open	pick lock	192	0	0	1	1
101	Fence with loose pannels	I wonder if I could get through this	push	examine	192	10	10	1	1
102	Door	The door is shut	open	examine	192	0	0	1	1
103	Door	The door is shut	open	examine	192	0	0	1	1
104	Door	The door is shut	open	examine	192	0	0	1	1
105	Door	The door is shut	open	examine	192	0	0	1	1
106	Door	The door is shut	open	examine	192	0	0	1	1
107	Door	The door is shut	open	examine	192	0	0	1	1
108	Door	The door is shut	open	examine	192	0	0	1	1
109	Door	The door is shut	open	examine	192	0	0	1	1
110	Door	The door is shut	open	examine	192	0	0	1	1
111	rat cage	The rat's have damaged the bars	search	examine	192	12	12	1	1
112	Door	The door is shut	open	examine	192	0	0	1	1
113	Door	The door is shut	open	examine	192	0	0	1	1
114	Door	The door is shut	open	examine	192	0	0	1	1
115	Door	The door is shut	open	examine	192	0	0	1	1
116	Door	The door is shut	open	examine	192	0	0	1	1
117	Door	The door is shut	open	examine	192	0	0	1	1
118	arrowslit		walkto	examine	192	44	44	1	0
119	solidblank		walkto	examine	192	12345678	12345678	1	0
120	Door	The door is shut	open	examine	192	0	0	1	1
121	Door	The door is shut	open	examine	192	0	0	1	1
122	Door	The door is shut	open	examine	192	0	0	1	1
123	Door	The door is shut	open	examine	192	0	0	1	1
124	loose panel	The panel has worn with age	break	examine	192	3	3	1	1
125	Door	The door is shut	open	examine	192	0	0	1	1
126	plankswindow		walkto	examine	192	45	45	1	0
127	Low Fence		walkto	examine	96	10	10	1	0
128	Door	The door is shut	open	examine	192	0	0	1	1
129	Door	The door is shut	open	examine	192	0	0	1	1
130	Door	The door is shut	open	examine	192	0	0	1	1
131	Door	The door is shut	open	examine	192	0	0	1	1
132	Door	The door is shut	open	examine	192	0	0	1	1
133	Door	The door is shut	open	examine	192	0	0	1	1
134	Door	The door is shut	open	examine	192	0	0	1	1
135	Door	The door is shut	open	examine	192	0	0	1	1
136	Door	The door is shut	open	examine	192	0	0	1	1
137	Cooking pot	Smells good!	walkto	examine	96	10	10	1	1
138	Door	The door is shut	open	examine	192	0	0	1	1
139	Door	The door is shut	open	examine	192	0	0	1	1
140	Door	The door is shut	open	examine	192	0	0	1	1
141	Door	The door is shut	open	examine	192	0	0	1	1
142	Door	The door is shut	open	examine	192	0	0	1	1
143	Door	The door is shut	open	examine	192	0	0	1	1
145	Door	The door is shut	open	examine	192	0	0	1	1
146	Door	The door is shut	open	examine	192	0	0	1	1
147	magic portal		enter	examine	192	17	17	1	1
148	magic portal		enter	examine	192	17	17	1	1
149	magic portal		enter	examine	192	17	17	1	1
150	Door	The door is shut	open	examine	192	0	0	1	1
151	Cavern wall	It looks as if it is covered in some fungus.	walkto	search	192	29	29	1	1
152	Door	The door is shut	open	examine	192	0	0	1	1
153	Door	the door is shut	walk through	examine	192	3	3	1	1
154	Door	The door is shut	walk through	examine	192	0	0	1	1
155	Door	The door is shut	walk through	examine	192	0	0	1	1
156	Door	The door is shut	walk through	examine	192	0	0	1	1
157	Door	The door is shut	walk through	examine	192	0	0	1	1
158	Door	The door is shut	walk through	examine	192	0	0	1	1
159	Door	The door is shut	walk through	examine	192	0	0	1	1
160	Door	The door is shut	walk through	examine	192	0	0	1	1
161	Door	The door is shut	open	examine	192	0	0	1	1
162	Door	The door is shut	open	pick lock	192	0	0	1	1
163	Low wall		jump	examine	70	2	2	1	1
164	Low wall		jump	examine	70	2	2	1	1
165	Blacksmiths Door	The door is shut	open	examine	192	0	0	1	1
166	railings		pick lock	examine	192	12	12	1	1
167	railings		pick lock	examine	192	12	12	1	1
168	railings		pick lock	search	192	12	12	1	1
169	railings		pick lock	search	192	12	12	1	1
170	railings		pick lock	search	192	12	12	1	1
171	railings		walkto	search	192	12	12	1	1
172	railings		walkto	look through	192	12	12	1	1
173	Door	The door is shut	open	knock on	192	0	0	1	1
174	Doorframe		walkto	close	192	4	4	1	1
175	Tent		walkto	examine	192	36	36	1	0
176	Jail Door	The door is shut	open	examine	192	0	0	1	1
177	Jail Door	The door is shut	open	examine	192	0	0	1	1
178	Window	A barred window	walkto	search	192	27	27	1	1
179	magic portal	A magical barrier shimmers with power	walkto	examine	192	17	17	1	1
180	Jail Door	A solid iron gate	open	examine	192	12	12	1	1
181	railings		walkto	search	192	12	12	1	1
182	railings		walkto	search	192	12	12	1	1
183	railings		walkto	search	192	12	12	1	1
184	railings		walkto	search	192	12	12	1	1
185	railings		walkto	search	192	12	12	1	1
186	railings		walkto	search	192	12	12	1	1
187	Cave exit	The way out	leave	examine	192	26	26	0	1
188	Cave exit	The way out	leave	examine	192	26	26	0	1
189	Cave exit	The way out	leave	examine	192	26	26	0	1
190	Cave exit	The way out	leave	examine	192	26	26	0	1
191	Cave exit	The way out	leave	examine	192	26	26	0	1
192	Cave exit	The way out	leave	examine	192	26	26	0	1
193	railings		walkto	search	192	12	12	1	1
194	Door	The door is shut	open	examine	192	0	0	1	1
195	battlement	This is blocking your path	climb-over	examine	70	2	2	1	1
196	Tent Door	An entrance into the tent	go through	examine	192	50	50	1	1
197	Door	The door is shut	open	examine	192	0	0	1	1
198	Tent Door	An entrance into the tent	go through	examine	192	50	50	1	1
199	Low Fence	A damaged wooden fence	search	examine	96	10	10	1	1
200	Sturdy Iron Gate	A solid iron gate	open	examine	192	12	12	1	1
201	battlement	this low wall blocks your path	climb over	examine	70	2	2	1	1
202	Water	My waterfall boundary!	walkto	examine	192	25	25	1	0
203	Wheat	Test Boundary!	walkto	examine	192	24	24	1	0
204	Jungle	Thick inpenetrable jungle	chop	examine	192	8	8	1	1
205	Window	you can see a vicious looking guard dog right outside	investigate	examine	192	5	5	1	1
206	Rut	Looks like a small rut carved into the ground.	walkto	search	96	51	51	1	0
207	Crumbled Cavern 1		walkto	examine	192	52	52	0	0
208	Crumbled Cavern 2		walkto	examine	192	53	53	0	0
209	cavernhole		walkto	examine	192	54	54	1	0
210	flamewall	A supernatural fire of incredible intensity	touch	investigate	192	54	54	1	1
211	Ruined wall	Some ancient wall structure - it doesn't look too high.	walkto	jump	192	28	28	1	1
212	Ancient Wall	An ancient - slightly higher wall with some strange markings on it	use	search	275	2	2	1	1
213	Door	The door is shut	open	examine	192	0	0	1	1
\.


--
-- Data for Name: game_object_locations; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.game_object_locations (id, x, y, direction, boundary) FROM stdin;
1	332	565	0	1
1	315	561	3	1
1	326	540	3	1
1	346	554	0	0
1	317	566	0	1
693	339	555	6	0
29	332	566	0	0
1	312	567	0	1
7	323	545	2	0
63	327	552	4	0
1	340	562	0	0
29	332	555	2	0
3	322	546	2	0
16	320	548	2	0
7	317	543	4	0
1	346	570	0	0
1	349	545	0	0
1	350	560	0	0
1	320	543	0	1
1	350	572	0	0
29	332	550	2	0
1	312	544	1	1
139	327	550	4	0
3	317	544	0	0
45	313	556	2	0
7	323	546	2	0
26	313	539	0	0
16	321	548	2	0
45	313	555	2	0
63	313	558	2	0
45	313	553	2	0
16	322	548	2	0
7	322	544	4	0
145	312	571	0	0
41	316	546	0	0
3	322	545	2	0
5	318	561	0	0
45	313	554	2	0
1	350	537	0	0
5	333	569	0	0
1	324	547	1	1
0	345	538	0	0
5	328	538	7	0
26	327	545	4	0
5	317	568	0	0
0	343	544	0	0
0	110	633	0	0
0	109	638	0	0
0	108	638	0	0
0	109	636	0	0
1	106	642	0	0
0	105	638	0	0
34	111	638	4	0
1	110	642	0	0
34	106	647	4	0
0	108	633	0	0
34	110	635	4	0
34	105	635	4	0
0	111	641	0	0
1	105	645	0	0
1	105	640	0	0
34	107	632	4	0
0	107	644	0	0
1	106	637	0	0
1	108	641	0	0
1	110	640	0	0
3	110	651	0	0
7	110	650	4	0
1	107	640	0	0
45	108	655	0	0
45	107	655	0	0
1	136	652	0	1
1	130	636	0	0
0	122	632	0	0
89	131	642	4	0
45	107	657	4	0
45	105	655	0	0
23	111	664	4	0
0	106	663	0	0
23	111	662	4	0
19	110	668	0	0
4	104	653	0	0
1	131	646	0	0
0	122	636	0	0
1	129	636	0	0
1	131	633	0	0
1	120	633	0	0
95	121	668	4	0
7	134	655	4	0
37	120	652	0	0
1	130	635	0	0
29	135	641	2	0
38	127	632	0	0
37	120	653	0	0
0	127	668	0	0
12	107	666	0	0
45	106	655	0	0
3	117	669	0	0
23	111	666	4	0
34	116	640	4	0
34	116	641	4	0
23	109	666	4	0
7	135	655	4	0
13	107	668	0	0
12	107	667	0	0
0	117	639	0	0
13	108	666	0	0
34	114	641	1	0
3	112	653	0	0
1	112	636	0	0
1	114	640	0	0
0	114	638	0	0
1	118	637	0	0
12	105	667	0	0
2	117	650	5	0
45	105	657	4	0
1	112	633	0	0
1	118	635	0	0
1	115	635	0	0
34	112	642	0	0
0	119	640	0	0
1	114	634	0	0
12	105	669	0	0
23	109	662	4	0
13	105	668	0	0
5	118	655	0	0
0	112	637	0	0
61	109	658	6	0
13	108	671	0	0
5	118	661	0	0
45	106	657	4	0
23	109	664	4	0
45	108	657	4	0
7	120	665	5	0
0	128	634	0	0
131	125	646	3	0
1	135	647	0	0
1	128	637	0	0
20	127	666	0	0
15	120	666	0	0
37	121	652	0	0
26	124	652	0	0
0	129	633	0	0
1	127	637	0	0
119	131	660	0	0
20	123	650	0	0
11	115	669	0	0
1	127	636	0	0
29	123	668	6	0
7	135	657	0	0
3	134	660	0	0
25	130	650	0	0
7	133	655	4	0
25	134	665	0	0
25	130	665	0	0
63	128	658	4	0
7	133	657	0	0
25	132	665	0	0
7	134	657	0	0
3	133	661	0	0
25	134	650	0	0
3	138	668	0	0
1	140	641	0	0
1	137	643	0	0
5	139	666	0	0
1	111	654	0	1
0	129	668	0	0
1	131	666	0	0
1	129	667	0	0
27	137	652	0	0
10	132	656	0	0
9	133	656	0	0
0	142	642	0	0
1	132	648	0	0
25	132	650	0	0
7	136	657	0	0
1	128	666	0	0
1	138	639	0	0
1	142	645	0	0
5	139	648	0	0
0	143	635	0	0
1	132	659	0	1
1	117	666	3	1
7	136	655	4	0
1	115	664	1	1
6	136	660	0	0
1	136	663	0	1
1	116	669	0	1
1	117	633	0	1
1	132	641	2	1
1	122	669	1	1
1	111	660	0	1
1	120	661	1	1
1	124	664	0	1
1	120	655	1	1
34	103	638	4	0
4	103	647	0	0
34	103	641	4	0
34	103	644	4	0
4	100	656	0	0
12	103	668	0	0
12	103	670	0	0
13	103	669	0	0
4	102	658	0	0
4	103	635	0	0
12	103	671	0	0
4	103	650	0	0
180	92	649	0	0
164	89	662	5	0
618	88	663	5	0
164	90	662	5	0
12	110	676	0	0
39	103	675	0	0
12	103	673	0	0
13	110	674	0	0
12	110	675	0	0
13	110	672	0	0
13	106	674	0	0
12	105	673	0	0
12	108	676	0	0
13	108	675	0	0
12	108	673	0	0
1	106	675	1	1
13	105	672	0	0
34	87	649	0	0
32	86	667	0	0
34	86	641	0	0
29	84	675	0	0
30	81	664	4	0
164	86	662	5	0
118	85	679	4	0
34	85	644	0	0
34	83	650	0	0
34	86	659	0	0
164	87	664	5	0
34	86	656	0	0
1	83	674	1	1
32	73	646	0	0
32	79	638	0	0
33	75	665	4	0
32	74	663	0	0
32	79	640	0	0
34	74	664	0	0
33	78	654	0	0
11	73	669	0	0
1	75	667	0	1
29	76	676	2	0
5	78	677	2	0
32	72	673	2	0
29	76	674	2	0
1	76	679	0	1
1	74	676	1	1
34	69	663	0	0
35	71	649	0	0
30	71	667	6	0
35	68	656	0	0
32	71	670	6	0
3	66	677	3	0
7	67	676	3	0
3	64	667	4	0
2	71	679	2	0
33	65	673	2	0
22	66	666	4	0
22	67	666	4	0
33	67	664	4	0
1	68	669	1	1
1	67	679	2	1
46	68	685	0	0
35	64	683	0	0
26	71	687	0	0
32	75	681	2	0
33	90	683	0	0
34	84	684	0	0
32	77	681	2	0
29	87	682	2	0
46	76	685	0	0
46	68	687	0	0
11	87	685	2	0
46	76	687	0	0
15	80	685	2	0
1	85	683	1	1
64	71	694	2	0
46	68	693	4	0
46	68	691	0	0
46	70	693	4	0
46	74	693	4	0
46	72	693	2	0
26	72	689	0	0
63	86	695	4	0
15	80	691	2	0
139	86	693	4	0
15	80	688	2	0
46	76	691	0	0
63	77	688	4	0
46	76	689	2	0
46	76	693	4	0
3	78	693	0	0
3	79	693	0	0
63	67	688	0	0
46	68	689	2	0
9	77	698	0	0
9	77	696	0	0
33	82	706	0	0
33	75	708	0	0
35	86	709	0	0
35	78	711	0	0
33	66	704	0	0
33	90	709	0	0
35	82	708	0	0
35	67	714	0	0
33	76	718	0	0
193	89	718	0	0
35	80	717	0	0
35	83	713	0	0
193	85	719	0	0
55	64	725	2	0
35	96	726	7	0
15	63	691	6	0
35	59	699	0	0
55	57	688	0	0
55	60	725	5	0
34	56	719	2	0
34	57	717	2	0
34	61	688	0	0
34	56	714	2	0
34	57	696	0	0
34	56	716	2	0
15	63	688	6	0
35	58	688	0	0
33	53	692	0	0
8	52	714	0	0
8	53	714	2	0
34	54	713	2	0
34	55	718	2	0
34	49	713	2	0
11	52	718	2	0
34	49	710	2	0
34	52	711	2	0
34	50	712	2	0
1	51	715	0	1
1	55	693	0	1
35	52	698	0	0
3	54	688	0	0
21	60	728	2	0
35	51	735	0	0
916	62	733	6	0
55	59	730	5	0
942	58	731	7	0
21	60	732	2	0
35	84	731	0	0
35	78	734	0	0
21	64	728	0	0
21	64	732	0	0
55	59	728	7	0
176	66	729	1	1
35	81	740	0	0
35	74	740	0	0
944	66	741	0	0
35	62	738	0	0
35	61	749	0	0
35	70	758	0	0
35	79	752	3	0
35	55	761	0	0
35	75	766	3	0
35	78	782	0	0
35	51	783	0	0
35	62	786	0	0
35	74	788	0	0
35	61	797	0	0
35	70	806	0	0
35	55	809	0	0
7	84	800	5	0
35	84	779	0	0
3	84	801	4	0
958	81	801	7	0
105	81	807	0	0
101	82	809	4	0
5	86	799	0	0
103	86	812	0	0
111	83	812	4	0
35	81	788	0	0
105	81	806	0	0
101	83	810	4	0
1025	84	807	6	0
111	84	812	4	0
3	85	801	6	0
103	86	811	0	0
1	85	804	0	1
35	87	766	3	0
35	81	761	3	0
35	88	749	3	0
35	89	741	0	0
35	91	730	0	0
35	93	753	3	0
35	92	758	4	0
945	90	763	3	0
35	91	778	0	0
35	89	789	0	0
177	88	801	1	1
932	92	807	0	0
954	92	800	0	0
953	91	801	4	0
178	90	802	1	1
35	103	791	4	0
35	103	770	5	0
35	101	802	2	0
35	99	779	5	0
35	98	812	1	0
967	90	811	0	0
967	90	810	0	0
966	89	810	0	0
35	110	777	6	0
35	108	811	3	0
35	106	765	1	0
35	99	761	2	0
35	101	750	0	0
35	106	746	0	0
35	99	740	0	0
35	106	734	0	0
35	109	725	0	0
106	110	705	0	0
36	111	706	0	0
106	110	704	0	0
35	117	736	0	0
35	118	748	6	0
3	117	712	0	0
35	112	756	0	0
35	118	766	5	0
35	120	728	0	0
35	126	755	0	0
35	125	771	5	0
35	117	770	6	0
35	112	787	7	0
35	116	797	0	0
35	119	810	1	0
35	130	810	2	0
35	129	798	3	0
35	128	784	4	0
35	134	772	6	0
35	133	765	5	0
35	128	744	0	0
35	135	743	3	0
35	128	731	0	0
35	138	734	0	0
35	137	754	0	0
35	141	775	0	0
35	137	781	7	0
35	144	774	7	0
35	147	788	0	0
35	138	792	1	0
35	149	798	0	0
35	139	802	2	0
35	147	809	2	0
35	142	814	3	0
35	154	813	1	0
35	154	794	0	0
35	154	782	0	0
35	157	773	0	0
35	165	784	0	0
35	166	796	6	0
395	166	798	2	0
35	160	804	0	0
35	166	814	5	0
21	170	806	0	0
21	172	793	4	0
1006	170	792	4	0
21	172	804	0	0
405	172	800	2	0
21	170	807	0	0
21	169	791	0	0
32	174	798	4	0
21	172	791	6	0
33	170	796	2	0
395	169	801	4	0
21	173	804	0	0
21	173	809	2	0
21	172	809	2	0
21	175	807	4	0
21	169	793	2	0
21	175	806	4	0
196	170	805	2	1
198	171	794	0	1
35	176	792	0	0
35	181	813	5	0
35	183	791	3	0
35	168	776	0	0
35	176	779	0	0
35	186	782	0	0
70	197	759	0	0
35	185	802	0	0
5	54	683	2	0
29	56	680	2	0
89	59	682	0	0
29	56	682	2	0
55	56	685	0	0
15	63	685	6	0
32	50	682	2	0
35	59	687	0	0
1	58	686	1	1
1	59	681	1	1
32	55	676	2	0
34	54	673	2	0
55	51	677	2	0
34	51	672	2	0
32	51	679	2	0
48	60	673	0	0
29	50	674	2	0
34	59	676	2	0
55	52	677	2	0
1	53	674	1	1
1	63	673	1	1
32	61	667	0	0
34	63	665	0	0
34	60	668	0	0
32	57	668	0	0
33	57	664	0	0
34	56	666	0	0
34	62	663	0	0
34	62	662	0	0
34	60	661	0	0
35	63	655	0	0
35	62	647	0	0
33	53	639	0	0
33	71	636	0	0
33	68	633	0	0
33	70	633	0	0
35	68	626	0	0
35	58	631	0	0
35	71	621	0	0
35	70	620	0	0
34	56	620	0	0
34	56	611	0	0
34	52	596	0	0
102	68	598	0	0
33	60	598	0	0
110	68	593	0	0
33	55	592	0	0
106	68	588	0	0
34	50	589	0	0
32	55	585	0	0
195	68	585	0	0
33	63	585	0	0
102	68	589	0	0
102	68	586	0	0
110	69	582	0	0
100	70	581	0	0
104	70	583	0	0
32	56	579	0	0
34	38	571	5	0
284	39	573	0	0
0	38	568	7	0
284	38	573	0	0
284	42	571	2	0
284	42	570	2	0
284	42	569	2	0
284	35	573	0	0
283	34	568	0	0
283	34	571	7	0
284	42	572	2	0
284	34	573	0	0
34	36	569	3	0
284	36	573	0	0
284	37	573	0	0
284	33	573	0	0
284	32	573	0	0
283	42	573	2	0
34	32	570	7	0
20	52	575	0	0
284	42	568	2	0
1	51	569	0	0
20	55	575	0	0
1	59	569	0	0
20	49	575	0	0
37	63	569	7	0
20	61	575	0	0
34	55	568	7	0
20	64	575	0	0
1089	59	573	0	0
34	67	574	4	0
20	67	575	0	0
20	58	575	0	0
20	70	575	0	0
34	64	574	4	0
283	40	569	3	0
284	41	573	0	0
284	40	573	0	0
283	39	566	5	0
284	39	562	0	0
284	34	562	0	0
284	33	562	0	0
284	36	562	0	0
284	38	562	0	0
34	37	560	2	0
283	37	561	2	0
284	37	562	0	0
34	35	561	6	0
284	41	562	0	0
283	42	566	2	0
283	42	562	2	0
284	32	562	0	0
283	39	560	2	0
284	35	562	0	0
34	38	561	1	0
283	32	561	2	0
284	42	567	2	0
34	41	561	0	0
0	51	567	7	0
1	51	565	0	0
34	54	562	7	0
0	68	564	7	0
0	57	564	7	0
37	71	561	7	0
34	58	561	7	0
0	53	561	7	0
34	60	561	7	0
37	54	564	7	0
37	50	563	7	0
37	60	563	7	0
34	40	561	2	0
37	55	560	7	0
1064	42	565	6	0
34	61	566	7	0
284	40	562	0	0
37	68	563	7	0
34	66	565	7	0
0	64	563	0	0
34	67	567	7	0
1	68	560	0	0
20	79	575	0	0
98	72	581	0	0
102	75	584	0	0
100	74	582	0	0
1	72	563	0	0
112	75	597	0	0
35	78	596	0	0
195	75	589	0	0
108	72	582	0	0
98	74	587	0	0
102	74	586	0	0
1	78	564	7	0
34	79	560	7	0
133	78	574	0	0
1	77	566	0	0
0	74	566	0	0
20	73	575	0	0
1	75	562	7	0
34	77	562	7	0
35	77	606	0	0
35	73	621	0	0
35	72	619	0	0
35	78	627	0	0
32	85	631	0	0
35	84	607	0	0
70	88	596	0	0
70	84	591	0	0
70	81	579	0	0
133	90	572	6	0
66	93	575	0	0
61	93	572	2	0
34	83	568	7	0
20	82	575	0	0
20	88	575	0	0
20	85	575	0	0
1	86	565	7	0
34	85	561	7	0
34	89	563	7	0
34	91	564	7	0
0	80	562	7	0
4	92	562	7	0
1	81	562	7	0
34	82	562	7	0
1	81	566	0	0
0	83	563	7	0
0	102	563	4	0
34	101	567	7	0
20	99	575	0	0
20	102	575	0	0
191	100	588	0	0
20	96	575	0	0
0	102	585	0	0
0	101	576	4	0
191	100	591	0	0
191	100	583	0	0
191	102	588	0	0
37	102	578	4	0
191	98	594	0	0
191	96	588	0	0
191	96	597	0	0
191	100	597	0	0
191	102	591	0	0
37	97	576	4	0
191	102	597	0	0
191	100	585	0	0
191	98	585	0	0
191	96	594	0	0
191	96	591	0	0
191	98	588	0	0
191	98	597	0	0
1	97	564	4	0
54	96	582	4	0
191	98	591	0	0
1	97	577	4	0
37	99	577	4	0
191	100	594	0	0
191	102	594	0	0
191	102	603	0	0
191	96	603	0	0
191	98	600	0	0
191	98	603	0	0
191	100	600	0	0
191	102	600	0	0
191	96	600	0	0
191	100	603	0	0
0	101	612	0	0
0	97	611	0	0
0	100	619	0	0
106	110	698	0	0
106	110	697	0	0
108	111	699	0	0
0	111	681	0	0
0	111	684	0	0
1	111	680	0	0
36	119	694	6	0
0	113	680	0	0
1	112	682	0	0
36	118	693	6	0
34	113	708	0	0
34	117	705	0	0
34	112	707	0	0
36	116	704	0	0
110	116	703	0	0
110	113	703	0	0
7	117	711	6	0
110	112	702	0	0
110	114	703	0	0
110	115	702	0	0
110	117	703	0	0
34	118	695	6	0
36	112	705	0	0
36	112	704	0	0
1	116	709	0	1
1	111	629	0	0
34	109	631	6	0
34	108	631	4	0
0	108	625	0	0
59	105	619	0	0
1	111	620	4	0
34	109	613	4	0
34	108	614	4	0
34	110	612	4	0
191	104	600	0	0
191	106	603	0	0
59	108	601	0	0
191	106	600	0	0
191	104	603	0	0
34	110	593	1	0
0	104	597	0	0
59	105	587	0	0
34	111	591	1	0
20	107	573	0	0
20	111	575	6	0
20	104	573	0	0
59	105	571	6	0
34	111	563	7	0
34	108	566	7	0
27	56	1623	0	0
15	57	1627	0	0
3	74	1618	0	0
7	75	1618	2	0
6	78	1621	0	0
6	54	1627	2	0
1	59	558	0	0
1	60	557	0	0
34	59	556	7	0
0	61	554	7	0
4	59	553	7	0
0	67	553	7	0
0	89	553	7	0
34	88	553	7	0
34	88	554	7	0
4	92	553	7	0
34	74	552	7	0
4	84	553	7	0
1	91	558	0	0
1	79	554	7	0
1	74	558	0	0
0	89	555	7	0
34	84	557	7	0
1	79	559	7	0
1	90	559	0	0
0	84	554	7	0
1	92	554	7	0
34	75	552	7	0
1	76	552	7	0
1	76	559	0	0
1	93	555	7	0
34	87	559	7	0
34	72	554	7	0
1	67	555	0	0
4	92	559	7	0
4	94	556	7	0
37	75	557	7	0
4	95	552	7	0
1	87	552	7	0
34	72	555	7	0
34	86	556	7	0
1	74	554	0	0
34	79	555	7	0
0	57	556	7	0
1	86	558	0	0
34	68	559	7	0
34	65	554	7	0
1	69	555	0	0
34	67	557	7	0
34	69	554	7	0
1	65	558	0	0
34	81	557	7	0
37	58	554	7	0
34	66	553	7	0
0	54	559	7	0
0	54	553	7	0
34	54	556	7	0
34	53	555	7	0
34	52	559	7	0
34	52	554	7	0
0	50	556	7	0
37	51	552	7	0
283	40	555	0	0
34	40	559	6	0
283	38	555	0	0
37	38	553	0	0
283	37	558	2	0
283	35	554	0	0
37	36	554	0	0
283	35	552	0	0
283	32	553	0	0
283	35	559	2	0
309	37	553	0	0
1	55	551	0	0
34	32	544	0	0
1	36	551	0	0
34	55	548	7	0
34	51	549	7	0
1	60	547	0	0
4	54	550	7	0
34	51	547	7	0
37	52	546	7	0
1	51	551	0	0
34	55	545	7	0
4	61	549	7	0
1	57	549	0	0
37	63	547	7	0
34	62	547	7	0
37	57	545	7	0
100	70	544	4	0
1	65	546	0	0
100	69	544	4	0
104	70	546	4	0
1	52	547	0	0
1	67	551	0	0
34	30	545	0	0
34	28	550	7	0
1072	30	544	0	0
34	29	550	6	0
34	29	544	0	0
34	27	554	2	0
34	25	547	0	0
285	29	566	2	0
33	27	553	5	0
284	31	562	0	0
34	25	550	0	0
284	26	562	0	0
283	30	569	2	0
145	30	572	2	0
145	30	571	2	0
17	28	573	0	0
283	31	570	2	0
145	30	573	2	0
284	28	562	0	0
284	25	562	0	0
284	31	573	0	0
34	26	544	0	0
37	24	547	0	0
17	26	573	0	0
283	24	574	0	0
17	27	573	0	0
25	24	573	6	0
284	29	562	0	0
34	29	553	5	0
1044	24	545	0	0
285	28	567	2	0
284	30	562	0	0
17	29	573	0	0
283	28	574	0	0
285	27	566	2	0
17	25	573	0	0
284	24	562	0	0
284	27	562	0	0
145	30	570	2	0
285	28	565	2	0
283	24	563	0	0
283	25	574	0	0
182	27	571	6	0
71	29	570	6	0
283	24	569	0	0
283	26	574	0	0
182	28	571	6	0
1	24	566	1	1
1	27	570	0	1
402	23	544	6	0
34	22	550	0	0
37	23	548	0	0
402	21	545	6	0
37	21	549	0	0
1090	17	573	2	0
47	17	571	2	0
284	19	562	0	0
71	17	563	6	0
1078	23	569	4	0
34	19	561	0	0
284	17	562	0	0
284	21	562	0	0
29	19	563	2	0
1104	17	569	6	0
24	17	568	2	0
284	18	562	0	0
47	19	573	2	0
1104	17	570	6	0
7	21	563	4	0
47	19	571	2	0
34	21	560	0	0
27	20	563	0	0
3	21	564	2	0
1	19	560	0	0
284	16	562	0	0
29	19	565	2	0
283	22	574	0	0
284	20	562	0	0
278	21	566	6	0
179	22	573	0	0
278	21	569	6	0
1	17	560	0	0
278	19	568	6	0
7	23	565	2	0
283	18	574	0	0
7	23	564	2	0
402	19	552	6	0
284	23	562	0	0
0	18	559	0	0
3	22	564	2	0
71	23	567	4	0
37	17	553	0	0
284	22	562	0	0
402	18	550	6	0
7	22	563	4	0
402	18	553	6	0
402	17	550	6	0
402	17	547	6	0
1044	19	547	2	0
1044	18	549	2	0
34	21	559	0	0
402	19	548	6	0
1072	16	553	2	0
1044	16	551	0	0
402	17	551	6	0
402	20	550	6	0
34	18	551	0	0
1	20	546	1	0
1072	19	551	7	0
1072	17	552	1	0
1044	20	547	0	0
1072	16	549	5	0
37	16	554	0	0
37	19	550	0	0
1	19	549	1	0
402	20	548	6	0
34	22	548	0	0
37	20	549	0	0
1044	22	546	0	0
1042	14	551	0	0
283	12	546	1	0
0	11	547	7	0
1058	9	550	2	0
0	12	556	1	0
1042	14	554	0	0
1058	9	546	7	0
1072	14	552	4	0
34	15	555	0	0
1042	15	552	0	0
284	15	562	0	0
309	10	556	2	0
34	9	555	7	0
34	9	558	0	0
284	15	570	2	0
284	15	566	2	0
0	14	566	0	0
284	15	571	2	0
284	15	569	2	0
0	15	561	0	0
284	15	568	2	0
1042	15	553	0	0
0	14	569	0	0
1072	10	574	2	0
0	11	572	0	0
0	10	571	0	0
1072	10	573	2	0
1072	13	572	5	0
1072	13	570	3	0
1042	13	553	0	0
283	13	571	4	0
1058	9	545	0	0
1042	13	555	0	0
1058	8	544	0	0
98	8	551	5	0
1042	12	554	0	0
98	8	549	0	0
284	15	567	2	0
284	15	564	2	0
284	15	563	2	0
284	15	565	2	0
34	11	559	7	0
284	15	574	2	0
0	14	572	6	0
284	15	572	2	0
0	14	574	0	0
283	13	566	4	0
284	15	573	2	0
0	14	571	0	0
0	13	573	5	0
0	12	570	2	0
1072	9	572	7	0
283	12	568	4	0
0	12	573	4	0
0	12	574	0	0
1072	13	569	3	0
1072	6	565	7	0
1	4	566	2	0
283	1	557	2	0
0	3	563	1	0
34	5	561	7	0
1	2	564	2	0
98	6	553	4	0
1058	7	553	7	0
286	2	553	2	0
34	5	565	3	0
34	6	566	5	0
1058	6	554	0	0
34	3	556	3	0
1058	6	552	4	0
0	2	560	0	0
34	1	567	7	0
0	7	574	0	0
1058	7	551	0	0
34	1	547	7	0
283	1	549	2	0
34	0	546	4	0
1	1	551	2	0
1058	7	550	6	0
1058	6	544	0	0
1058	7	544	0	0
34	3	547	5	0
34	4	548	3	0
1	2	549	2	0
1044	30	540	2	0
1	57	543	0	0
34	59	541	7	0
34	53	538	7	0
34	58	540	7	0
283	51	540	0	0
1	51	539	0	0
1	50	542	0	0
404	27	541	6	0
1044	32	537	0	0
1	55	540	0	0
1056	33	540	5	0
34	54	536	0	0
1050	24	536	6	0
34	24	537	0	0
1044	31	538	2	0
4	63	538	7	0
1056	31	543	5	0
34	63	541	7	0
34	48	539	0	0
401	24	542	6	0
1	59	543	0	0
34	50	536	0	0
283	52	539	0	0
0	53	543	7	0
34	34	540	0	0
1072	31	542	0	0
1072	34	538	0	0
34	32	543	0	0
1073	32	541	0	0
1042	34	536	0	0
404	28	541	4	0
402	25	543	6	0
1044	26	542	2	0
1044	27	542	0	0
1042	23	537	0	0
402	23	543	6	0
1042	22	538	0	0
1075	20	538	0	0
1042	21	538	0	0
1042	20	537	0	0
1083	19	539	0	0
1082	18	540	0	0
1050	19	536	6	0
34	17	542	0	0
1082	19	540	0	0
1075	18	538	1	0
21	16	542	4	0
21	16	540	4	0
1075	19	538	0	0
1075	19	537	0	0
21	16	538	4	0
1058	4	542	5	0
34	5	536	0	0
1182	7	537	5	0
283	0	536	1	0
34	5	537	1	0
283	4	536	1	0
34	4	538	6	0
98	2	540	4	0
34	3	536	4	0
34	2	539	0	0
98	2	538	7	0
3	13	541	0	0
34	3	538	6	0
98	3	537	6	0
34	3	540	2	0
283	3	539	1	0
3	15	541	0	0
34	2	537	6	0
1058	0	540	6	0
98	1	538	3	0
21	10	540	0	0
8	8	540	6	0
97	10	536	4	0
21	14	538	4	0
98	0	539	7	0
8	8	539	4	0
98	1	539	7	0
1058	5	543	0	0
97	11	536	3	0
21	12	538	7	0
1085	10	541	6	0
196	9	537	1	1
21	10	542	2	0
1044	31	528	2	0
34	22	531	4	0
1044	31	531	2	0
1046	8	528	0	0
1067	15	528	4	0
34	23	532	5	0
1046	8	532	0	0
1052	14	534	0	0
1042	18	535	0	0
1044	13	533	0	0
1065	16	528	5	0
34	29	530	0	0
1044	10	528	2	0
1044	10	531	2	0
1	19	534	7	0
1042	25	534	0	0
1056	29	531	5	0
1044	38	534	0	0
1042	18	534	0	0
1066	17	528	6	0
1052	16	534	0	0
1044	16	533	0	0
34	8	533	0	0
34	13	534	0	0
1060	13	528	4	0
1046	8	530	0	0
1075	26	534	2	0
34	28	530	0	0
1042	25	535	0	0
1044	35	535	0	0
34	34	531	0	0
34	36	530	0	0
1042	43	531	0	0
1066	27	528	5	0
34	36	531	0	0
34	28	531	0	0
34	37	530	0	0
1050	25	533	6	0
1067	25	528	3	0
1042	41	529	0	0
1042	41	533	0	0
1042	41	531	0	0
1042	41	530	0	0
1042	42	531	0	0
1060	24	528	4	0
1042	42	532	0	0
1042	42	530	0	0
1042	41	528	0	0
1044	28	533	0	0
1046	6	530	0	0
34	6	535	7	0
1046	6	532	0	0
34	7	533	0	0
0	6	534	7	0
1046	6	528	0	0
34	4	534	6	0
34	3	529	6	0
283	4	531	1	0
402	2	530	2	0
283	0	534	1	0
34	2	532	0	0
402	0	533	5	0
283	0	529	1	0
402	3	528	7	0
0	3	532	2	0
283	1	530	5	0
402	4	530	5	0
402	2	529	4	0
402	1	531	7	0
283	1	534	1	0
1066	15	526	2	0
1044	7	522	2	0
1072	2	520	3	0
1048	13	525	0	0
1048	23	523	2	0
1048	23	522	2	0
1063	23	521	2	0
1066	15	527	2	0
1048	23	524	2	0
1048	13	523	2	0
1048	13	522	2	0
1065	20	524	0	0
1067	14	523	2	0
1067	14	527	0	0
1048	14	525	0	0
1048	13	524	2	0
1065	15	520	2	0
1048	15	525	0	0
1044	7	525	2	0
1067	20	525	0	0
1066	15	522	2	0
215	15	521	5	0
1067	15	523	4	0
1072	1	523	7	0
1067	21	520	2	0
1066	14	521	4	0
1048	13	521	2	0
1067	21	523	2	0
1061	21	525	4	0
1066	21	521	2	0
1048	28	525	0	0
1048	28	524	6	0
1067	14	520	2	0
1044	31	522	2	0
1044	31	525	2	0
1065	13	527	0	0
1066	16	521	0	0
1065	16	520	0	0
1048	16	525	0	0
1048	17	524	6	0
1067	16	523	0	0
1066	19	525	0	0
1067	17	526	0	0
1048	17	525	0	0
1065	20	521	0	0
1066	20	520	2	0
1067	17	527	4	0
1048	17	523	6	0
1065	19	523	0	0
1065	19	522	0	0
1048	17	522	6	0
1067	16	526	0	0
1065	19	521	2	0
1061	17	521	6	0
1067	19	524	2	0
1066	21	522	2	0
1067	20	522	2	0
1065	21	524	0	0
1044	10	522	2	0
1046	33	521	0	0
1065	13	526	2	0
1044	10	525	2	0
1048	38	523	2	0
1048	38	520	2	0
1048	38	521	2	0
1046	33	525	0	0
1046	33	523	0	0
1067	24	520	4	0
1046	35	525	0	0
1048	38	522	2	0
1065	24	523	4	0
1048	38	524	2	0
1048	23	525	0	0
1048	38	527	2	0
1	36	527	7	0
1048	38	525	2	0
1048	38	526	2	0
1067	24	526	5	0
1065	24	527	4	0
1065	24	524	4	0
1048	24	525	0	0
1067	25	526	0	0
1048	26	525	0	0
214	26	523	1	0
1066	24	522	2	0
1067	24	521	2	0
1049	25	521	2	0
1048	25	525	0	0
1065	25	523	4	0
1067	25	524	4	0
1067	27	520	2	0
1067	27	527	2	0
1067	27	521	2	0
1067	26	527	2	0
1048	28	520	6	0
1063	28	522	6	0
1048	28	523	6	0
1066	26	526	2	0
1067	26	522	0	0
1066	27	522	2	0
1065	27	523	0	0
1048	27	525	0	0
1066	27	524	0	0
1048	17	517	6	0
1048	14	515	4	0
1044	23	513	0	0
1048	17	518	6	0
1048	23	518	2	0
1048	13	518	2	0
1067	21	519	0	0
1072	1	517	1	0
1065	14	518	4	0
1048	23	517	4	0
285	2	512	2	0
1048	15	515	4	0
402	11	517	0	0
285	1	514	2	0
1048	13	515	4	0
1044	7	513	2	0
1061	19	517	0	0
1048	13	516	2	0
1048	25	517	4	0
1044	31	513	2	0
1044	31	519	2	0
1044	31	516	2	0
1048	28	519	6	0
1067	26	519	0	0
1048	28	518	6	0
1	27	514	7	0
1059	27	515	4	0
1067	16	519	4	0
1048	17	515	4	0
1048	26	517	4	0
1048	27	517	4	0
1048	17	516	6	0
1053	26	515	2	0
1067	25	515	0	0
1059	26	516	6	0
1065	19	519	4	0
1067	25	514	6	0
1065	25	516	6	0
1044	26	513	0	0
1048	24	517	4	0
1048	23	519	2	0
1066	15	517	4	0
1048	17	519	6	0
1059	23	515	1	0
1067	25	519	0	0
1067	16	517	2	0
1066	24	519	0	0
1066	21	517	2	0
1048	13	517	2	0
34	18	512	0	0
1061	13	519	2	0
1044	20	513	0	0
1066	19	518	2	0
1065	20	517	2	0
1065	20	519	0	0
214	15	518	5	0
1048	28	514	6	0
1065	27	519	0	0
1048	28	517	4	0
1048	28	515	6	0
1066	15	519	4	0
1066	26	514	0	0
1048	16	515	4	0
1059	20	518	4	0
285	0	512	2	0
1044	10	519	2	0
1072	2	514	3	0
1044	10	513	2	0
1062	24	515	2	0
1044	7	516	2	0
1044	10	516	2	0
1044	7	519	2	0
34	4	507	2	0
1044	7	507	2	0
1044	7	510	2	0
1044	7	504	2	0
285	4	506	2	0
1066	15	508	0	0
1	11	508	2	0
34	11	510	0	0
1042	11	511	0	0
1048	13	508	2	0
1065	15	504	0	0
34	2	506	2	0
1048	13	507	2	0
1044	13	511	0	0
34	1	505	7	0
308	1	508	2	0
1059	15	505	6	0
34	0	510	2	0
1067	15	506	0	0
34	1	506	0	0
1065	15	507	0	0
0	21	511	7	0
1062	23	509	4	0
1072	2	510	5	0
285	3	510	2	0
34	3	504	2	0
34	2	507	4	0
34	11	509	4	0
1048	21	510	0	0
1066	23	505	2	0
1066	23	508	4	0
1044	19	511	2	0
1044	10	510	2	0
1048	19	507	2	0
1048	20	510	0	0
1044	10	507	2	0
34	10	508	2	0
1065	23	506	4	0
402	12	504	3	0
1048	13	504	2	0
1095	14	504	4	0
1048	13	505	2	0
1062	14	508	4	0
1066	14	506	0	0
1059	14	507	1	0
1048	19	508	2	0
1067	20	506	0	0
1066	20	507	0	0
1059	20	508	3	0
1048	19	506	2	0
1048	21	504	4	0
1048	13	506	2	0
1067	20	505	0	0
1048	19	510	0	0
1048	19	509	2	0
1048	19	505	2	0
1048	28	509	6	0
1048	16	508	6	0
1048	19	504	4	0
1065	27	505	0	0
1059	24	505	0	0
1044	31	507	2	0
1048	16	505	6	0
1048	27	504	4	0
1066	27	506	2	0
1066	26	506	4	0
1044	16	511	0	0
1065	27	508	2	0
1065	24	507	2	0
1065	26	505	0	0
1067	24	506	4	0
1067	21	505	2	0
1048	28	508	6	0
1065	21	507	0	0
1048	22	504	4	0
1048	22	510	0	0
1	9	504	3	0
1048	23	504	4	0
1067	21	508	4	0
1048	20	504	4	0
1067	22	507	0	0
1067	22	506	4	0
1048	26	510	0	0
1065	25	507	4	0
1065	26	508	4	0
1059	25	508	2	0
1065	26	507	4	0
1048	31	504	4	0
34	29	505	2	0
1048	28	505	6	0
1044	31	510	2	0
1048	26	504	4	0
1048	24	504	4	0
1048	16	507	6	0
1048	25	510	0	0
1048	16	506	6	0
1048	25	504	4	0
1048	28	506	6	0
1048	16	504	6	0
1065	24	508	4	0
1048	28	504	4	0
0	18	505	3	0
1048	30	504	4	0
1048	28	507	6	0
1048	28	510	0	0
1048	29	504	4	0
1048	17	504	6	0
1048	18	504	4	0
1048	27	510	0	0
1048	7	500	2	0
1042	15	499	6	0
1048	7	499	2	0
1076	15	502	0	0
1048	15	503	4	0
1048	7	501	2	0
1048	7	498	2	0
1075	14	502	3	0
34	0	503	2	0
1066	13	497	4	0
1067	14	497	0	0
1057	12	497	4	0
1051	14	498	2	0
1048	14	500	0	0
1067	10	498	0	0
1066	13	496	2	0
1063	11	499	4	0
1048	10	500	0	0
1066	11	496	2	0
1042	9	499	6	0
1066	11	497	2	0
1042	9	497	6	0
1048	7	502	2	0
1042	9	498	6	0
1042	9	496	6	0
1048	19	500	6	0
1048	21	502	4	0
1048	22	502	4	0
1048	22	499	4	0
1048	7	503	2	0
34	1	496	1	0
1048	7	496	2	0
1042	15	498	6	0
1048	22	496	0	0
1042	15	497	6	0
1048	7	497	2	0
1075	16	501	6	0
1048	13	500	0	0
285	2	501	2	0
1048	19	499	6	0
1048	17	499	6	0
1051	9	500	0	0
1048	19	498	6	0
34	9	503	7	0
1048	14	503	4	0
1042	15	496	6	0
1048	19	497	6	0
1048	17	498	6	0
1067	13	498	0	0
1050	10	496	6	0
1048	24	496	0	0
1048	17	503	6	0
1067	10	497	0	0
1065	14	496	0	0
1048	29	499	4	0
1048	13	503	4	0
1048	26	496	0	0
1048	26	499	4	0
1048	28	502	4	0
1048	26	502	4	0
1048	25	502	4	0
1048	23	499	4	0
1048	20	502	4	0
286	8	502	4	0
1048	31	502	4	0
1048	23	496	0	0
1048	17	500	6	0
1048	23	502	4	0
1048	28	499	4	0
1048	19	496	6	0
1048	21	496	0	0
1048	29	502	4	0
1048	20	496	0	0
1048	24	502	4	0
1048	30	499	4	0
1048	25	499	4	0
1048	27	502	4	0
1048	30	502	4	0
1048	30	496	0	0
1048	31	499	4	0
1048	29	496	0	0
1048	28	496	0	0
1048	17	497	6	0
1048	31	496	0	0
1048	17	501	6	0
1048	17	502	6	0
1048	19	502	6	0
1048	17	496	6	0
1075	16	502	2	0
1048	16	503	4	0
1048	19	501	6	0
1048	27	496	0	0
1048	24	499	4	0
1048	27	499	4	0
1048	25	496	0	0
1048	7	495	2	0
1048	7	493	2	0
1048	7	492	4	0
1048	7	494	2	0
34	0	492	1	0
1042	15	495	6	0
1048	13	492	4	0
1065	12	495	2	0
1042	11	495	6	0
1048	17	493	6	0
1048	10	492	4	0
1048	17	492	4	0
34	10	495	3	0
0	18	495	7	0
1048	14	492	4	0
1048	17	494	6	0
1048	16	492	4	0
34	9	493	5	0
217	14	495	1	0
1048	15	492	4	0
1050	14	494	6	0
1042	12	494	6	0
1065	13	495	2	0
1042	11	494	6	0
1048	12	492	4	0
1048	11	492	4	0
1048	8	492	4	0
1048	9	492	4	0
1048	17	495	6	0
1042	13	494	6	0
34	8	493	5	0
34	4	487	1	0
0	6	483	1	0
1072	3	480	5	0
37	13	485	1	0
1072	13	484	7	0
34	3	486	1	0
34	12	486	1	0
0	14	487	1	0
1	10	485	1	0
1072	2	481	7	0
34	9	484	1	0
1042	21	486	2	0
1	15	481	1	0
37	12	482	1	0
37	3	484	1	0
1072	16	486	7	0
1042	25	486	2	0
1	17	486	1	0
38	14	485	1	0
0	12	484	1	0
37	18	483	1	0
0	9	486	1	0
1042	28	484	2	0
1042	24	485	2	0
1072	19	485	0	0
1042	21	484	2	0
34	31	485	7	0
1042	22	483	2	0
0	24	487	1	0
1048	39	496	0	0
1048	37	502	4	0
1048	33	502	4	0
1048	38	510	2	0
1048	33	499	4	0
1048	36	499	4	0
1048	35	502	4	0
1048	36	502	4	0
1048	34	496	0	0
1048	33	496	0	0
1048	38	508	2	0
1046	33	507	0	0
1048	38	509	2	0
1048	36	496	0	0
1048	38	507	2	0
1048	38	503	2	0
1046	33	505	0	0
1048	38	514	2	0
1048	38	519	2	0
1048	38	512	2	0
1048	32	504	4	0
1046	33	519	0	0
1048	38	513	2	0
1048	38	499	4	0
1048	37	499	4	0
1048	34	502	4	0
1048	32	499	4	0
1048	32	496	0	0
1048	38	511	2	0
1048	38	505	2	0
1046	33	513	0	0
1048	38	517	2	0
1048	38	518	2	0
1046	33	509	0	0
1048	38	515	2	0
1048	38	516	2	0
1048	38	504	2	0
1046	33	517	0	0
1048	35	496	0	0
1048	38	506	2	0
1048	34	499	4	0
1048	37	496	0	0
1048	38	496	0	0
1048	39	499	4	0
1048	32	502	4	0
1046	33	511	0	0
1048	35	499	4	0
1048	43	518	2	0
1048	43	522	2	0
1048	43	524	2	0
1048	43	513	2	0
1048	43	514	2	0
1048	43	516	2	0
1048	41	518	2	0
1048	43	512	2	0
1048	41	517	2	0
1048	43	515	2	0
1048	41	512	2	0
1048	43	523	2	0
1048	41	516	2	0
1048	41	519	2	0
34	47	512	2	0
1048	43	525	2	0
34	46	513	3	0
1048	41	513	2	0
1048	43	519	2	0
1048	43	517	2	0
1048	41	524	2	0
1048	43	521	2	0
1048	41	514	2	0
1048	41	515	2	0
1048	41	527	2	0
1048	41	522	2	0
1048	41	521	2	0
1048	41	520	2	0
1048	41	523	2	0
1048	41	526	2	0
1048	43	520	2	0
1048	41	525	2	0
34	45	504	7	0
1048	41	504	2	0
1048	41	511	2	0
1048	41	509	2	0
1048	41	508	2	0
1048	43	510	2	0
1048	41	510	2	0
1048	43	508	2	0
1048	41	506	2	0
1048	41	505	2	0
1	47	508	3	0
1048	43	507	2	0
1048	43	506	2	0
1072	47	511	0	0
1072	46	507	0	0
1048	43	504	2	0
1048	43	505	2	0
34	47	505	6	0
1048	41	507	2	0
34	47	506	3	0
1048	43	511	2	0
1048	43	509	2	0
1072	45	510	0	0
1048	43	501	2	0
1048	43	502	2	0
1048	43	500	2	0
1048	43	499	2	0
34	47	502	0	0
1048	43	498	2	0
1048	41	500	2	0
1048	43	503	2	0
1048	40	496	0	0
1048	41	496	0	0
34	40	500	3	0
1048	43	496	2	0
1048	42	496	0	0
1048	43	497	2	0
1048	41	501	2	0
1048	41	503	2	0
1048	40	499	4	0
1048	41	502	2	0
1042	35	492	3	0
1042	38	493	2	0
38	52	534	0	0
1	53	533	0	0
283	54	534	0	0
283	51	533	0	0
34	48	529	0	0
283	53	532	0	0
1	52	524	1	0
1	61	525	1	0
1	58	522	1	0
34	62	533	0	0
401	63	529	0	0
402	61	535	0	0
1	62	515	1	0
34	61	516	1	0
37	56	513	1	0
1	56	517	1	0
2	53	513	1	0
1	51	517	1	0
34	48	514	3	0
1	55	509	1	0
1	59	509	1	0
283	51	507	5	0
34	48	506	2	0
283	49	505	3	0
283	51	505	5	0
1	50	511	1	0
34	49	507	1	0
34	48	511	7	0
34	49	510	6	0
622	56	504	4	0
283	48	509	5	0
34	50	503	5	0
34	53	497	1	0
1	61	503	1	0
38	55	492	1	0
1	53	494	1	0
1	51	491	1	0
1	61	488	1	0
1	53	495	1	0
1	58	491	1	0
306	40	485	2	0
1	54	486	1	0
1	61	484	1	0
34	37	483	0	0
34	36	484	0	0
1056	35	486	5	0
34	39	487	2	0
1	39	485	7	0
1072	35	483	0	0
38	56	482	1	0
34	39	483	1	0
0	38	486	1	0
34	41	485	4	0
34	37	486	1	0
0	37	485	1	0
1072	42	487	0	0
34	40	486	6	0
34	59	482	1	0
34	34	485	4	0
1	56	474	1	0
1	51	474	1	0
1	61	472	1	0
0	63	474	0	0
1	53	470	1	0
1	51	471	1	0
1	58	470	1	0
0	60	467	0	0
1	57	464	1	0
1	70	467	1	0
1	67	466	0	0
38	68	492	1	0
38	67	489	1	0
0	67	473	0	0
0	64	470	0	0
0	70	474	0	0
1	70	484	1	0
1	50	459	1	0
0	54	457	0	0
1	55	456	0	0
1	55	460	0	0
1	48	457	1	0
1	58	458	1	0
0	65	458	0	0
1	67	460	0	0
1	68	456	0	0
1	69	463	0	0
1	64	461	0	0
0	71	462	0	0
0	70	457	0	0
0	59	459	0	0
1	61	462	0	0
1	59	456	0	0
0	68	462	0	0
1	69	459	0	0
1	68	463	0	0
0	68	461	0	0
1	61	456	0	0
1	64	458	0	0
4	62	449	0	0
0	62	452	0	0
0	71	454	0	0
38	56	450	0	0
0	69	454	0	0
0	60	455	0	0
34	58	453	0	0
1041	61	448	0	0
34	60	449	0	0
0	69	451	0	0
0	68	453	0	0
0	64	454	0	0
0	66	448	0	0
1	64	455	0	0
1	71	450	0	0
0	52	441	0	0
1041	63	442	0	0
1041	63	446	0	0
8	55	440	0	0
34	55	447	0	0
1041	63	441	0	0
1041	63	445	0	0
1039	59	445	0	0
1039	59	443	0	0
21	62	445	0	0
21	61	442	0	0
1039	59	442	0	0
4	58	445	0	0
1039	64	442	0	0
8	67	446	6	0
1039	65	445	0	0
21	65	444	0	0
4	68	443	6	0
0	58	442	0	0
8	66	447	0	0
1039	67	440	0	0
8	56	441	1	0
1039	60	445	0	0
21	63	440	0	0
21	59	441	0	0
1039	58	443	0	0
1	69	443	0	0
1040	64	445	0	0
1041	65	446	0	0
1039	68	440	0	0
8	66	446	0	0
21	61	447	0	0
1039	60	446	0	0
1041	64	441	0	0
1	75	441	0	0
1	72	441	0	0
1	73	441	0	0
1	73	445	0	0
0	75	447	0	0
7	78	444	6	0
7	79	443	4	0
1	72	444	0	0
1	74	444	0	0
3	79	444	2	0
1	76	443	0	0
0	76	440	4	0
1	76	447	0	0
41	78	448	0	0
1	75	452	0	0
1	74	454	0	0
1	73	454	0	0
34	79	471	4	0
1	74	461	0	0
1	72	447	0	0
205	73	469	0	0
283	77	474	7	0
34	79	472	1	0
20	73	459	0	0
0	78	461	0	0
1	78	457	0	0
1	78	463	0	0
1	73	449	0	0
283	78	473	6	0
283	79	474	5	0
34	78	474	2	0
1	75	455	0	0
0	74	458	0	0
1	77	450	0	0
3	79	453	0	0
1	74	468	1	0
7	78	453	6	0
34	78	472	1	0
283	79	473	2	0
1	72	483	1	0
20	75	481	0	0
1	72	494	1	0
1	79	493	1	0
34	76	500	1	0
34	76	507	1	0
1	78	517	1	0
1	76	513	1	0
1	65	512	1	0
37	71	512	1	0
1	71	518	1	0
1	74	524	1	0
34	66	520	1	0
37	70	523	1	0
1	74	522	1	0
1	80	515	1	0
34	81	508	1	0
34	84	504	1	0
1	84	489	1	0
34	87	494	1	0
511	86	524	0	0
1	85	494	1	0
1	82	493	1	0
1088	80	505	0	0
3	85	521	0	0
5	82	525	0	0
1	86	513	1	1
1	86	522	1	1
1	86	525	1	1
47	85	532	6	0
71	85	535	2	0
34	65	533	0	0
22	85	531	6	0
22	86	531	6	0
23	80	531	0	0
14	79	532	4	0
402	64	531	0	0
34	64	530	0	0
3	81	532	4	0
402	65	530	0	0
1	83	532	0	1
145	83	534	1	1
1	87	533	1	1
1	81	529	0	1
25	84	531	4	0
25	82	532	6	0
3	83	535	6	0
24	84	535	4	0
34	67	538	7	0
100	70	543	4	0
34	65	541	7	0
0	76	539	7	0
34	82	540	7	0
102	76	543	4	0
102	75	543	4	0
34	82	543	7	0
34	75	541	7	0
34	84	538	7	0
4	64	541	7	0
0	87	543	7	0
34	73	539	7	0
0	84	543	7	0
0	79	539	7	0
4	81	543	7	0
4	65	539	7	0
402	65	536	0	0
1	79	542	7	0
105	78	546	0	0
101	72	549	0	0
105	78	545	0	0
101	74	549	0	0
101	73	549	0	0
105	79	546	0	0
1	87	546	7	0
104	74	545	4	0
34	79	549	7	0
34	72	546	7	0
1	77	551	0	0
101	73	547	0	0
101	75	546	0	0
34	72	551	7	0
4	87	544	7	0
1	85	548	7	0
103	76	544	0	0
37	72	544	7	0
100	76	547	4	0
1	84	551	7	0
4	80	549	0	0
34	80	548	7	0
1	87	545	7	0
1	81	545	7	0
1	83	547	7	0
34	85	546	7	0
34	82	551	7	0
16	92	527	6	0
513	93	521	4	0
556	93	525	2	0
16	94	527	6	0
278	92	520	2	0
22	91	520	6	0
22	90	520	6	0
3	93	530	6	0
0	94	542	7	0
4	92	547	7	0
1	88	540	7	0
20	88	536	0	0
281	94	532	6	0
0	89	543	7	0
1	89	529	1	1
3	91	530	6	0
1	89	534	1	1
4	89	544	7	0
22	89	520	6	0
2	89	521	0	0
34	90	540	7	0
0	92	542	7	0
0	92	538	7	0
20	94	536	0	0
0	89	551	7	0
1	89	549	7	0
281	94	534	6	0
1	90	546	7	0
7	98	523	4	0
5	102	532	0	0
143	103	528	0	0
23	98	534	0	0
1	97	549	4	0
15	102	534	0	0
23	98	532	0	0
143	101	535	6	0
1	97	554	4	0
0	97	547	4	0
1	99	548	4	0
23	96	534	0	0
1	97	556	4	0
0	98	544	4	0
144	97	530	0	0
0	98	552	4	0
1	103	530	0	1
50	102	520	0	0
1	99	526	1	1
29	102	524	0	0
3	98	524	0	0
1	101	530	0	1
143	101	533	6	0
0	98	554	4	0
5	96	523	0	0
23	96	532	0	0
25	99	528	0	0
0	101	551	4	0
0	98	558	4	0
0	100	546	4	0
71	96	526	0	0
143	101	528	0	0
25	96	528	0	0
143	101	531	6	0
1	103	557	4	0
0	101	554	4	0
20	103	532	0	1
68	110	557	2	0
68	106	555	1	0
68	110	545	2	0
68	104	551	0	0
69	110	551	0	0
68	106	547	3	0
3	107	533	0	0
6	106	534	0	0
11	110	534	0	0
15	111	529	0	0
1	106	532	0	1
1	111	532	0	1
1	108	530	0	1
1	108	535	1	1
1	110	533	1	1
3	111	527	0	0
7	111	526	4	0
15	109	521	4	0
3	106	526	0	0
3	106	522	4	0
1	104	525	1	1
1	109	523	1	1
1	109	520	1	1
1	105	524	0	1
3	116	523	0	0
11	113	521	0	0
3	119	522	0	0
3	116	522	0	0
15	116	534	0	0
5	119	526	0	0
3	113	520	0	0
68	114	547	1	0
7	115	529	4	0
7	114	535	0	0
3	114	534	0	0
68	114	555	3	0
1	118	558	7	0
11	116	525	0	0
6	118	520	0	0
3	119	523	0	0
5	113	533	0	0
15	112	534	0	0
68	116	551	0	0
3	115	530	0	0
1	119	525	0	1
1	113	531	0	1
1	116	524	1	1
1	113	532	0	1
1	115	532	0	1
1	115	526	1	1
2	114	523	1	1
0	118	562	7	0
1	118	563	7	0
34	117	565	1	0
1	119	565	7	0
1	118	568	0	0
0	116	570	0	0
1	118	572	0	0
0	117	573	0	0
20	116	575	0	0
0	113	573	0	0
20	124	575	0	0
34	121	563	1	0
0	124	572	0	0
1	125	569	0	0
1	122	559	0	0
0	126	571	0	0
0	120	559	7	0
34	125	563	7	0
1	121	556	0	0
0	124	555	7	0
1	127	543	0	0
34	124	557	1	0
1	120	569	0	0
1	122	568	0	0
20	120	575	0	0
0	120	573	0	0
1	121	570	0	0
34	121	566	7	0
1	122	572	0	0
59	127	558	4	0
34	133	542	0	0
72	135	559	3	0
72	134	558	3	0
72	135	558	3	0
20	129	536	0	0
34	129	540	0	0
72	134	559	3	0
72	134	565	7	0
72	134	567	3	0
72	131	559	3	0
20	133	536	0	0
72	135	560	3	0
72	134	566	7	0
72	130	558	3	0
66	132	536	0	0
72	131	558	3	0
72	130	559	3	0
72	133	567	3	0
72	131	564	3	0
72	133	560	7	0
72	131	567	3	0
72	133	563	3	0
72	130	565	3	0
72	134	560	3	0
72	131	565	3	0
72	130	567	3	0
72	132	567	7	0
72	130	566	7	0
72	132	565	3	0
34	128	542	0	0
72	132	560	3	0
34	131	544	0	0
72	134	561	3	0
72	135	562	3	0
72	134	562	3	0
72	134	564	3	0
72	133	559	3	0
72	135	561	3	0
72	135	564	7	0
72	134	563	7	0
72	133	558	3	0
72	132	559	3	0
72	132	558	3	0
72	135	565	3	0
72	130	564	3	0
72	135	563	3	0
72	135	567	3	0
72	130	568	3	0
34	128	552	1	0
72	135	566	3	0
72	130	561	3	0
72	130	569	3	0
72	132	570	3	0
34	128	561	1	0
72	132	566	3	0
72	133	570	3	0
72	130	570	3	0
1	128	573	0	0
72	129	565	3	0
72	135	568	3	0
20	129	575	0	0
72	130	562	3	0
72	134	569	3	0
72	130	563	3	0
72	133	569	3	0
72	129	566	3	0
72	134	568	3	0
72	131	566	3	0
72	133	568	3	0
72	129	567	3	0
72	132	571	7	0
72	133	562	3	0
72	131	570	7	0
72	131	561	3	0
72	132	564	3	0
72	131	563	3	0
72	129	564	3	0
72	131	562	3	0
72	129	568	3	0
72	132	562	3	0
72	131	569	3	0
72	131	568	7	0
72	132	568	3	0
72	133	564	7	0
72	129	563	3	0
72	131	560	7	0
72	132	569	3	0
72	132	561	3	0
72	129	569	3	0
72	130	571	3	0
72	129	570	3	0
72	133	565	3	0
72	130	560	3	0
20	130	573	0	0
72	133	566	7	0
72	132	563	3	0
72	131	571	3	0
72	133	561	3	0
9	134	538	0	1
0	128	531	0	0
7	124	534	0	0
3	126	534	0	0
11	121	533	0	0
5	121	530	0	0
1	124	532	1	1
1	121	532	1	1
1	135	533	3	1
1	127	532	1	1
73	127	525	0	0
3	124	525	0	0
7	124	524	4	0
3	123	522	0	0
92	132	524	4	0
20	130	525	0	0
3	125	521	0	0
20	130	523	0	0
22	120	520	0	0
7	127	523	0	0
16	121	520	0	0
7	127	521	4	0
3	127	522	0	0
3	122	526	0	0
24	122	520	0	0
3	129	526	0	0
29	121	521	2	0
16	121	523	0	0
7	125	522	0	0
7	129	525	4	0
1	133	525	1	1
1	127	524	1	1
72	137	559	3	0
29	136	525	2	0
5	136	530	0	0
3	138	532	1	0
5	138	522	0	0
72	136	558	3	0
72	136	559	3	0
1	141	533	1	1
1	143	522	0	1
1	142	524	0	1
7	145	526	6	0
41	145	529	4	0
145	148	555	6	0
145	148	557	6	0
3	146	526	4	0
41	151	558	0	0
0	149	552	4	0
0	150	540	4	0
145	148	559	6	0
1	148	558	1	1
21	148	533	0	1
44	150	554	0	1
1	145	533	1	1
72	141	563	3	0
72	141	565	3	0
72	140	565	3	0
72	140	566	3	0
72	141	564	3	0
72	139	566	3	0
72	139	565	3	0
72	139	564	7	0
72	138	565	3	0
72	140	564	7	0
72	140	563	3	0
72	138	561	3	0
72	138	563	3	0
72	138	562	3	0
72	139	563	3	0
72	136	560	3	0
72	137	561	7	0
72	136	567	3	0
72	137	562	3	0
72	136	563	7	0
72	137	560	3	0
72	136	562	3	0
72	136	561	3	0
72	140	562	3	0
72	138	566	3	0
72	139	561	3	0
72	139	562	3	0
72	138	564	3	0
72	137	566	3	0
20	147	563	4	0
72	136	566	3	0
72	136	565	7	0
72	138	560	3	0
72	137	563	3	0
72	137	564	7	0
72	136	564	3	0
72	137	565	3	0
1	150	561	0	1
34	146	574	6	0
20	115	582	0	0
20	113	577	0	0
20	112	580	6	0
145	152	557	2	0
0	158	535	4	0
0	156	549	4	0
20	152	563	4	0
36	159	560	0	0
34	158	564	6	0
38	155	544	4	0
36	159	559	0	0
34	159	548	4	0
145	152	555	2	0
34	156	556	4	0
11	152	561	0	0
4	157	522	0	0
1	157	526	0	0
1	159	521	0	0
286	159	523	3	0
1	167	520	0	0
1	165	525	0	0
0	166	521	0	0
0	166	524	0	0
0	164	522	0	0
1	162	521	0	0
4	162	524	0	0
105	163	534	0	0
0	167	541	4	0
115	161	534	0	0
105	161	535	0	0
1	161	523	0	0
99	164	532	0	0
105	163	535	0	0
0	165	531	4	0
99	160	542	0	0
99	162	540	0	0
196	163	546	0	0
99	160	538	0	0
36	160	559	0	0
103	161	538	0	0
196	164	547	0	0
80	162	556	4	0
0	166	540	4	0
115	162	541	0	0
80	161	556	4	0
105	167	544	0	0
103	165	546	0	0
105	167	545	0	0
105	163	536	0	0
105	164	544	0	0
196	163	543	0	0
34	161	560	0	0
34	163	560	0	0
34	165	564	6	0
34	155	568	6	0
0	165	571	6	0
34	157	568	6	0
1	157	571	6	0
0	155	582	0	0
0	170	537	4	0
0	174	551	4	0
0	170	564	6	0
105	168	544	0	0
0	172	567	6	0
34	175	570	6	0
34	169	565	6	0
34	170	556	4	0
34	173	544	4	0
34	169	548	4	0
34	173	564	6	0
20	174	563	4	0
42	151	1502	0	0
145	148	1501	6	0
45	150	1501	0	0
145	148	1499	6	0
45	150	1504	4	0
158	148	1507	6	0
29	150	1499	6	0
145	148	1505	6	0
45	152	1501	0	0
145	151	1507	0	0
45	148	1506	4	0
45	149	1504	6	0
145	148	1503	6	0
45	149	1506	4	0
45	149	1503	6	0
45	149	1501	6	0
45	149	1502	6	0
45	151	1501	0	0
0	177	540	4	0
88	183	563	0	0
70	180	561	0	0
70	180	567	6	0
70	182	565	6	0
70	182	575	6	0
70	182	572	6	0
70	180	570	6	0
20	177	573	0	0
20	177	569	0	0
0	182	531	4	0
70	188	546	0	0
70	188	556	6	0
88	190	553	0	0
70	187	566	6	0
70	185	550	0	0
70	188	542	0	0
0	187	530	4	0
70	191	550	6	0
70	184	540	0	0
88	190	545	0	0
70	191	559	6	0
88	187	561	0	0
70	190	563	6	0
70	190	544	0	0
70	190	565	6	0
70	184	559	6	0
0	171	521	0	0
0	171	525	0	0
286	169	523	3	0
4	169	520	0	0
34	174	522	0	0
1	189	521	0	0
0	168	521	0	0
1	185	525	0	0
1	167	518	0	0
34	186	519	0	0
7	168	513	3	0
3	168	514	0	0
1	172	519	0	0
0	163	519	0	0
37	161	516	4	0
1	171	516	0	0
0	171	513	4	0
8	161	518	0	0
283	159	512	3	0
15	163	516	0	0
0	170	514	0	0
286	158	513	3	0
0	158	515	5	0
1	157	519	0	0
283	160	512	3	0
1	173	513	0	0
7	167	514	0	0
283	159	516	3	0
0	174	518	0	0
286	160	514	3	0
15	167	516	0	0
1	167	515	1	1
3	164	514	0	0
1	163	513	1	1
1	165	515	1	1
1	166	515	0	1
20	155	505	4	0
7	153	504	2	0
20	155	510	4	0
132	156	506	2	0
37	166	508	4	0
37	162	508	4	0
37	164	508	4	0
4	160	509	0	0
283	160	511	3	0
4	170	511	0	0
1	188	510	0	0
34	185	508	0	0
1	164	511	0	1
50	145	514	0	0
50	145	512	0	0
63	150	507	6	0
29	148	504	2	0
139	148	507	6	0
50	148	513	0	0
1	149	512	1	1
1	146	510	0	1
5	136	517	1	0
1	136	514	2	1
283	88	516	0	0
29	89	516	0	0
278	89	515	0	0
278	90	515	0	0
278	91	516	0	0
1183	91	518	6	0
283	88	518	0	0
279	88	517	2	0
20	89	511	1	0
20	89	506	1	0
49	102	500	0	0
27	101	496	4	0
24	98	496	4	0
29	101	514	0	0
139	100	509	2	0
63	102	509	2	0
45	99	489	2	0
45	99	488	2	0
63	103	494	0	0
27	97	494	2	0
39	97	488	2	0
43	97	490	4	0
45	99	486	2	0
45	99	485	2	0
23	98	480	0	0
45	103	485	4	0
45	101	484	6	0
45	99	484	2	0
46	97	485	0	0
34	89	484	1	0
23	101	480	0	0
45	99	487	2	0
18	97	484	4	0
46	97	486	0	0
45	101	485	6	0
45	102	485	4	0
1	83	483	1	0
1	81	485	1	0
23	98	478	0	0
23	101	476	0	0
25	102	473	0	0
19	99	474	4	0
0	88	472	0	0
34	80	473	0	0
23	101	478	0	0
25	97	478	0	0
34	82	472	5	0
283	81	472	2	0
1	80	472	7	0
25	98	473	0	0
23	98	476	0	0
0	84	474	0	0
1	103	474	1	1
283	82	471	0	0
4	98	464	0	0
4	93	467	0	0
1	88	466	0	0
1	82	466	0	0
4	95	464	0	0
34	103	466	0	0
1	85	464	0	0
4	102	465	0	0
38	100	467	0	0
283	80	470	3	0
34	81	470	1	0
283	81	471	3	0
20	83	471	0	0
34	80	471	3	0
4	101	466	0	0
1	87	463	0	0
1	86	459	0	0
4	100	463	0	0
4	100	461	0	0
34	101	459	0	0
1	89	458	0	0
4	94	458	0	0
4	99	457	0	0
4	92	463	0	0
1	85	456	0	0
4	94	463	0	0
1	82	461	0	0
0	82	459	0	0
4	102	460	0	0
34	103	463	0	0
34	97	462	0	0
8	103	459	2	0
34	98	459	0	0
8	103	458	1	0
4	96	461	0	0
16	85	452	2	0
34	102	450	0	0
37	101	453	0	0
7	84	448	0	0
1	88	450	4	0
1	102	452	0	0
1	92	450	4	0
0	94	448	4	0
1	94	451	4	0
1	98	448	0	0
0	91	450	4	0
37	99	448	0	0
34	91	448	4	0
16	82	452	2	0
1	90	452	4	0
1	89	455	0	0
71	83	454	2	0
29	83	452	0	0
34	95	449	4	0
7	85	446	2	0
7	84	445	4	0
63	81	441	2	0
7	82	446	6	0
7	82	447	6	0
3	84	447	2	0
3	84	446	2	0
3	83	446	2	0
3	83	447	2	0
7	84	442	4	0
3	84	443	2	0
1	87	443	4	0
76	80	441	2	0
1	92	446	4	0
0	88	447	4	0
1	91	447	4	0
34	93	441	4	0
1	92	442	4	0
0	90	442	4	0
0	102	444	0	0
1	101	442	0	0
0	88	440	4	0
34	99	447	0	0
0	98	443	0	0
37	90	447	4	0
34	89	444	4	0
1	99	441	0	0
0	89	445	4	0
1	94	445	4	0
34	93	443	4	0
37	94	446	4	0
1	94	440	4	0
1039	68	439	0	0
8	69	437	0	0
209	71	432	6	0
209	70	432	4	0
0	70	434	0	0
21	68	437	0	0
34	81	432	4	0
37	80	434	4	0
1	80	433	4	0
37	94	433	4	0
0	87	438	4	0
1041	66	439	0	0
0	77	435	4	0
1041	67	436	0	0
1041	66	436	3	0
1	90	433	4	0
0	76	437	4	0
0	76	433	4	0
209	78	432	2	0
21	64	437	0	0
4	68	435	3	0
1	75	438	0	0
209	86	433	7	0
1	72	435	0	0
0	77	438	4	0
34	78	434	4	0
34	87	433	4	0
1041	65	438	0	0
8	66	435	0	0
209	76	432	0	0
1041	65	439	0	0
209	82	433	3	0
1	78	435	4	0
0	78	436	4	0
0	79	432	4	0
209	83	433	6	0
0	75	435	4	0
209	73	432	6	0
209	74	432	1	0
1	73	439	0	0
0	79	436	4	0
20	84	439	0	0
20	81	439	0	0
1	89	439	4	0
0	88	432	4	0
34	93	433	4	0
34	89	433	4	0
1	92	439	4	0
1	102	434	0	0
0	92	432	4	0
1	97	432	0	0
70	90	428	0	0
4	89	426	0	0
70	92	425	0	0
70	64	430	0	0
70	90	427	0	0
208	77	426	0	0
205	73	428	0	0
70	72	424	0	0
70	64	429	0	0
70	68	426	0	0
70	84	426	0	0
70	92	428	0	0
205	85	431	0	0
70	81	426	0	0
70	81	422	0	0
70	68	420	0	0
70	64	417	0	0
70	82	423	0	0
70	82	417	0	0
4	89	420	0	0
70	89	422	0	0
38	92	420	0	0
70	89	417	0	0
70	80	423	0	0
70	85	418	0	0
70	97	421	0	0
38	97	417	0	0
70	99	420	0	0
70	98	417	0	0
4	63	426	0	0
1039	61	436	0	0
1039	63	436	4	0
8	57	436	0	0
1039	62	438	0	0
38	58	439	0	0
1039	62	437	0	0
0	61	435	0	0
205	58	426	0	0
4	56	434	4	0
21	63	437	0	0
1039	58	438	0	0
1041	58	437	0	0
1039	59	438	0	0
21	59	439	0	0
1041	60	437	0	0
1041	59	437	0	0
1041	61	438	0	0
8	61	439	0	0
70	61	424	0	0
1041	60	436	0	0
38	56	430	0	0
70	62	429	0	0
38	56	424	0	0
1039	62	436	0	0
8	58	435	0	0
5	60	439	0	0
5	60	1384	0	0
6	60	1383	0	0
45	79	1391	4	0
45	78	1391	4	0
15	78	1386	4	0
42	78	1392	0	0
15	84	1391	2	0
15	82	1386	6	0
3	85	1393	0	0
3	82	1388	4	0
2	82	1393	1	1
3	80	1386	0	0
3	83	1395	4	0
1	84	1389	0	1
1	82	1397	1	1
15	84	1395	4	0
1	79	1389	0	1
45	80	1392	2	0
45	80	1394	2	0
45	80	1393	2	0
4	57	422	0	0
205	63	420	0	0
70	57	417	0	0
4	58	421	0	0
70	59	408	0	0
205	56	413	0	0
70	95	412	0	0
70	86	412	0	0
70	90	412	0	0
70	56	414	0	0
38	72	409	0	0
70	82	409	0	0
4	82	412	0	0
4	83	411	0	0
70	94	415	0	0
70	82	410	0	0
70	63	409	0	0
70	55	415	0	0
70	52	418	0	0
205	54	418	0	0
4	51	413	0	0
4	54	415	0	0
70	48	422	0	0
38	49	425	0	0
70	54	428	0	0
70	50	423	0	0
0	52	436	0	0
4	55	437	2	0
199	51	438	1	1
38	50	420	0	0
37	107	463	0	0
1	107	433	0	0
20	108	465	0	0
38	104	457	0	0
34	104	444	0	0
34	106	440	0	0
34	106	450	0	0
4	105	465	0	0
8	107	460	3	0
37	104	446	0	0
8	107	459	0	0
4	105	459	0	0
1	107	445	0	0
0	107	448	0	0
34	104	448	0	0
38	106	442	0	0
4	107	462	0	0
4	104	462	0	0
0	105	443	0	0
70	107	428	0	0
70	108	416	0	0
70	109	423	0	0
70	108	418	0	0
70	105	421	0	0
70	105	419	0	0
4	108	419	0	0
38	111	417	0	0
71	104	476	2	0
63	104	478	0	0
5	104	473	0	0
3	106	473	2	0
79	111	474	0	0
1	107	474	1	1
63	107	480	0	0
45	104	487	6	0
3	105	482	0	0
45	104	485	4	0
45	104	486	6	0
3	104	482	0	0
77	117	461	4	0
6	119	450	0	0
1	115	448	0	0
1	112	458	0	0
1	112	456	0	0
3	118	460	0	0
9	119	475	0	0
1	112	464	1	1
11	118	465	0	0
7	119	470	4	0
7	119	476	0	0
48	118	461	0	0
7	118	459	4	0
15	118	456	6	0
1	118	467	1	1
1	119	461	0	1
1	117	484	6	0
20	114	479	0	0
1	119	480	6	0
1	119	484	6	0
34	114	449	0	0
11	118	463	0	0
9	119	471	0	0
7	119	474	4	0
7	119	472	0	0
1	116	475	0	0
1	117	480	6	0
20	114	485	0	0
1	112	450	0	0
1	115	484	6	0
1	115	472	0	0
1	115	480	6	0
1	118	444	0	0
34	117	447	0	0
37	117	440	0	0
34	117	445	0	0
1	118	441	0	0
38	115	440	0	0
34	114	444	0	0
34	114	441	0	0
34	116	440	0	0
1	116	443	0	0
1	118	438	0	0
1	118	435	0	0
20	115	439	0	0
20	115	436	0	0
34	117	436	0	0
0	118	433	0	0
20	115	433	0	0
20	115	426	0	0
20	115	430	0	0
4	112	425	0	0
20	118	424	0	0
38	119	419	0	0
70	116	418	0	0
208	116	410	2	0
70	102	408	0	0
4	101	412	0	0
70	107	413	0	0
70	99	414	0	0
70	99	415	0	0
70	106	411	0	0
70	105	410	0	0
70	97	410	0	0
38	104	413	0	0
70	97	411	0	0
70	97	412	0	0
70	96	409	0	0
1	122	447	0	0
1	120	440	0	0
70	125	409	0	0
1	121	443	0	0
1	125	442	0	0
1	126	446	0	0
1	127	434	0	0
0	127	438	0	0
70	120	426	0	0
70	126	416	0	0
70	125	419	0	0
20	122	423	0	0
4	122	421	0	0
70	126	408	0	0
70	121	429	0	0
205	124	410	0	0
220	120	410	0	0
34	121	439	0	0
1	120	445	0	0
38	127	436	0	0
34	122	434	0	0
70	122	419	0	0
1	126	440	0	0
70	126	424	0	0
70	125	431	0	0
20	127	423	0	0
0	126	436	0	0
1	124	433	0	0
70	123	408	0	0
1	122	437	0	0
47	126	455	0	0
1	121	452	0	0
1	124	455	0	1
7	122	459	5	0
5	123	463	0	0
47	126	458	0	0
3	120	463	0	0
3	121	463	0	0
15	122	456	4	0
3	123	459	0	0
1	122	461	0	1
33	127	467	7	0
34	127	468	7	0
37	127	465	0	0
37	126	464	0	0
7	122	470	4	0
3	121	465	0	0
3	120	464	0	0
7	120	470	4	0
3	120	465	0	0
3	121	464	6	0
7	121	470	4	0
1	124	466	1	1
1	121	469	0	1
1	124	470	1	1
10	127	476	2	0
46	127	479	4	0
10	126	476	2	0
7	122	474	4	0
7	120	474	4	0
7	120	472	0	0
145	124	472	6	0
7	122	476	0	0
7	122	472	0	0
145	124	474	6	0
145	124	476	6	0
7	121	476	0	0
46	124	479	4	0
7	121	474	4	0
7	121	472	0	0
7	120	476	0	0
1	126	472	0	1
1	121	484	6	0
1	121	480	6	0
1	135	448	0	0
47	133	455	0	0
67	132	455	4	0
34	128	464	7	0
47	129	455	4	0
47	133	458	0	0
47	130	455	0	0
46	133	479	4	0
145	128	474	2	0
25	134	475	0	0
34	131	465	7	0
37	128	466	0	0
25	134	477	0	0
37	130	465	0	0
25	129	475	0	0
34	129	467	7	0
25	129	473	0	0
145	128	476	2	0
20	130	487	0	0
1	134	451	0	0
25	129	477	0	0
25	134	473	0	0
46	130	479	4	0
145	128	472	2	0
47	129	458	4	0
43	131	470	4	0
63	131	478	6	0
47	130	458	0	0
20	133	487	0	0
47	132	458	4	0
1	130	469	0	1
1	130	463	0	1
1	128	460	0	1
1	131	438	0	0
0	134	439	0	0
1	128	442	0	0
0	130	435	0	0
0	133	434	0	0
34	132	440	0	0
37	131	434	0	0
1	131	442	0	0
34	129	439	0	0
34	130	433	0	0
1	134	446	0	0
0	135	435	0	0
0	129	438	0	0
4	129	424	0	0
20	133	420	0	0
38	132	418	0	0
20	130	420	0	0
70	134	415	0	0
70	133	415	0	0
70	133	414	0	0
205	131	414	0	0
20	137	420	0	0
4	141	413	0	0
70	141	411	0	0
38	138	415	0	0
0	142	438	0	0
38	143	412	0	0
70	138	421	0	0
70	141	417	0	0
0	140	436	0	0
0	142	433	0	0
70	143	420	0	0
1	138	446	0	0
1	139	434	0	0
1	141	440	0	0
70	138	431	0	0
37	138	436	0	0
0	137	433	0	0
70	138	417	0	0
0	142	441	0	0
0	139	442	0	0
0	138	438	0	0
38	136	426	0	0
0	137	441	0	0
4	137	431	0	0
38	138	439	0	0
1	141	448	0	0
5	142	454	0	0
27	138	462	2	0
14	140	460	4	0
1	139	456	3	1
70	151	417	0	0
205	147	418	0	0
4	151	419	0	0
1	147	437	2	0
1	145	434	2	0
0	144	435	0	0
1	151	445	0	0
70	149	427	0	0
70	145	419	0	0
70	148	430	0	0
70	145	427	0	0
1	151	448	0	0
20	144	443	2	0
1	145	441	2	0
20	150	452	0	0
20	148	449	2	0
1	145	450	2	0
20	146	445	2	0
1	146	447	0	0
1	151	441	0	0
1	148	441	2	0
1	148	454	2	0
70	150	410	0	0
70	147	411	0	0
70	154	413	0	0
70	155	410	0	0
70	158	428	0	0
70	154	408	0	0
1	154	434	0	0
70	159	427	0	0
1	159	439	0	0
0	152	439	0	0
1	154	445	0	0
4	157	424	0	0
1	158	440	0	0
4	157	425	0	0
1	155	443	0	0
0	154	438	0	0
1	159	435	0	0
1	159	444	0	0
70	153	413	0	0
1	152	437	0	0
70	158	431	0	0
1	156	436	0	0
0	157	443	0	0
70	157	430	0	0
1	156	434	0	0
1	155	448	0	0
1	156	451	0	0
1	155	453	2	0
1	158	451	0	0
1	157	455	0	0
20	150	456	2	0
37	149	460	2	0
34	150	460	2	0
37	150	462	2	0
20	151	458	0	0
205	157	461	0	0
1	158	458	0	0
20	158	463	6	0
1	145	463	2	0
1	145	456	2	0
20	151	461	0	0
1	154	461	0	0
70	161	428	0	0
1	163	443	0	0
1	163	450	0	0
1	167	453	0	0
1	161	433	0	0
1	166	441	0	0
391	160	453	0	0
1	160	450	0	0
1	162	456	0	0
20	161	463	0	0
1	160	459	0	0
1	162	457	0	0
70	161	424	0	0
1	163	437	0	0
70	164	430	0	0
1	166	456	0	0
205	164	461	0	0
1	165	460	0	0
1	163	460	0	0
1	165	435	0	0
1	161	446	0	0
1	172	458	0	0
0	175	438	0	0
1	170	445	0	0
1	174	450	0	0
205	172	430	0	0
1	174	452	0	0
38	168	428	0	0
1	172	456	0	0
1	168	449	0	0
1	169	442	0	0
1	170	462	0	0
1	169	434	0	0
1	170	436	0	0
1	174	440	0	0
1	175	447	0	0
205	175	430	0	0
70	175	425	0	0
1	171	452	0	0
0	175	459	0	0
70	174	419	0	0
70	170	420	0	0
70	162	421	0	0
4	181	421	0	0
70	181	425	0	0
70	182	423	0	0
70	183	420	0	0
4	180	430	0	0
205	178	418	0	0
1	176	439	0	0
1	178	436	0	0
1	183	435	0	0
1	180	441	0	0
1	177	443	0	0
1	179	446	0	0
0	183	448	0	0
1	182	448	0	0
1	176	434	0	0
70	182	430	0	0
70	176	420	0	0
70	181	427	0	0
70	177	417	0	0
1	182	454	0	0
1	182	452	0	0
1	183	438	0	0
0	180	450	0	0
1	178	451	0	0
1	182	463	0	0
1	176	456	0	0
1	176	461	0	0
0	177	457	0	0
0	179	462	0	0
1	179	456	0	0
0	182	460	0	0
205	154	465	0	0
1	167	470	0	0
1	167	467	0	0
205	165	465	0	0
15	158	465	4	0
205	157	464	0	0
143	164	467	6	0
20	151	470	0	0
205	162	464	0	0
1	157	471	0	0
20	151	464	0	0
143	163	465	0	0
34	149	464	2	0
20	151	467	0	0
34	149	467	2	0
37	150	465	2	0
20	150	471	0	0
37	150	468	2	0
1	178	466	0	0
1	179	468	0	0
1	177	468	0	0
1	181	471	0	0
143	161	469	4	0
146	160	470	6	0
1	154	470	0	0
143	158	467	2	0
1	144	468	2	0
1	183	467	0	0
1	174	466	0	0
1	170	465	0	0
1	168	470	0	0
109	161	465	0	1
146	160	469	2	0
24	139	465	4	0
14	140	466	4	0
7	138	471	0	0
25	141	465	0	0
3	138	470	4	0
145	138	467	6	0
1	138	464	1	1
1	138	468	1	1
17	141	471	0	0
20	150	479	0	0
46	139	479	4	0
46	136	479	4	0
20	150	476	0	0
20	150	473	0	0
26	147	476	0	0
1	168	474	0	0
20	162	476	0	0
1	160	474	0	0
20	166	476	0	0
1	165	472	0	0
1	173	472	0	0
26	147	473	0	0
1	175	472	0	0
20	174	476	0	0
20	158	476	0	0
1	153	474	0	0
20	170	476	0	0
20	154	476	0	0
1	137	472	0	1
20	177	476	0	0
20	180	476	0	0
1	180	472	0	0
20	182	474	0	0
3	176	480	0	0
5	180	486	0	0
74	182	482	0	0
48	176	483	0	0
20	175	480	0	0
53	179	481	0	0
20	175	487	2	0
20	150	486	0	0
20	150	482	0	0
34	181	495	0	0
7	171	494	4	0
24	172	494	0	0
191	160	490	0	0
37	162	492	0	0
1	160	492	4	0
191	161	490	0	0
8	167	492	0	0
20	149	489	4	0
47	174	494	4	0
20	147	493	0	0
37	163	491	0	0
1	157	495	4	0
1	166	495	2	1
43	179	488	0	1
1	184	449	0	0
1	189	452	0	0
1	189	454	0	0
0	189	449	0	0
1	185	452	0	0
1	189	461	0	0
1	187	463	0	0
0	186	459	0	0
1	189	459	0	0
20	191	470	0	0
20	189	472	0	0
0	185	466	0	0
1	185	468	0	0
1	185	460	0	0
1	185	456	0	0
1	191	465	0	0
1	188	470	0	0
20	184	472	0	0
20	186	472	0	0
1	187	466	0	0
1	188	446	0	0
34	196	443	0	0
0	187	443	0	0
1	185	445	0	0
1	187	442	0	0
0	199	445	0	0
34	195	441	0	0
1	199	454	0	0
1	199	453	0	0
1	194	449	0	0
0	195	460	0	0
0	195	454	0	0
1	198	455	0	0
1	196	446	0	0
0	197	467	0	0
0	194	462	0	0
0	196	450	0	0
1	189	436	0	0
1	187	439	0	0
1	187	433	0	0
1	191	439	0	0
4	184	426	0	0
70	191	431	0	0
70	191	430	0	0
70	191	426	0	0
70	190	422	0	0
4	184	423	0	0
70	195	423	0	0
70	194	417	0	0
70	198	425	0	0
70	195	429	0	0
70	197	431	0	0
70	194	427	0	0
70	193	427	0	0
70	193	431	0	0
45	197	434	0	0
45	195	436	4	0
45	194	436	4	0
45	194	434	0	0
45	196	436	4	0
45	197	436	4	0
70	192	430	0	0
45	195	434	0	0
45	196	434	0	0
1	197	485	0	0
0	199	487	0	0
0	196	483	0	0
1	197	482	0	0
0	206	448	0	0
34	204	455	0	0
1	200	454	0	0
0	201	451	0	0
34	204	466	0	0
0	205	487	0	0
0	207	485	0	0
6	203	482	4	0
1	205	461	0	0
1	206	468	0	0
34	205	459	0	0
34	205	470	0	0
0	203	466	0	0
1	207	482	0	0
1	203	458	0	0
1	201	461	0	0
34	207	461	0	0
34	206	464	0	0
34	203	468	0	0
1	200	466	0	0
34	207	487	0	0
23	202	485	0	1
0	204	442	0	0
58	204	435	4	0
34	206	438	0	0
4	207	430	0	0
209	206	425	7	0
209	203	425	0	0
4	206	424	0	0
209	200	426	0	0
209	202	426	1	0
209	201	429	0	0
38	205	419	0	0
70	204	423	0	0
70	206	421	0	0
209	215	425	3	0
209	212	429	1	0
209	211	426	0	0
38	211	425	0	0
4	213	418	0	0
209	214	426	6	0
34	214	436	4	0
34	209	437	0	0
34	209	435	0	0
34	208	442	4	0
34	214	447	4	0
34	209	454	0	0
29	213	449	2	0
63	107	490	0	0
45	104	488	6	0
29	105	492	0	0
210	456	3709	1	1
45	104	489	6	0
37	125	493	0	0
29	105	493	2	0
34	126	490	1	0
20	114	489	0	0
34	124	490	1	0
34	125	489	1	0
33	124	491	1	0
1	123	493	1	0
34	125	491	1	0
34	125	494	1	0
51	93	3300	0	0
51	100	3288	0	0
51	88	3288	2	0
51	90	3292	4	0
51	103	3288	0	0
51	109	3291	6	0
58	107	3297	6	0
5	111	3306	0	0
51	117	3298	4	0
51	116	3293	2	0
80	116	3303	4	0
51	105	3291	2	0
80	116	3295	6	0
51	107	3288	0	0
80	114	3303	4	0
51	113	3308	4	0
51	116	3308	4	0
51	110	3299	0	0
51	127	3292	4	0
80	122	3298	2	0
51	122	3294	6	0
20	121	3288	0	0
51	127	3290	0	0
117	134	3293	2	0
51	106	3285	6	0
5	119	3282	0	0
51	106	3282	6	0
20	122	3286	0	0
51	117	3285	2	0
20	124	3286	4	0
51	132	3285	0	0
51	123	3284	0	0
20	120	3286	0	0
58	128	3285	4	0
51	141	3286	6	0
116	137	3290	2	0
58	138	3294	4	0
51	138	3283	0	0
51	138	3288	4	0
24	137	3290	0	1
51	101	3284	4	0
51	90	3283	2	0
51	98	3275	0	0
51	90	3279	2	0
51	103	3274	0	0
51	112	3272	0	0
51	115	3272	0	0
51	94	3278	4	0
51	112	3279	4	0
58	108	3275	4	0
51	147	3292	0	0
51	155	3291	0	0
51	153	3301	4	0
51	166	3290	6	0
51	164	3287	0	0
51	161	3290	2	0
51	164	3303	6	0
51	160	3303	2	0
51	161	3296	6	0
51	161	3306	4	0
1	129	489	0	0
1	133	489	0	0
1	129	492	0	0
1	133	492	0	0
49	106	500	7	0
49	104	500	0	0
49	106	502	6	0
34	124	497	1	0
50	118	501	0	0
33	124	496	1	0
30	126	502	0	0
50	116	499	0	0
210	456	3706	3	1
50	114	499	0	0
1	113	501	1	1
45	99	1431	4	0
45	98	1431	4	0
45	97	1431	4	0
45	100	1428	2	0
45	100	1431	2	0
45	100	1430	2	0
45	100	1429	2	0
45	98	1437	4	0
45	97	1437	4	0
45	99	1434	0	0
45	99	1437	4	0
45	100	1437	2	0
15	122	1400	4	0
45	100	1435	2	0
6	123	1407	0	0
45	100	1436	2	0
15	118	1400	4	0
2	121	1404	1	1
6	104	1417	0	0
45	100	1434	2	0
44	97	1434	4	0
15	118	1405	4	0
15	104	1419	2	0
1	124	1409	1	1
2	121	1408	1	1
2	122	1404	0	1
55	119	504	0	0
90	116	510	2	0
49	106	504	6	0
30	124	505	2	0
1	115	511	0	1
1	104	506	0	1
1	117	506	0	1
3	106	518	0	0
3	110	518	2	0
7	109	518	6	0
50	105	518	0	0
5	106	519	0	0
29	115	514	0	0
7	113	519	4	0
5	113	518	2	0
29	104	514	0	0
1	104	518	0	1
1	127	513	0	1
15	96	1470	6	0
6	82	1478	0	0
6	94	1473	0	0
3	102	1479	0	0
6	106	1463	0	0
15	104	1466	6	0
3	106	1480	0	0
15	112	1463	2	0
17	106	1467	0	0
6	119	1470	0	0
24	98	1468	6	0
6	96	1467	4	0
6	113	1462	0	0
6	102	1476	0	0
6	82	1469	0	0
3	104	1463	0	0
15	116	1464	4	0
2	119	1466	0	1
1	119	1467	1	1
3	118	1470	0	0
15	116	1469	0	0
6	107	1479	0	0
3	107	1480	0	0
3	110	1476	0	0
17	111	1476	0	0
15	112	1478	0	0
6	113	1477	0	0
7	124	1470	0	0
1	124	1467	1	1
6	121	1474	0	0
15	125	1469	0	0
1	123	1466	0	1
15	121	1469	0	0
7	123	1470	0	0
15	125	1464	4	0
2	121	1466	0	1
3	126	1461	0	0
1	122	1468	0	1
187	126	1467	4	0
15	124	1460	0	0
15	125	1477	0	0
1	122	1476	0	1
7	124	1457	3	0
3	120	1464	0	0
2	124	1475	1	1
15	125	1474	0	0
15	121	1477	0	0
51	111	3368	6	0
3	111	3375	0	0
51	108	3375	0	0
51	103	3377	0	0
51	105	3379	6	0
14	101	3377	0	0
51	113	3371	6	0
3	108	3380	0	0
22	122	3355	0	0
22	121	3355	0	0
5	106	3366	0	0
51	111	3366	6	0
81	101	3380	0	0
51	121	3352	0	0
3	112	3376	0	0
14	111	3377	0	0
3	111	3376	0	0
27	110	3377	0	0
51	106	3375	0	0
51	108	3366	0	0
51	113	3374	6	0
5	118	3352	0	0
25	112	3380	0	0
3	112	3375	0	0
19	110	3370	0	1
1	108	3377	0	1
25	109	3380	0	0
1	105	3377	0	1
1	134	519	1	0
5	129	518	0	0
91	134	516	4	0
89	130	514	0	0
1	130	515	1	1
3	129	1458	0	0
7	129	1459	0	0
15	133	1470	0	0
11	133	1467	0	0
3	135	1467	0	0
7	135	1468	0	0
7	129	1457	4	0
6	129	1462	0	0
1	128	1461	0	1
29	138	504	0	0
30	128	507	0	0
96	139	507	6	0
30	130	507	0	0
1	140	510	0	1
1	140	507	0	1
1	138	496	1	0
1	133	499	0	0
29	149	499	4	0
34	137	497	1	0
30	128	502	0	0
1	133	496	0	0
29	148	501	2	0
6	147	498	6	0
20	144	496	0	0
1	129	496	0	0
58	131	501	6	0
1	129	499	0	0
6	136	1461	1	0
6	136	1474	0	0
15	137	1470	0	0
85	145	1477	0	0
6	138	1466	0	0
6	141	1456	0	0
15	138	1459	1	0
42	145	1473	4	0
1	148	1477	1	1
34	137	495	1	0
1	137	491	1	0
37	136	490	0	0
32	139	491	1	0
34	136	493	1	0
34	137	490	1	0
37	136	495	0	0
34	138	489	1	0
34	138	494	1	0
37	138	493	0	0
34	138	492	1	0
34	136	492	1	0
23	135	1402	6	0
23	135	1399	6	0
23	130	1399	6	0
23	133	1399	6	0
23	133	1402	6	0
15	141	1401	2	0
1	136	1414	1	1
19	137	1401	6	0
47	141	1412	4	0
1	128	1414	1	1
1	138	1413	1	1
25	137	1403	6	0
15	140	1404	2	0
45	130	1414	0	0
3	141	1403	2	0
15	138	1407	0	0
44	131	1414	4	0
15	140	1407	0	0
45	134	1415	2	0
23	130	1402	6	0
45	129	1415	6	0
25	137	1399	6	0
45	129	1414	6	0
3	141	1406	2	0
45	134	1414	2	0
45	133	1414	0	0
6	142	1398	0	0
63	131	1404	6	0
14	139	1409	2	0
45	131	1417	4	0
45	133	1417	4	0
45	132	1417	4	0
45	130	1417	4	0
27	137	1418	0	0
14	136	1421	0	0
25	140	1423	4	0
2	138	1418	0	1
45	134	1417	2	0
45	134	1416	2	0
45	129	1417	6	0
25	140	1419	4	0
56	140	1420	4	0
45	129	1416	6	0
7	153	503	2	0
7	153	502	2	0
7	153	501	2	0
1	160	503	4	0
45	167	503	4	0
14	165	499	0	0
45	166	502	6	0
45	166	503	6	0
47	174	496	4	0
14	173	499	0	0
3	170	501	2	0
45	170	503	2	0
45	168	503	4	0
45	169	503	4	0
45	170	502	2	0
1	168	502	0	1
1	170	496	1	1
1	168	499	0	1
1	170	498	1	1
7	171	501	0	0
7	171	497	0	0
3	172	497	0	0
7	173	497	0	0
4	57	401	0	0
70	58	403	0	0
70	49	407	0	0
38	58	405	0	0
4	71	400	0	0
70	68	402	0	0
70	68	400	0	0
205	69	406	0	0
4	51	395	0	0
38	68	398	0	0
70	48	397	0	0
70	57	398	0	0
70	65	398	0	0
4	53	396	0	0
70	52	392	0	0
205	71	394	0	0
205	70	398	0	0
70	53	392	0	0
4	52	393	0	0
4	53	393	0	0
70	54	388	0	0
70	55	390	0	0
70	53	391	0	0
38	55	386	0	0
70	54	391	0	0
38	53	386	0	0
70	69	388	0	0
70	56	387	0	0
70	53	385	0	0
205	52	384	0	0
38	64	387	0	0
4	68	388	0	0
205	62	385	0	0
70	66	389	0	0
70	73	402	0	0
70	78	401	0	0
70	73	399	0	0
70	78	394	0	0
205	78	385	0	0
38	78	386	0	0
70	73	384	0	0
38	75	384	0	0
70	55	383	0	0
70	49	379	0	0
70	70	376	0	0
70	68	382	0	0
4	57	376	0	0
70	59	381	0	0
70	63	382	0	0
38	61	378	0	0
70	51	382	0	0
70	62	377	0	0
70	48	377	0	0
70	52	383	0	0
4	78	377	0	0
70	70	373	0	0
70	58	372	0	0
38	65	370	0	0
70	59	373	0	0
70	59	372	0	0
38	48	374	0	0
70	56	370	0	0
70	70	374	0	0
70	52	373	0	0
70	71	375	0	0
70	55	375	0	0
70	49	373	0	0
70	49	368	0	0
205	74	373	0	0
4	74	371	0	0
4	77	373	0	0
70	53	363	0	0
4	53	361	0	0
70	75	367	0	0
70	76	361	0	0
70	68	366	0	0
70	78	363	0	0
205	69	363	0	0
70	49	362	0	0
70	71	361	0	0
4	55	366	0	0
70	77	364	0	0
70	77	362	0	0
4	72	362	0	0
205	72	367	0	0
4	52	363	0	0
4	55	353	0	0
70	48	358	0	0
70	50	359	0	0
70	52	358	0	0
70	62	354	0	0
70	65	359	0	0
205	66	352	0	0
38	77	357	0	0
70	77	353	0	0
4	71	359	0	0
70	73	358	0	0
205	48	357	0	0
70	58	356	0	0
4	75	357	0	0
38	61	353	0	0
4	63	348	0	0
70	58	350	0	0
70	54	351	0	0
70	66	351	0	0
70	68	345	0	0
4	48	336	0	0
70	71	336	0	0
205	62	338	0	0
70	54	340	0	0
70	69	337	0	0
38	72	340	0	0
70	64	339	0	0
70	74	343	0	0
70	59	333	0	0
4	63	331	0	0
205	74	330	0	0
70	61	335	0	0
38	62	335	0	0
70	72	330	0	0
70	62	332	0	0
70	72	331	0	0
70	55	323	0	0
4	53	326	0	0
70	61	326	0	0
205	55	326	0	0
207	60	323	0	0
4	60	321	0	0
70	74	323	0	0
70	68	321	0	0
205	56	324	0	0
205	49	322	0	0
205	59	322	0	0
70	57	326	0	0
70	70	323	0	0
70	52	316	0	0
70	60	312	0	0
4	62	313	0	0
206	64	319	0	0
38	76	312	0	0
205	59	312	0	0
38	61	316	0	0
70	58	317	0	0
4	61	318	0	0
4	65	315	0	0
38	65	313	0	0
206	65	318	0	0
204	66	319	0	0
4	87	316	0	0
70	87	313	0	0
70	82	313	0	0
70	86	318	0	0
205	83	330	0	0
70	80	324	0	0
70	86	324	0	0
205	82	329	0	0
70	86	350	0	0
70	82	338	0	0
70	80	337	0	0
70	80	349	0	0
70	83	327	0	0
70	84	346	0	0
70	82	340	0	0
4	83	326	0	0
70	80	343	0	0
38	87	327	0	0
70	84	330	0	0
70	91	316	0	0
70	88	327	0	0
70	90	320	0	0
38	91	322	0	0
70	88	315	0	0
4	90	326	0	0
70	88	319	0	0
70	94	329	0	0
205	88	336	0	0
70	92	351	0	0
4	94	321	0	0
70	88	322	0	0
4	96	319	0	0
70	97	319	0	0
38	100	315	0	0
70	100	318	0	0
4	96	316	0	0
70	102	318	0	0
205	102	331	0	0
205	100	332	0	0
205	96	335	0	0
38	96	336	0	0
70	101	341	0	0
205	100	338	0	0
205	101	337	0	0
70	101	350	0	0
70	98	350	0	0
70	103	342	0	0
4	70	311	0	0
205	73	310	0	0
4	83	305	0	0
70	83	310	0	0
70	86	311	0	0
38	89	304	0	0
70	89	307	0	0
4	90	311	0	0
205	88	308	0	0
4	92	305	0	0
205	102	305	0	0
205	98	305	0	0
38	81	311	0	0
70	88	304	0	0
4	58	305	0	0
70	60	305	0	0
70	58	309	0	0
70	60	308	0	0
70	62	305	0	0
70	59	299	0	0
70	71	297	0	0
70	71	298	0	0
70	69	299	0	0
4	63	296	0	0
70	60	303	0	0
38	70	302	0	0
38	86	298	0	0
4	91	297	0	0
70	94	303	0	0
205	92	300	0	0
38	66	301	0	0
4	65	296	0	0
70	72	297	0	0
70	88	301	0	0
70	67	297	0	0
205	93	296	0	0
70	85	296	0	0
70	57	295	0	0
205	94	295	0	0
70	91	288	0	0
70	74	295	0	0
70	80	292	0	0
70	81	288	0	0
70	82	295	0	0
70	60	289	0	0
205	73	295	0	0
4	73	294	0	0
70	84	290	0	0
70	49	299	0	0
70	55	297	0	0
70	53	301	0	0
70	50	302	0	0
70	54	296	0	0
4	48	293	0	0
70	48	281	0	0
205	55	284	0	0
4	71	287	0	0
70	50	287	0	0
214	61	286	6	0
70	53	280	0	0
215	66	284	2	0
216	63	286	7	0
205	60	286	0	0
214	67	286	7	0
38	60	282	0	0
214	68	284	2	0
38	57	286	0	0
70	58	286	0	0
4	54	281	0	0
70	58	283	0	0
205	53	287	0	0
4	53	281	0	0
217	65	286	6	0
215	69	285	7	0
217	65	284	2	0
70	63	282	0	0
205	48	279	0	0
70	69	273	0	0
70	70	278	0	0
70	57	278	0	0
70	49	273	0	0
70	54	276	0	0
70	50	272	0	0
70	54	278	0	0
4	63	268	0	0
70	61	271	0	0
70	61	266	0	0
70	52	269	0	0
205	54	269	0	0
4	66	268	0	0
70	66	271	0	0
70	69	266	0	0
70	70	265	0	0
4	61	264	0	0
38	63	258	0	0
70	70	258	0	0
70	48	261	0	0
70	67	257	0	0
4	66	258	0	0
70	50	258	0	0
4	50	261	0	0
70	60	262	0	0
70	60	263	0	0
70	53	251	0	0
4	66	250	0	0
4	64	251	0	0
70	57	250	0	0
70	50	252	0	0
205	62	248	0	0
70	59	254	0	0
70	63	244	0	0
70	56	241	0	0
38	65	242	0	0
205	65	243	0	0
70	51	246	0	0
205	64	240	0	0
70	64	242	0	0
70	66	236	0	0
70	71	238	0	0
205	51	234	0	0
70	64	232	0	0
4	55	237	0	0
70	65	233	0	0
205	56	230	0	0
4	68	227	0	0
38	57	226	0	0
70	61	230	0	0
70	50	227	0	0
70	65	229	0	0
4	50	224	0	0
70	50	218	0	0
4	53	223	0	0
70	68	218	0	0
70	71	221	0	0
70	56	216	0	0
70	67	219	0	0
70	58	217	0	0
70	71	208	0	0
4	54	214	0	0
38	54	212	0	0
70	68	214	0	0
70	50	208	0	0
70	63	211	0	0
38	61	214	0	0
70	71	211	0	0
38	60	211	0	0
4	50	205	0	0
70	50	204	0	0
70	58	205	0	0
70	56	202	0	0
70	69	205	0	0
70	65	202	0	0
205	66	205	0	0
70	65	200	0	0
70	55	196	0	0
70	55	194	0	0
70	48	196	0	0
4	55	192	0	0
70	51	197	0	0
70	65	199	0	0
70	54	191	0	0
70	61	187	0	0
70	57	188	0	0
4	71	191	0	0
70	71	190	0	0
4	71	187	0	0
70	68	190	0	0
4	70	191	0	0
205	70	185	0	0
70	69	186	0	0
70	64	188	0	0
4	54	185	0	0
70	70	182	0	0
205	63	183	0	0
70	59	177	0	0
38	69	179	0	0
70	71	180	0	0
70	69	176	0	0
70	68	182	0	0
70	51	175	0	0
70	52	172	0	0
70	59	169	0	0
4	58	168	0	0
70	66	169	0	0
70	71	170	0	0
70	57	168	0	0
70	53	164	0	0
38	60	160	0	0
4	48	166	0	0
70	53	165	0	0
38	58	165	0	0
70	58	162	0	0
38	57	163	0	0
4	51	163	0	0
70	51	158	0	0
4	55	159	0	0
70	63	156	0	0
70	71	157	0	0
70	59	159	0	0
70	65	157	0	0
70	69	157	0	0
38	67	153	0	0
70	52	152	0	0
70	62	145	0	0
70	53	150	0	0
70	64	148	0	0
4	60	151	0	0
70	55	145	0	0
70	55	147	0	0
4	66	149	0	0
38	67	149	0	0
70	65	151	0	0
70	54	144	0	0
70	71	138	0	0
70	58	139	0	0
20	54	142	0	0
20	60	142	0	0
70	60	136	0	0
70	69	137	0	0
70	70	139	0	0
70	51	137	0	0
20	67	142	0	0
70	54	128	0	0
205	65	131	0	0
70	76	131	0	0
205	76	136	0	0
70	76	138	0	0
70	74	163	0	0
70	74	146	0	0
70	76	147	0	0
4	77	144	0	0
70	76	148	0	0
4	75	157	0	0
70	76	167	0	0
4	76	153	0	0
70	74	156	0	0
20	75	142	0	0
70	86	139	0	0
205	85	163	0	0
70	84	128	0	0
38	84	157	0	0
207	86	165	0	0
70	82	150	0	0
205	84	144	0	0
70	81	136	0	0
207	80	160	0	0
70	81	156	0	0
207	84	160	3	0
70	84	131	0	0
4	82	149	0	0
70	80	157	0	0
20	83	142	0	0
207	83	165	5	0
70	86	152	0	0
205	88	134	0	0
70	92	139	0	0
205	89	160	0	0
70	89	157	0	0
205	92	148	0	0
70	95	157	0	0
4	92	146	0	0
70	92	134	0	0
70	91	129	0	0
70	94	151	0	0
70	91	128	0	0
70	90	158	0	0
70	95	156	0	0
4	93	153	0	0
70	91	156	0	0
20	89	142	0	0
70	95	149	0	0
70	91	135	0	0
70	85	183	0	0
70	84	183	0	0
70	75	179	0	0
70	82	179	0	0
70	81	182	0	0
70	80	185	0	0
38	80	189	0	0
38	75	188	0	0
70	76	185	0	0
70	76	197	0	0
70	74	192	0	0
70	86	192	0	0
4	81	194	0	0
70	86	205	0	0
70	74	202	0	0
70	84	207	0	0
4	72	204	0	0
70	76	203	0	0
70	84	203	0	0
70	83	212	0	0
4	74	210	0	0
38	75	219	0	0
70	74	218	0	0
4	73	222	0	0
70	72	220	0	0
70	76	220	0	0
70	84	219	0	0
70	86	218	0	0
70	83	223	0	0
4	83	229	0	0
70	87	231	0	0
70	72	228	0	0
70	79	225	0	0
70	79	227	0	0
70	87	236	0	0
70	81	238	0	0
70	73	233	0	0
70	73	236	0	0
38	73	243	0	0
70	87	240	0	0
70	82	240	0	0
38	85	241	0	0
70	86	243	0	0
70	76	248	0	0
4	74	254	0	0
4	75	249	0	0
38	81	248	0	0
4	72	254	0	0
4	77	255	0	0
70	79	251	0	0
4	74	252	0	0
205	78	253	0	0
70	74	248	0	0
70	85	254	0	0
4	85	249	0	0
70	74	257	0	0
70	77	263	0	0
4	72	263	0	0
38	87	263	0	0
70	76	261	0	0
70	86	262	0	0
70	78	269	0	0
205	72	265	0	0
70	80	271	0	0
70	72	266	0	0
70	86	267	0	0
70	81	264	0	0
4	79	274	0	0
215	87	278	0	0
70	77	272	0	0
38	77	277	0	0
215	87	279	5	0
214	86	277	7	0
70	82	277	0	0
214	85	279	7	0
70	86	284	0	0
217	87	282	0	0
216	74	283	0	0
215	75	281	7	0
214	74	286	7	0
216	84	281	0	0
217	78	281	0	0
205	77	283	0	0
70	82	357	0	0
70	89	356	0	0
38	89	357	0	0
70	99	353	0	0
70	93	354	0	0
70	97	353	0	0
70	101	356	0	0
70	92	357	0	0
70	94	353	0	0
38	84	358	0	0
70	95	356	0	0
38	99	360	0	0
70	93	367	0	0
70	82	364	0	0
70	93	360	0	0
70	103	362	0	0
70	89	362	0	0
38	97	364	0	0
70	95	362	0	0
4	102	362	0	0
70	100	367	0	0
70	81	360	0	0
70	100	363	0	0
70	81	363	0	0
4	95	366	0	0
38	80	366	0	0
70	82	368	0	0
70	84	373	0	0
38	81	372	0	0
70	99	372	0	0
205	94	374	0	0
70	100	369	0	0
4	98	371	0	0
70	93	371	0	0
205	103	370	0	0
70	101	374	0	0
4	83	377	0	0
4	93	383	0	0
70	80	376	0	0
70	87	379	0	0
70	85	381	0	0
205	84	382	0	0
38	87	378	0	0
70	84	388	0	0
70	95	389	0	0
4	80	384	0	0
38	89	388	0	0
70	90	391	0	0
70	87	386	0	0
70	95	386	0	0
205	84	396	0	0
4	83	395	0	0
4	85	392	0	0
70	93	393	0	0
38	81	405	0	0
4	85	401	0	0
4	83	405	0	0
70	95	406	0	0
70	80	401	0	0
38	87	401	0	0
70	87	402	0	0
4	97	299	0	0
70	101	298	0	0
70	98	303	0	0
70	102	288	0	0
70	99	290	0	0
70	99	295	0	0
70	99	293	0	0
70	94	283	0	0
38	98	287	0	0
4	96	287	0	0
205	92	281	0	0
70	90	282	0	0
214	97	281	7	0
70	93	284	0	0
214	97	283	6	0
70	88	285	0	0
70	100	280	0	0
70	100	287	0	0
70	88	280	0	0
215	97	275	5	0
215	90	279	5	0
215	99	275	6	0
216	97	277	5	0
214	93	278	6	0
70	93	272	0	0
70	91	278	0	0
70	101	274	0	0
70	101	279	0	0
70	88	276	0	0
70	101	275	0	0
217	102	276	5	0
70	94	275	0	0
215	99	272	6	0
70	93	279	0	0
70	100	269	0	0
70	93	271	0	0
4	88	265	0	0
70	103	271	0	0
70	90	267	0	0
38	90	270	0	0
70	92	266	0	0
70	98	267	0	0
70	102	257	0	0
70	91	260	0	0
70	101	260	0	0
70	96	257	0	0
4	97	262	0	0
4	96	259	0	0
4	95	255	0	0
70	95	248	0	0
70	88	248	0	0
70	93	254	0	0
4	90	255	0	0
70	88	240	0	0
70	90	240	0	0
70	92	247	0	0
205	92	235	0	0
70	93	237	0	0
70	92	238	0	0
70	91	239	0	0
38	92	232	0	0
205	92	227	0	0
4	93	231	0	0
70	95	222	0	0
70	93	219	0	0
70	91	221	0	0
70	95	210	0	0
4	94	210	0	0
70	89	212	0	0
70	95	209	0	0
70	88	210	0	0
4	92	207	0	0
70	88	203	0	0
70	89	201	0	0
70	91	205	0	0
70	91	192	0	0
70	90	197	0	0
70	93	199	0	0
70	88	197	0	0
70	93	185	0	0
205	88	185	0	0
70	90	184	0	0
70	94	181	0	0
70	94	183	0	0
70	92	174	0	0
70	102	170	0	0
70	101	188	0	0
38	99	163	0	0
205	103	160	0	0
70	103	166	0	0
70	100	196	0	0
70	98	188	0	0
4	98	190	0	0
70	102	197	0	0
70	100	188	0	0
70	110	164	0	0
205	111	188	0	0
70	107	165	0	0
70	106	181	0	0
70	106	175	0	0
205	111	178	0	0
70	108	182	0	0
70	107	180	0	0
70	108	194	0	0
70	110	191	0	0
70	109	194	0	0
38	111	185	0	0
70	106	185	0	0
70	108	185	0	0
70	110	161	0	0
4	110	192	0	0
70	104	189	0	0
70	98	158	0	0
70	99	158	0	0
38	99	159	0	0
4	100	155	0	0
205	110	158	0	0
70	111	149	0	0
70	101	147	0	0
4	110	151	0	0
347	111	142	6	0
70	100	139	0	0
20	98	142	0	0
70	110	139	0	0
70	98	141	0	0
20	107	142	0	0
205	100	135	0	0
70	110	133	0	0
38	105	128	0	0
70	103	204	0	0
70	96	204	0	0
70	96	209	0	0
70	99	213	0	0
70	98	214	0	0
88	102	219	2	0
88	102	222	2	0
70	97	217	0	0
70	101	222	0	0
88	103	217	2	0
70	96	221	0	0
70	96	218	0	0
70	96	220	0	0
205	102	229	0	0
70	101	225	0	0
4	99	226	0	0
70	102	228	0	0
38	98	227	0	0
70	99	234	0	0
4	96	244	0	0
70	98	248	0	0
70	102	249	0	0
70	100	249	0	0
70	102	383	0	0
205	97	385	0	0
70	96	391	0	0
70	97	387	0	0
4	98	388	0	0
70	99	384	0	0
70	103	385	0	0
205	99	404	0	0
4	103	400	0	0
4	102	396	0	0
70	110	333	0	0
70	111	335	0	0
70	107	334	0	0
70	108	328	0	0
70	108	332	0	0
70	109	340	0	0
70	104	331	0	0
70	109	344	0	0
70	107	329	0	0
70	107	354	0	0
70	110	362	0	0
4	106	352	0	0
70	106	373	0	0
205	106	374	0	0
70	106	375	0	0
38	105	379	0	0
70	111	378	0	0
70	106	377	0	0
70	111	381	0	0
70	111	388	0	0
38	107	389	0	0
70	108	387	0	0
70	111	394	0	0
70	105	397	0	0
70	107	397	0	0
70	105	396	0	0
4	106	395	0	0
70	106	404	0	0
220	116	405	2	0
219	116	407	0	0
218	116	406	2	0
70	115	400	0	0
70	116	401	0	0
70	113	392	0	0
4	112	395	0	0
70	112	397	0	0
70	118	397	0	0
70	116	392	0	0
38	114	388	0	0
4	117	390	0	0
70	112	380	0	0
70	118	377	0	0
205	117	381	0	0
70	115	381	0	0
4	116	379	0	0
70	118	375	0	0
205	117	371	0	0
20	119	364	0	0
20	119	366	0	0
296	116	366	0	0
20	119	362	0	0
20	114	364	6	0
20	114	366	6	0
20	114	362	6	0
143	113	362	2	0
143	113	365	2	0
70	118	358	0	0
63	116	359	2	0
205	119	351	0	0
70	116	351	0	0
4	113	349	0	0
205	115	351	0	0
38	117	342	0	0
70	113	338	0	0
38	113	342	0	0
70	116	343	0	0
143	120	362	6	0
38	124	375	0	0
205	127	380	0	0
70	125	373	0	0
70	126	382	0	0
38	122	388	0	0
70	125	388	0	0
143	120	365	6	0
70	125	379	0	0
70	125	353	0	0
70	121	384	0	0
70	120	375	0	0
70	123	363	0	0
70	122	383	0	0
4	123	359	0	0
70	125	365	0	0
70	120	388	0	0
70	127	386	0	0
70	127	398	0	0
70	126	398	0	0
205	127	392	0	0
38	125	394	0	0
70	123	392	0	0
70	121	398	0	0
70	134	363	0	0
70	134	362	0	0
70	130	362	0	0
70	134	373	0	0
70	133	369	0	0
70	132	392	0	0
70	130	391	0	0
4	131	383	0	0
70	130	390	0	0
38	132	395	0	0
70	128	370	0	0
70	131	369	0	0
70	132	375	0	0
4	135	392	0	0
70	133	383	0	0
4	135	394	0	0
70	129	393	0	0
70	126	404	0	0
70	121	404	0	0
70	128	404	0	0
70	128	406	0	0
38	121	400	0	0
4	140	393	0	0
4	142	393	0	0
70	136	399	0	0
70	140	390	0	0
4	142	391	0	0
4	143	376	0	0
70	140	379	0	0
38	139	382	0	0
4	139	378	0	0
70	141	381	0	0
70	138	378	0	0
70	139	379	0	0
70	141	372	0	0
4	137	368	0	0
70	136	374	0	0
70	142	367	0	0
70	143	366	0	0
70	140	360	0	0
4	142	365	0	0
205	110	320	0	0
70	111	319	0	0
70	105	324	0	0
70	106	324	0	0
70	109	316	0	0
38	110	316	0	0
70	109	320	0	0
70	104	312	0	0
205	107	308	0	0
70	106	311	0	0
70	104	309	0	0
70	104	308	0	0
70	109	296	0	0
4	104	302	0	0
70	104	303	0	0
70	107	303	0	0
70	105	303	0	0
70	108	295	0	0
70	108	282	0	0
70	110	282	0	0
214	110	272	4	0
70	106	278	0	0
70	108	270	0	0
70	108	266	0	0
215	109	269	4	0
70	110	267	0	0
70	106	264	0	0
70	110	271	0	0
70	105	266	0	0
70	105	264	0	0
70	107	271	0	0
70	110	261	0	0
70	107	258	0	0
4	104	262	0	0
70	109	254	0	0
38	105	251	0	0
70	108	252	0	0
70	109	242	0	0
70	107	236	0	0
205	109	235	0	0
38	108	234	0	0
4	106	231	0	0
88	107	224	2	0
88	105	224	2	0
208	109	224	4	0
88	106	215	2	0
88	109	213	2	0
208	104	220	2	0
108	109	222	2	0
70	108	221	0	0
4	110	218	0	0
70	111	212	0	0
88	104	214	2	0
108	108	223	2	0
88	105	217	2	0
88	107	212	2	0
88	107	219	2	0
88	106	213	2	0
70	80	122	0	0
205	75	123	0	0
205	73	120	0	0
205	83	121	0	0
70	75	124	0	0
205	95	125	0	0
70	74	124	0	0
70	106	124	0	0
205	72	124	0	0
70	74	126	0	0
70	90	126	0	0
70	97	125	0	0
70	98	124	0	0
70	116	137	0	0
70	116	139	0	0
70	114	126	0	0
20	116	142	0	0
70	115	134	0	0
70	116	157	0	0
70	116	149	0	0
4	116	144	0	0
70	116	155	0	0
70	118	149	0	0
70	115	153	0	0
70	115	158	0	0
70	113	158	0	0
70	113	160	0	0
70	116	162	0	0
70	119	175	0	0
70	118	175	0	0
70	112	172	0	0
70	119	182	0	0
70	112	185	0	0
70	115	186	0	0
70	114	187	0	0
70	119	191	0	0
205	113	194	0	0
70	115	196	0	0
70	115	195	0	0
70	118	194	0	0
70	112	196	0	0
70	112	194	0	0
38	115	201	0	0
38	112	215	0	0
70	112	209	0	0
70	116	212	0	0
88	113	213	2	0
70	117	215	0	0
88	115	214	2	0
88	116	222	2	0
88	117	219	2	0
208	114	217	2	0
70	116	219	0	0
208	112	223	0	0
70	114	332	0	0
4	117	323	0	0
205	113	326	0	0
70	112	332	0	0
4	119	318	0	0
205	115	318	0	0
4	116	317	0	0
4	112	306	0	0
70	114	308	0	0
70	114	301	0	0
70	117	302	0	0
4	117	294	0	0
205	119	294	0	0
70	112	292	0	0
4	114	291	0	0
4	114	284	0	0
70	115	280	0	0
4	118	284	0	0
70	118	283	0	0
214	114	272	4	0
70	118	274	0	0
38	119	272	0	0
217	112	274	5	0
216	119	264	7	0
70	117	266	0	0
214	112	268	3	0
217	113	267	6	0
216	116	269	6	0
214	116	265	3	0
70	113	271	0	0
70	119	259	0	0
215	116	263	2	0
205	115	259	0	0
217	117	261	0	0
4	113	253	0	0
38	113	252	0	0
70	116	255	0	0
70	112	240	0	0
70	114	247	0	0
70	114	242	0	0
70	123	172	0	0
70	122	176	0	0
70	122	172	0	0
70	122	191	0	0
70	121	191	0	0
38	123	189	0	0
70	125	203	0	0
70	121	194	0	0
70	121	186	0	0
70	127	168	0	0
4	123	213	0	0
205	124	214	0	0
205	125	217	0	0
4	124	220	0	0
4	127	222	0	0
70	126	222	0	0
70	119	226	0	0
205	116	230	0	0
70	116	231	0	0
70	117	224	0	0
38	119	224	0	0
38	115	225	0	0
70	121	226	0	0
70	120	224	0	0
70	121	227	0	0
70	124	232	0	0
70	123	235	0	0
70	126	233	0	0
70	123	232	0	0
70	116	232	0	0
70	127	233	0	0
38	120	239	0	0
70	118	237	0	0
70	126	240	0	0
70	124	245	0	0
70	120	248	0	0
70	123	255	0	0
70	122	252	0	0
70	122	250	0	0
70	124	258	0	0
70	126	263	0	0
215	122	260	3	0
214	122	263	4	0
215	127	268	2	0
70	123	267	0	0
70	122	267	0	0
70	121	271	0	0
215	120	268	2	0
70	122	266	0	0
70	122	276	0	0
70	121	273	0	0
70	125	279	0	0
70	125	282	0	0
70	124	280	0	0
70	127	295	0	0
70	120	288	0	0
205	125	295	0	0
205	123	290	0	0
70	120	289	0	0
70	123	303	0	0
70	127	309	0	0
205	122	305	0	0
4	122	304	0	0
70	124	307	0	0
38	122	312	0	0
70	122	317	0	0
205	122	315	0	0
70	126	312	0	0
4	127	320	0	0
205	127	326	0	0
70	126	323	0	0
70	124	327	0	0
70	125	331	0	0
205	122	335	0	0
70	127	329	0	0
4	120	335	0	0
70	122	330	0	0
70	120	337	0	0
70	124	337	0	0
4	120	340	0	0
70	125	340	0	0
4	127	346	0	0
70	127	349	0	0
38	123	350	0	0
70	121	348	0	0
70	135	330	0	0
205	133	359	0	0
38	130	336	0	0
70	133	345	0	0
205	128	337	0	0
38	134	351	0	0
70	132	345	0	0
4	129	346	0	0
70	128	333	0	0
205	134	231	0	0
70	135	228	0	0
70	128	216	0	0
4	135	221	0	0
4	129	220	0	0
70	134	247	0	0
70	131	233	0	0
205	134	225	0	0
70	128	255	0	0
205	133	230	0	0
4	128	254	0	0
4	131	240	0	0
70	129	252	0	0
70	132	251	0	0
70	135	247	0	0
70	135	212	0	0
70	133	209	0	0
205	129	207	0	0
70	130	205	0	0
70	134	205	0	0
70	130	200	0	0
38	130	198	0	0
70	129	195	0	0
205	121	167	0	0
70	127	160	0	0
205	126	160	0	0
70	125	162	0	0
70	122	160	0	0
70	123	166	0	0
38	120	158	0	0
70	123	158	0	0
70	123	159	0	0
70	130	152	0	0
70	133	161	0	0
38	128	166	0	0
4	129	161	0	0
70	132	159	0	0
38	130	153	0	0
70	133	157	0	0
70	132	155	0	0
38	133	170	0	0
70	132	177	0	0
70	128	169	0	0
70	124	149	0	0
4	131	147	0	0
70	133	147	0	0
70	135	147	0	0
4	125	149	0	0
38	132	148	0	0
70	132	151	0	0
38	129	149	0	0
70	133	150	0	0
70	122	151	0	0
4	121	145	0	0
70	128	146	0	0
38	122	147	0	0
4	134	146	0	0
205	126	139	0	0
70	121	137	0	0
20	123	142	0	0
20	132	142	0	0
70	130	137	0	0
38	128	129	0	0
70	123	128	0	0
70	134	134	0	0
70	129	134	0	0
70	120	133	0	0
70	137	131	0	0
20	138	142	0	0
70	136	129	0	0
70	142	144	0	0
4	136	154	0	0
70	136	165	0	0
70	139	157	0	0
70	141	157	0	0
70	138	131	0	0
70	136	137	0	0
70	137	138	0	0
4	139	167	0	0
70	139	148	0	0
205	141	150	0	0
70	143	157	0	0
70	137	151	0	0
70	136	169	0	0
70	140	169	0	0
93	140	180	2	0
38	139	177	0	0
70	137	176	0	0
38	143	196	0	0
38	139	192	0	0
70	136	196	0	0
70	141	196	0	0
4	136	203	0	0
4	137	211	0	0
4	141	212	0	0
70	142	209	0	0
70	139	211	0	0
4	130	261	0	0
214	131	263	3	0
70	134	258	0	0
70	129	257	0	0
217	128	262	1	0
70	134	259	0	0
70	135	262	0	0
216	133	268	2	0
217	133	264	0	0
205	131	269	0	0
215	129	268	2	0
205	135	271	0	0
70	134	271	0	0
214	128	271	4	0
70	135	267	0	0
215	134	266	3	0
214	128	265	4	0
214	135	277	2	0
215	133	279	1	0
70	131	277	0	0
205	134	283	0	0
216	135	281	1	0
217	135	283	0	0
70	132	282	0	0
4	130	285	0	0
70	133	284	0	0
38	134	298	0	0
70	132	297	0	0
70	132	300	0	0
70	135	304	0	0
70	129	305	0	0
70	135	307	0	0
70	132	311	0	0
70	129	304	0	0
70	133	316	0	0
38	135	319	0	0
70	128	315	0	0
70	130	314	0	0
38	128	313	0	0
70	133	324	0	0
70	133	320	0	0
70	135	321	0	0
4	136	343	0	0
70	142	347	0	0
4	136	336	0	0
70	136	332	0	0
70	138	337	0	0
70	140	339	0	0
70	138	342	0	0
70	137	340	0	0
70	140	333	0	0
70	137	322	0	0
70	138	321	0	0
70	139	313	0	0
70	142	315	0	0
205	139	306	0	0
70	142	311	0	0
4	138	311	0	0
70	140	294	0	0
70	141	289	0	0
70	138	287	0	0
4	138	281	0	0
205	140	274	0	0
70	136	273	0	0
215	137	278	1	0
38	141	274	0	0
70	141	277	0	0
214	138	269	3	0
70	141	268	0	0
70	137	266	0	0
4	139	268	0	0
215	141	270	2	0
70	141	257	0	0
70	142	259	0	0
4	138	256	0	0
38	141	256	0	0
70	143	261	0	0
38	141	252	0	0
205	139	255	0	0
4	142	242	0	0
70	142	245	0	0
70	142	232	0	0
70	136	225	0	0
70	136	228	0	0
70	137	217	0	0
70	137	218	0	0
205	139	223	0	0
70	146	201	0	0
4	145	200	0	0
70	147	212	0	0
38	151	206	0	0
4	146	204	0	0
205	147	219	0	0
70	151	239	0	0
70	147	237	0	0
38	146	234	0	0
205	146	213	0	0
4	146	233	0	0
70	151	218	0	0
205	150	226	0	0
70	149	223	0	0
70	146	243	0	0
218	150	255	0	0
218	149	254	2	0
219	150	253	0	0
205	149	252	0	0
219	149	253	4	0
70	147	252	0	0
218	150	251	2	0
220	151	255	4	0
70	145	251	0	0
70	144	248	0	0
118	146	254	0	0
4	147	251	0	0
218	150	252	2	0
219	149	255	6	0
220	151	261	6	0
219	151	262	0	0
218	151	258	2	0
218	151	259	2	0
218	151	260	2	0
218	151	257	2	0
70	144	261	0	0
206	146	257	4	0
70	144	263	0	0
219	147	262	4	0
218	151	256	2	0
218	148	262	0	0
218	150	262	0	0
218	149	262	0	0
70	149	263	0	0
70	145	271	0	0
205	144	269	0	0
38	145	264	0	0
70	146	264	0	0
70	150	268	0	0
70	144	272	0	0
70	150	278	0	0
70	151	286	0	0
70	151	280	0	0
70	145	287	0	0
70	145	280	0	0
70	147	292	0	0
70	144	289	0	0
70	151	292	0	0
70	149	295	0	0
70	149	293	0	0
38	151	303	0	0
4	150	297	0	0
70	150	308	0	0
4	151	304	0	0
70	147	310	0	0
4	147	313	0	0
70	147	324	0	0
4	146	327	0	0
38	147	320	0	0
70	149	333	0	0
70	151	334	0	0
4	146	329	0	0
4	145	332	0	0
70	149	336	0	0
70	144	340	0	0
70	144	339	0	0
205	147	342	0	0
70	150	351	0	0
70	145	344	0	0
38	150	349	0	0
70	150	356	0	0
70	148	358	0	0
4	145	362	0	0
38	150	366	0	0
70	151	367	0	0
205	144	362	0	0
70	145	361	0	0
4	147	363	0	0
70	148	372	0	0
70	145	369	0	0
70	148	377	0	0
70	146	391	0	0
70	148	394	0	0
70	147	404	0	0
4	145	405	0	0
70	151	401	0	0
70	157	400	0	0
38	153	405	0	0
38	155	406	0	0
38	159	397	0	0
70	157	394	0	0
70	152	394	0	0
4	155	394	0	0
70	156	393	0	0
70	153	384	0	0
70	153	385	0	0
70	157	377	0	0
70	159	378	0	0
70	153	379	0	0
70	153	381	0	0
70	159	375	0	0
4	159	372	0	0
70	156	372	0	0
4	158	373	0	0
70	152	370	0	0
70	152	373	0	0
70	158	374	0	0
38	156	363	0	0
70	154	365	0	0
4	156	362	0	0
70	153	360	0	0
70	154	363	0	0
70	152	362	0	0
205	156	356	0	0
205	159	352	0	0
4	153	356	0	0
70	157	344	0	0
70	152	344	0	0
70	153	349	0	0
70	158	344	0	0
70	159	343	0	0
70	155	340	0	0
70	154	331	0	0
38	154	328	0	0
70	154	264	0	0
70	158	269	0	0
58	155	266	6	0
70	153	268	0	0
70	158	267	0	0
70	155	280	0	0
205	159	276	0	0
70	155	279	0	0
38	153	279	0	0
70	157	282	0	0
70	155	276	0	0
4	157	291	0	0
70	154	292	0	0
70	159	280	0	0
70	152	289	0	0
70	154	283	0	0
70	153	296	0	0
70	157	300	0	0
4	154	298	0	0
70	154	284	0	0
70	159	260	0	0
218	152	261	0	0
205	156	259	0	0
218	153	261	0	0
38	158	256	2	0
218	152	255	0	0
218	154	255	0	0
205	155	254	0	0
70	159	253	0	0
218	153	255	0	0
205	157	244	0	0
38	157	240	0	0
70	156	243	0	0
70	159	235	0	0
70	152	233	0	0
4	157	235	0	0
70	157	227	0	0
70	158	227	0	0
70	154	229	0	0
70	159	220	0	0
70	158	220	0	0
4	154	221	0	0
4	156	223	0	0
205	155	210	0	0
205	159	212	0	0
70	159	207	0	0
70	155	309	0	0
70	155	307	0	0
70	152	314	0	0
70	157	319	0	0
70	156	315	0	0
70	156	316	0	0
205	166	413	0	0
70	161	410	0	0
70	173	408	0	0
70	173	415	0	0
70	166	406	0	0
70	163	400	0	0
4	160	407	0	0
205	160	403	0	0
38	162	400	0	0
70	174	404	0	0
70	164	402	0	0
70	162	395	0	0
70	166	398	0	0
70	164	384	0	0
70	167	384	0	0
70	174	388	0	0
70	162	384	0	0
4	160	391	0	0
4	170	385	0	0
70	164	383	0	0
70	170	381	0	0
70	161	381	0	0
70	175	379	0	0
70	170	376	0	0
70	168	380	0	0
70	173	382	0	0
38	164	375	0	0
70	164	369	0	0
70	164	373	0	0
4	161	370	0	0
38	167	366	0	0
70	164	366	0	0
70	172	373	0	0
38	168	361	0	0
38	170	371	0	0
205	161	352	0	0
70	166	353	0	0
70	175	353	0	0
38	167	346	0	0
70	167	349	0	0
70	166	347	0	0
70	162	343	0	0
4	162	340	0	0
38	160	339	0	0
70	163	335	0	0
13	167	328	1	0
38	161	331	0	0
13	165	328	2	0
38	163	333	0	0
38	162	324	0	0
13	167	322	2	0
70	161	327	0	0
213	167	325	0	0
12	165	325	0	0
13	165	322	2	0
12	167	319	0	0
70	165	317	0	0
12	165	314	0	0
13	165	312	0	0
70	160	310	0	0
70	161	308	0	0
38	167	307	0	0
13	165	308	0	0
12	165	306	0	0
70	165	303	0	0
70	161	303	0	0
70	160	298	0	0
70	164	294	0	0
70	162	289	0	0
70	163	286	0	0
4	166	287	0	0
70	161	284	0	0
70	162	277	0	0
70	166	268	0	0
70	161	270	0	0
38	165	268	0	0
70	165	270	0	0
70	164	267	0	0
70	166	265	0	0
70	164	261	0	0
217	161	263	5	0
208	167	262	0	0
70	162	262	0	0
38	164	260	0	0
70	161	260	0	0
204	164	258	5	0
1	166	259	0	1
1	164	256	1	1
70	162	254	0	0
4	164	250	0	0
15	164	254	4	0
70	162	253	0	0
70	167	251	0	0
70	161	249	0	0
4	166	250	0	0
70	160	246	0	0
38	165	247	0	0
70	167	244	0	0
205	160	244	0	0
70	161	243	0	0
70	162	245	0	0
38	163	232	0	0
70	161	239	0	0
70	160	225	0	0
205	162	227	0	0
4	164	226	0	0
70	167	224	0	0
70	160	221	0	0
70	167	217	0	0
70	165	218	0	0
4	166	219	0	0
70	162	220	0	0
70	161	208	0	0
38	166	215	0	0
4	162	208	0	0
70	169	230	0	0
70	175	216	0	0
70	173	227	0	0
38	174	222	0	0
4	173	228	0	0
38	168	213	0	0
38	175	222	0	0
70	168	225	0	0
4	172	225	0	0
70	169	216	0	0
70	169	220	0	0
70	171	236	0	0
38	168	235	0	0
70	169	244	0	0
70	174	224	0	0
70	169	211	0	0
70	174	241	0	0
38	170	208	0	0
70	168	254	0	0
219	171	254	2	0
218	170	254	0	0
70	171	248	0	0
208	172	252	2	0
205	171	250	0	0
218	171	255	2	0
70	173	261	0	0
219	171	258	3	0
70	174	256	0	0
219	173	258	5	0
4	174	261	0	0
218	171	256	2	0
70	182	224	0	0
205	176	233	0	0
70	179	235	0	0
70	181	245	0	0
70	176	253	0	0
70	177	253	0	0
4	177	232	0	0
4	176	263	0	0
70	181	255	0	0
205	183	227	0	0
70	179	228	0	0
38	183	255	0	0
70	183	238	0	0
38	180	243	0	0
4	183	252	0	0
38	183	259	0	0
208	178	258	0	0
70	178	263	0	0
70	178	243	0	0
38	177	251	0	0
70	179	217	0	0
4	177	217	0	0
70	176	221	0	0
70	177	223	0	0
205	181	214	0	0
70	176	208	3	0
70	181	209	0	0
4	174	201	0	0
70	175	204	0	0
70	183	203	0	0
70	160	205	0	0
70	180	207	0	0
70	182	207	3	0
70	180	204	3	0
70	177	203	3	0
38	170	200	0	0
70	161	206	0	0
205	171	205	0	0
70	179	200	0	0
70	169	206	0	0
4	168	205	0	0
70	179	206	3	0
70	180	202	3	0
70	177	200	0	0
70	176	206	3	0
70	178	204	3	0
70	190	236	0	0
70	187	224	0	0
70	185	247	0	0
70	188	239	0	0
70	189	237	0	0
70	191	241	0	0
4	187	248	0	0
4	187	249	0	0
70	188	229	0	0
70	189	228	0	0
70	185	250	0	0
70	189	233	0	0
70	184	254	0	0
70	187	230	0	0
70	191	228	0	0
70	187	238	0	0
4	184	250	0	0
70	191	242	0	0
70	184	234	0	0
70	186	235	0	0
70	189	226	0	0
70	187	225	0	0
70	184	222	0	0
38	186	217	0	0
4	187	218	0	0
70	184	217	0	0
205	187	221	0	0
70	191	213	0	0
70	189	210	0	0
70	189	208	0	0
4	187	210	0	0
4	188	203	0	0
38	185	207	0	0
4	187	201	0	0
70	184	206	0	0
70	165	196	0	0
70	154	192	0	0
205	167	193	0	0
70	167	195	0	0
4	179	196	0	0
70	175	192	0	0
70	173	192	0	0
4	182	193	0	0
70	182	199	0	0
70	175	193	0	0
70	176	194	0	0
70	188	198	0	0
70	184	196	0	0
70	176	196	0	0
70	183	196	0	0
70	173	199	0	0
4	180	192	0	0
70	174	198	0	0
70	177	196	0	0
70	180	198	0	0
205	183	197	0	0
4	169	193	0	0
70	171	190	0	0
70	165	185	0	0
70	182	187	0	0
70	167	186	0	0
70	156	191	0	0
4	175	185	0	0
205	175	190	0	0
70	185	188	0	0
205	187	185	0	0
70	174	188	0	0
70	188	189	0	0
70	158	176	0	0
70	152	176	0	0
4	161	177	0	0
70	158	182	0	0
70	175	178	0	0
4	172	181	0	0
4	172	178	0	0
70	182	180	0	0
4	189	176	0	0
38	190	183	0	0
70	177	178	0	0
38	154	168	0	0
205	161	172	0	0
4	189	173	0	0
70	162	172	0	0
4	167	169	0	0
70	186	171	0	0
70	187	172	0	0
70	188	170	0	0
4	164	165	0	0
206	168	162	3	0
70	166	162	0	0
204	170	163	2	0
70	175	167	0	0
70	160	165	0	0
4	183	163	0	0
70	166	165	0	0
38	183	162	0	0
70	159	162	0	0
4	182	167	0	0
206	172	163	6	0
70	157	160	0	0
20	173	160	0	0
70	177	167	0	0
15	169	160	4	0
4	171	167	0	0
24	168	160	2	1
24	172	164	2	1
70	185	154	0	0
70	184	156	0	0
4	169	159	0	0
4	167	157	0	0
70	175	157	0	0
4	175	158	0	0
70	180	155	0	0
4	190	158	0	0
70	164	155	0	0
70	191	155	0	0
70	160	154	0	0
38	159	156	0	0
70	177	157	0	0
70	157	144	0	0
70	153	146	0	0
205	166	146	0	0
70	155	144	0	0
70	166	145	0	0
70	180	144	0	0
70	183	148	0	0
70	183	150	0	0
70	179	146	0	0
38	174	147	0	0
70	172	151	0	0
70	174	145	0	0
4	165	145	0	0
205	190	150	0	0
70	176	145	0	0
70	184	148	0	0
4	176	147	0	0
205	150	150	0	0
38	150	158	0	0
70	145	146	0	0
70	146	149	0	0
70	148	155	0	0
70	144	156	0	0
70	146	153	0	0
70	148	153	0	0
70	145	173	0	0
70	144	165	0	0
70	144	158	0	0
70	148	169	0	0
70	145	170	0	0
4	151	164	0	0
70	144	178	0	0
38	145	166	0	0
4	150	178	0	0
4	146	167	0	0
4	147	161	0	0
70	149	182	0	0
70	151	136	0	0
70	156	141	0	0
20	147	142	0	0
20	159	142	0	0
70	165	137	0	0
70	175	136	0	0
70	174	138	0	0
70	165	136	0	0
20	178	142	0	0
38	160	138	0	0
38	170	138	0	0
20	170	142	0	0
70	180	139	0	0
70	191	137	0	0
70	186	138	0	0
20	186	142	0	0
38	187	137	0	0
70	188	137	0	0
70	160	129	0	0
205	170	132	0	0
70	175	133	0	0
38	153	134	0	0
70	181	133	0	0
70	169	130	0	0
349	180	129	0	0
70	170	128	0	0
70	182	134	0	0
70	189	128	0	0
70	147	129	0	0
70	148	132	0	0
38	151	128	0	0
70	150	130	0	0
70	147	185	0	0
70	148	192	0	0
70	148	199	0	0
4	196	158	0	0
70	199	154	0	0
4	198	152	0	0
70	198	166	0	0
4	193	161	0	0
38	199	165	0	0
70	192	181	0	0
4	199	177	0	0
70	197	178	0	0
38	193	167	0	0
70	196	160	0	0
70	195	156	0	0
4	199	178	0	0
4	198	174	0	0
70	194	158	0	0
4	196	175	0	0
70	192	180	0	0
4	195	184	0	0
70	192	190	0	0
4	197	181	0	0
70	192	175	0	0
70	198	194	0	0
205	196	192	0	0
70	192	197	0	0
70	192	194	0	0
70	194	196	0	0
70	198	202	0	0
205	192	206	0	0
70	193	207	0	0
70	193	215	0	0
70	199	215	0	0
38	196	214	0	0
70	197	212	0	0
38	196	218	0	0
70	195	221	0	0
70	198	217	0	0
70	197	234	0	0
4	194	224	0	0
70	199	238	0	0
70	193	231	0	0
205	192	238	0	0
70	196	251	0	0
70	193	243	0	0
70	194	260	0	0
70	192	244	0	0
205	195	242	0	0
70	192	257	0	0
70	193	258	0	0
70	192	256	0	0
13	173	306	0	0
70	174	306	0	0
12	173	304	0	0
13	169	308	0	0
12	173	308	0	0
13	169	304	0	0
12	169	306	0	0
70	175	313	0	0
70	174	316	0	0
205	174	322	0	0
13	173	319	0	0
13	173	312	0	0
13	171	322	2	0
12	173	314	0	0
13	173	322	0	0
12	171	325	0	0
70	172	332	0	0
70	175	342	0	0
4	168	334	0	0
12	171	328	2	0
70	169	338	0	0
70	168	342	0	0
12	174	325	0	0
70	174	331	0	0
12	169	314	0	0
13	171	318	0	0
13	169	312	0	0
70	173	333	0	0
70	174	335	0	0
12	173	328	0	0
70	170	339	0	0
70	172	337	0	0
70	171	343	0	0
70	172	341	0	0
70	169	347	0	0
70	170	346	0	0
70	171	298	0	0
70	173	301	0	0
70	175	292	0	0
70	173	288	0	0
70	174	284	0	0
70	170	282	0	0
4	170	281	0	0
70	173	287	0	0
205	170	286	0	0
4	171	282	0	0
70	171	283	0	0
70	170	279	0	0
70	172	272	0	0
70	170	277	0	0
70	174	278	0	0
70	173	268	0	0
70	168	270	0	0
70	180	269	0	0
70	181	264	0	0
70	180	267	0	0
70	191	264	0	0
205	185	268	0	0
4	189	265	0	0
70	199	264	0	0
70	197	270	0	0
4	181	278	0	0
38	190	274	0	0
38	178	274	0	0
70	191	278	0	0
205	181	274	0	0
4	176	274	0	0
205	196	276	0	0
70	196	274	0	0
70	192	276	0	0
70	195	273	0	0
4	194	276	0	0
70	183	283	0	0
70	179	284	0	0
4	177	285	0	0
38	185	283	0	0
4	178	280	0	0
4	188	283	0	0
70	199	287	0	0
70	185	286	0	0
70	191	282	0	0
70	191	287	0	0
205	190	283	0	0
70	196	282	0	0
70	189	281	0	0
70	179	290	0	0
70	181	301	0	0
38	177	302	0	0
70	178	298	0	0
38	181	302	0	0
70	182	297	0	0
70	179	296	0	0
4	181	306	0	0
70	176	307	0	0
13	183	310	3	0
70	180	306	0	0
13	181	309	2	0
70	183	305	0	0
12	181	311	3	0
213	176	310	0	0
70	181	307	0	0
205	179	306	0	0
13	179	309	2	0
4	182	318	0	0
13	183	319	6	0
13	181	319	6	0
12	177	319	6	0
12	183	313	2	0
70	179	313	0	0
70	180	320	0	0
12	182	322	6	0
70	176	320	0	0
12	177	325	0	0
13	179	322	6	0
12	183	325	6	0
213	179	324	0	0
12	180	325	6	0
12	177	322	1	0
4	181	329	0	0
70	179	334	0	0
70	181	332	0	0
13	181	328	6	0
12	179	328	6	0
70	180	331	0	0
13	183	328	6	0
70	183	342	0	0
70	178	336	0	0
38	176	341	0	0
70	183	348	0	0
70	177	351	0	0
70	178	348	0	0
70	178	354	0	0
70	180	359	0	0
70	176	354	0	0
70	180	363	0	0
70	178	367	0	0
70	181	362	0	0
70	181	364	0	0
70	179	370	0	0
4	178	372	0	0
205	178	371	0	0
205	180	382	0	0
205	181	381	0	0
70	177	387	0	0
70	181	386	0	0
205	177	392	0	0
70	179	395	0	0
70	179	396	0	0
70	176	397	0	0
70	183	399	0	0
4	181	398	0	0
70	181	405	0	0
70	176	402	0	0
4	179	400	0	0
70	176	405	0	0
70	181	406	0	0
70	176	414	0	0
205	179	411	0	0
70	189	408	0	0
38	186	409	0	0
70	189	411	0	0
70	189	400	0	0
4	186	405	0	0
70	185	403	0	0
70	186	407	0	0
70	191	405	0	0
70	184	401	0	0
70	184	390	0	0
70	186	388	0	0
70	191	384	0	0
4	191	377	0	0
70	191	383	0	0
4	189	383	0	0
38	184	382	0	0
70	187	369	0	0
70	188	368	0	0
70	186	362	0	0
70	185	366	0	0
70	186	365	0	0
70	184	356	0	0
70	187	357	0	0
70	191	356	0	0
70	191	357	0	0
70	187	348	0	0
70	185	347	0	0
4	185	338	0	0
205	187	338	0	0
70	184	339	0	0
70	187	342	0	0
38	189	339	0	0
205	189	342	0	0
70	190	333	0	0
70	187	333	0	0
13	185	325	6	0
12	185	322	0	0
70	189	321	0	0
70	185	320	0	0
70	190	320	0	0
70	186	324	0	0
13	185	313	2	0
70	190	317	0	0
70	186	318	0	0
70	189	316	0	0
12	185	319	6	0
4	187	299	0	0
70	190	297	0	0
205	191	300	0	0
70	184	293	0	0
70	184	294	0	0
70	184	308	0	0
12	185	307	2	0
70	185	304	0	0
4	185	295	0	0
12	187	309	2	0
70	196	295	0	0
70	196	290	0	0
70	198	291	0	0
70	195	290	0	0
70	194	292	0	0
70	194	291	0	0
70	195	289	0	0
4	192	298	0	0
70	194	297	0	0
70	196	303	0	0
70	194	309	0	0
70	192	310	0	0
205	199	304	0	0
70	194	305	0	0
70	196	313	0	0
70	192	313	0	0
70	195	313	0	0
4	197	319	0	0
70	196	323	0	0
4	192	321	0	0
70	192	324	0	0
38	193	335	0	0
38	196	329	0	0
4	195	341	0	0
205	196	339	0	0
70	200	313	0	0
38	202	308	0	0
70	207	321	0	0
38	203	320	0	0
70	201	335	0	0
70	203	340	0	0
70	205	323	0	0
70	202	333	0	0
4	200	308	0	0
70	200	326	0	0
70	207	333	0	0
70	200	324	0	0
38	205	338	0	0
70	207	307	0	0
70	204	314	0	0
4	205	302	0	0
70	207	300	0	0
70	204	298	0	0
70	203	297	0	0
4	203	301	0	0
70	203	296	0	0
70	200	348	0	0
70	199	348	0	0
70	200	346	0	0
70	201	345	0	0
70	197	350	0	0
70	200	351	0	0
70	206	348	0	0
70	195	350	0	0
70	198	359	0	0
70	199	365	0	0
38	199	361	0	0
70	205	358	0	0
70	207	362	0	0
70	204	355	0	0
70	202	356	0	0
70	207	354	0	0
205	207	360	0	0
70	207	365	0	0
205	200	364	0	0
70	202	360	0	0
70	204	364	0	0
38	200	367	0	0
70	203	365	0	0
70	195	368	0	0
70	203	370	0	0
70	206	372	0	0
205	201	371	0	0
4	193	378	0	0
70	198	380	0	0
70	192	378	0	0
70	197	391	0	0
38	198	389	0	0
70	199	395	0	0
70	196	406	0	0
70	198	402	0	0
70	192	404	0	0
70	195	403	0	0
4	194	409	0	0
38	206	410	0	0
70	206	407	0	0
4	202	400	0	0
70	202	404	0	0
70	200	403	0	0
70	205	401	0	0
70	207	386	0	0
70	206	388	0	0
70	201	389	0	0
70	206	391	0	0
70	203	391	0	0
38	203	386	0	0
70	207	384	0	0
205	203	376	0	0
70	203	378	0	0
70	201	378	0	0
70	207	291	0	0
70	205	293	0	0
70	202	285	0	0
70	207	275	0	0
70	207	278	0	0
70	206	275	0	0
70	205	275	0	0
205	205	278	0	0
205	202	275	0	0
70	205	277	0	0
70	207	268	0	0
70	204	266	0	0
70	206	271	0	0
4	207	266	0	0
4	207	258	0	0
4	205	259	0	0
70	203	258	0	0
205	201	261	0	0
4	200	258	0	0
70	203	251	0	0
70	205	249	0	0
70	201	253	0	0
70	207	253	0	0
70	200	247	0	0
70	206	246	0	0
70	204	238	0	0
70	202	229	0	0
70	207	224	0	0
70	201	226	0	0
70	204	222	0	0
70	203	219	0	0
70	201	222	0	0
70	203	214	0	0
4	205	208	0	0
70	207	215	0	0
70	206	201	0	0
70	206	202	0	0
4	203	201	0	0
70	201	203	0	0
38	201	206	0	0
70	205	195	0	0
38	206	193	0	0
70	200	197	0	0
70	201	184	0	0
70	203	189	0	0
38	204	190	0	0
70	201	190	0	0
70	202	180	0	0
70	201	179	0	0
70	207	173	0	0
70	205	168	0	0
70	200	171	0	0
70	203	160	0	0
70	206	167	0	0
70	206	161	0	0
205	203	155	0	0
70	204	155	0	0
70	195	148	0	0
70	204	151	0	0
205	206	149	0	0
70	198	138	0	0
20	194	142	0	0
205	202	140	0	0
205	194	139	0	0
70	207	138	0	0
70	200	138	0	0
205	205	133	0	0
205	200	128	0	0
70	203	133	0	0
205	208	135	0	0
4	214	158	0	0
70	214	153	0	0
70	210	148	0	0
20	208	142	0	0
70	212	157	0	0
205	208	145	0	0
4	208	158	0	0
38	213	134	0	0
70	213	163	0	0
205	212	150	0	0
70	209	147	0	0
70	210	171	0	0
70	214	174	0	0
38	214	168	0	0
70	208	174	0	0
70	211	186	0	0
70	215	198	0	0
70	209	199	0	0
70	209	196	0	0
70	215	203	0	0
70	214	206	0	0
70	212	204	0	0
4	215	214	0	0
205	215	208	0	0
70	210	211	0	0
4	215	215	0	0
70	210	222	0	0
70	214	220	0	0
70	214	219	0	0
70	214	225	0	0
70	208	230	0	0
205	215	228	0	0
70	208	234	0	0
70	214	235	0	0
70	210	235	0	0
70	215	242	0	0
70	210	241	0	0
70	209	242	0	0
70	208	250	0	0
70	215	248	0	0
70	209	254	0	0
70	213	251	0	0
205	210	261	0	0
70	208	259	0	0
70	213	265	0	0
4	211	265	0	0
38	210	266	0	0
70	208	273	0	0
70	211	282	0	0
38	215	287	0	0
4	210	291	0	0
205	211	292	0	0
70	215	292	0	0
70	215	301	0	0
70	209	299	0	0
70	211	296	0	0
205	213	311	0	0
70	209	304	0	0
70	212	309	0	0
70	211	308	0	0
70	209	318	0	0
70	214	315	0	0
70	214	316	0	0
70	212	313	0	0
70	215	323	0	0
70	210	326	0	0
70	208	322	0	0
70	211	326	0	0
70	213	329	0	0
70	213	331	0	0
70	211	330	0	0
70	208	335	0	0
70	211	331	0	0
70	209	409	0	0
70	211	406	0	0
38	214	393	0	0
4	221	399	0	0
70	223	397	0	0
4	218	415	0	0
70	220	416	0	0
70	223	400	0	0
209	223	429	7	0
70	219	419	0	0
209	219	429	3	0
70	220	414	0	0
4	217	397	0	0
209	221	429	6	0
70	223	402	0	0
38	222	404	0	0
4	217	416	0	0
103	221	395	0	0
70	222	407	0	0
70	222	417	0	0
111	220	397	0	0
38	217	405	0	0
70	220	403	0	0
70	220	393	0	0
209	216	429	2	0
38	216	426	0	0
111	219	397	0	0
70	218	389	0	0
70	217	386	0	0
70	217	388	0	0
70	212	382	0	0
70	214	382	0	0
38	209	381	0	0
205	209	378	0	0
70	209	377	0	0
70	209	372	0	0
70	210	374	0	0
70	208	375	0	0
70	217	363	0	0
70	223	364	0	0
70	220	367	0	0
205	220	368	0	0
70	210	360	0	0
70	221	378	0	0
70	214	359	0	0
70	216	353	0	0
4	214	353	0	0
70	217	352	0	0
70	209	358	0	0
38	218	354	0	0
70	209	348	0	0
70	211	349	0	0
70	216	350	0	0
70	223	343	0	0
70	211	340	0	0
70	221	342	0	0
70	213	341	0	0
4	212	337	0	0
4	214	340	0	0
70	219	341	0	0
70	211	339	0	0
1010	217	135	6	0
1018	217	132	2	0
1010	217	133	6	0
1012	217	134	2	0
1020	219	129	4	0
1009	217	129	2	0
1010	223	141	0	0
1012	217	130	2	0
1009	221	138	4	0
1018	217	128	2	0
1011	220	139	4	0
70	217	156	0	0
20	219	142	0	0
4	221	154	0	0
38	219	157	0	0
1013	220	138	4	0
70	220	153	0	0
1011	222	141	4	0
1013	222	140	4	0
1009	217	131	2	0
1016	221	140	4	0
70	223	159	0	0
1013	218	136	4	0
1013	221	139	4	0
1013	219	137	4	0
1016	218	137	4	0
4	218	150	0	0
70	221	146	0	0
70	219	167	0	0
1011	217	136	4	0
70	218	166	0	0
1011	219	138	4	0
70	216	165	0	0
70	216	161	0	0
70	223	173	0	0
70	223	170	0	0
70	222	172	0	0
70	221	182	0	0
70	220	183	0	0
70	222	176	0	0
70	218	181	0	0
38	216	178	0	0
70	217	180	0	0
70	216	180	0	0
70	220	190	0	0
70	223	185	0	0
70	219	190	0	0
70	223	199	0	0
70	218	199	0	0
70	218	194	0	0
4	221	195	0	0
70	217	197	0	0
70	217	199	0	0
205	218	202	0	0
4	217	202	0	0
4	217	209	0	0
70	223	218	0	0
70	218	223	0	0
70	217	221	0	0
70	219	221	0	0
70	221	222	0	0
70	220	228	0	0
70	223	225	0	0
70	218	224	0	0
4	220	239	0	0
70	219	236	0	0
108	223	247	4	0
4	222	243	0	0
70	217	246	0	0
107	220	247	0	0
70	220	244	0	0
4	218	242	0	0
70	222	242	0	0
107	222	249	0	0
110	220	252	4	0
107	222	250	0	0
111	219	255	0	0
106	222	252	4	0
106	223	253	4	0
107	222	254	0	0
106	216	249	4	0
107	222	253	0	0
107	219	249	0	0
107	218	258	0	0
111	222	257	0	0
70	223	268	0	0
70	223	267	0	0
70	220	267	0	0
70	220	265	0	0
205	222	274	0	0
70	220	277	0	0
70	219	287	0	0
70	217	282	0	0
70	223	282	0	0
70	222	285	0	0
38	221	292	0	0
70	222	290	0	0
70	222	293	0	0
4	219	300	0	0
70	218	298	0	0
4	218	296	0	0
70	217	297	0	0
4	219	309	0	0
70	220	316	0	0
4	223	319	0	0
205	221	312	0	0
38	216	319	0	0
4	223	321	0	0
70	220	326	0	0
70	220	334	0	0
70	216	332	0	0
38	223	334	0	0
6	231	392	4	0
38	230	420	0	0
303	231	394	4	0
205	225	409	0	0
205	230	418	0	0
205	224	410	0	0
70	229	398	0	0
70	225	387	0	0
70	224	389	0	0
70	226	391	0	0
70	230	405	0	0
70	227	394	0	0
205	229	419	0	0
205	225	402	0	0
38	229	404	0	0
70	227	422	0	0
70	225	422	0	0
205	227	421	0	0
38	231	379	0	0
70	229	378	0	0
4	225	373	0	0
70	224	369	0	0
38	228	366	0	0
70	231	367	0	0
70	226	367	0	0
70	228	354	0	0
38	226	355	0	0
205	228	357	0	0
70	226	344	0	0
4	225	344	0	0
70	231	346	0	0
70	230	331	0	0
70	230	328	0	0
70	228	329	0	0
70	224	334	0	0
4	230	335	0	0
70	227	330	0	0
70	229	328	0	0
205	226	335	0	0
70	226	325	0	0
70	225	325	0	0
205	231	322	0	0
70	227	325	0	0
70	227	326	0	0
70	229	281	0	0
70	230	265	0	0
70	228	269	0	0
108	224	258	4	0
110	226	256	4	0
110	226	259	4	0
70	226	288	0	0
70	225	294	0	0
109	231	256	0	0
70	227	284	0	0
106	229	259	4	0
102	226	253	4	0
102	225	249	4	0
111	231	248	0	0
102	228	253	4	0
108	229	249	4	0
70	231	243	0	0
110	230	246	4	0
38	230	242	0	0
106	229	246	4	0
111	231	246	0	0
111	229	244	0	0
70	228	247	0	0
70	226	243	0	0
70	225	241	0	0
70	231	238	0	0
70	230	222	0	0
70	230	229	0	0
205	226	228	0	0
4	227	231	0	0
38	230	221	0	0
70	230	212	0	0
70	230	214	0	0
205	228	213	0	0
70	230	209	0	0
4	227	213	0	0
70	227	211	0	0
4	229	213	0	0
4	230	211	0	0
205	231	215	0	0
70	229	210	0	0
38	225	203	0	0
4	229	205	0	0
70	229	185	0	0
70	227	183	0	0
70	225	179	0	0
70	225	180	0	0
205	226	178	0	0
70	227	179	0	0
24	226	180	0	1
70	230	173	0	0
70	224	172	0	0
70	227	166	0	0
1010	231	141	0	0
1018	230	141	4	0
1009	229	141	4	0
1016	224	133	4	0
1012	228	141	4	0
1010	225	141	0	0
20	227	142	0	0
1009	227	141	4	0
1012	224	141	4	0
1018	226	141	4	0
70	226	147	0	0
70	230	145	0	0
205	226	144	0	0
4	229	154	0	0
70	227	159	0	0
1018	239	128	2	0
1009	236	138	4	0
1016	232	133	6	0
1019	237	129	0	0
1010	239	134	2	0
1010	232	141	0	0
1012	239	135	6	0
1011	239	137	6	0
70	235	165	0	0
70	234	144	0	0
1010	239	133	2	0
1013	237	138	6	0
1013	236	139	6	0
70	233	160	0	0
1016	236	140	6	0
70	238	150	0	0
38	235	166	0	0
4	235	160	0	0
1011	235	141	6	0
1009	239	129	2	0
1011	237	139	6	0
1018	239	132	2	0
1009	239	131	2	0
1012	239	130	6	0
1013	238	137	6	0
1018	239	136	2	0
1016	238	138	6	0
1010	234	141	0	0
1012	233	141	4	0
70	234	150	0	0
70	234	161	0	0
1013	235	140	6	0
70	237	146	0	0
70	237	167	0	0
205	239	159	0	0
4	235	169	0	0
4	236	171	0	0
222	233	180	0	0
70	236	176	0	0
70	237	191	0	0
70	236	190	0	0
70	239	188	0	0
70	237	189	0	0
70	239	193	0	0
70	238	200	0	0
4	237	208	0	0
70	237	199	0	0
38	234	221	0	0
70	239	223	0	0
70	234	207	0	0
70	233	198	0	0
4	239	204	0	0
70	236	193	0	0
70	236	195	0	0
4	237	214	0	0
24	236	194	1	1
205	232	200	0	0
70	235	195	0	0
38	239	208	0	0
4	234	205	0	0
70	235	193	0	0
70	233	229	0	0
70	237	227	0	0
70	234	238	0	0
38	235	235	0	0
4	232	234	0	0
70	236	240	0	0
103	234	246	0	0
205	232	246	0	0
111	232	249	0	0
70	232	254	0	0
106	233	250	4	0
70	233	253	0	0
205	236	258	0	0
205	232	262	0	0
70	236	269	0	0
70	233	268	0	0
70	239	269	0	0
205	238	265	0	0
70	237	270	0	0
4	239	266	0	0
70	233	277	0	0
70	232	275	0	0
70	237	274	0	0
205	235	280	0	0
70	236	287	0	0
70	232	280	0	0
4	238	293	0	0
70	237	291	0	0
70	235	298	0	0
70	233	303	0	0
38	228	301	0	0
70	237	299	0	0
70	230	297	0	0
70	228	303	0	0
70	230	301	0	0
70	235	307	0	0
70	225	309	0	0
70	230	307	0	0
70	232	315	0	0
70	234	312	0	0
70	233	314	0	0
70	235	312	0	0
70	235	316	0	0
38	228	312	0	0
38	234	327	0	0
70	237	325	0	0
4	237	327	0	0
70	239	328	0	0
70	234	337	0	0
38	239	337	0	0
70	232	337	0	0
205	236	340	0	0
70	236	341	0	0
70	236	350	0	0
70	237	347	0	0
70	235	351	0	0
209	230	430	1	0
70	231	427	0	0
70	225	429	0	0
70	225	426	0	0
70	226	424	0	0
70	224	425	0	0
209	227	430	7	0
0	222	433	4	0
70	227	433	4	0
70	223	436	4	0
34	217	437	4	0
0	216	433	4	0
5	226	439	0	0
4	237	403	0	0
70	236	407	0	0
205	239	420	0	0
70	232	411	0	0
70	236	414	0	0
70	235	418	0	0
209	232	430	0	0
70	234	422	0	0
70	232	431	0	0
0	234	437	4	0
70	232	422	0	0
70	235	428	0	0
70	238	429	0	0
70	233	417	0	0
70	240	420	0	0
70	242	428	0	0
70	240	408	0	0
70	246	425	0	0
4	245	417	0	0
70	241	415	0	0
4	246	403	0	0
70	247	420	0	0
205	243	423	0	0
70	240	416	0	0
70	246	402	0	0
38	245	419	0	0
70	244	407	0	0
70	244	413	0	0
70	247	401	0	0
34	223	447	4	0
29	223	440	0	0
2	222	445	0	0
63	217	447	2	0
34	227	444	4	0
15	245	444	0	0
7	244	442	4	0
54	218	442	0	0
1	242	440	0	0
34	224	446	4	0
139	219	447	2	0
1	225	444	0	1
1	242	443	1	1
70	236	394	0	0
70	242	398	0	0
70	245	397	0	0
38	246	399	0	0
205	239	399	0	0
4	245	399	0	0
70	247	396	0	0
4	233	389	0	0
70	236	391	0	0
70	247	376	0	0
70	236	383	0	0
205	238	380	0	0
70	238	379	0	0
38	253	399	0	0
70	254	376	0	0
70	249	382	0	0
70	252	376	0	0
70	249	392	0	0
70	254	402	0	0
38	249	394	0	0
70	248	409	0	0
70	252	409	0	0
70	249	404	0	0
70	242	368	0	0
70	246	372	0	0
70	251	372	0	0
205	250	374	0	0
38	253	375	0	0
70	252	370	0	0
38	244	364	0	0
70	238	367	0	0
4	234	356	0	0
205	241	358	0	0
70	235	357	0	0
70	238	359	0	0
1	250	357	0	1
1	252	358	1	1
70	242	350	0	0
70	241	347	0	0
70	245	350	0	0
4	241	348	0	0
70	243	348	0	0
70	245	341	0	0
70	240	338	0	0
38	241	339	0	0
4	244	340	0	0
205	240	341	0	0
70	245	332	0	0
70	242	328	0	0
4	241	335	0	0
70	247	325	0	0
4	242	313	0	0
70	240	317	0	0
70	244	304	0	0
70	247	307	0	0
70	243	309	0	0
70	240	305	0	0
205	243	308	0	0
70	247	298	0	0
70	243	296	0	0
70	243	299	0	0
4	242	295	0	0
38	241	291	0	0
4	242	280	0	0
205	246	283	0	0
38	243	285	0	0
70	247	280	0	0
70	247	286	0	0
70	242	277	0	0
4	244	279	0	0
4	240	277	0	0
70	247	277	0	0
70	240	268	0	0
38	246	269	0	0
70	246	266	0	0
38	240	265	0	0
70	247	259	0	0
4	240	256	0	0
4	242	263	0	0
70	243	255	0	0
70	247	250	0	0
70	244	250	0	0
70	241	244	0	0
70	242	234	0	0
205	243	228	0	0
38	244	230	0	0
4	244	231	0	0
70	245	231	0	0
70	241	231	0	0
70	240	219	0	0
70	240	220	0	0
4	244	223	0	0
38	244	222	0	0
205	244	212	0	0
70	247	208	0	0
4	242	210	0	0
70	240	208	0	0
70	242	207	0	0
70	245	207	0	0
70	244	203	0	0
4	243	198	0	0
70	246	192	0	0
4	244	193	0	0
38	242	186	0	0
4	244	191	0	0
70	252	210	0	0
4	250	197	0	0
70	254	191	0	0
70	249	187	0	0
70	252	199	0	0
70	250	196	0	0
70	248	223	0	0
70	249	209	0	0
70	252	204	0	0
70	248	201	0	0
38	249	206	0	0
4	249	202	0	0
4	250	204	0	0
70	252	225	0	0
4	253	231	0	0
70	254	231	0	0
70	252	236	0	0
4	248	234	0	0
205	248	233	0	0
4	249	241	0	0
4	253	241	0	0
70	251	242	0	0
38	248	244	0	0
205	250	252	0	0
70	251	252	0	0
38	254	263	0	0
4	253	262	0	0
70	253	263	0	0
70	250	261	0	0
4	253	265	0	0
70	253	268	0	0
205	252	268	0	0
4	248	265	0	0
38	251	265	0	0
70	252	279	0	0
70	250	286	0	0
70	251	282	0	0
205	252	285	0	0
70	253	282	0	0
193	254	290	4	0
70	248	301	0	0
70	254	301	0	0
4	252	309	0	0
70	248	304	0	0
70	252	304	0	0
38	254	319	0	0
38	252	318	0	0
70	255	316	0	0
70	252	314	0	0
70	252	326	0	0
4	250	326	0	0
70	250	325	0	0
70	250	327	0	0
70	249	329	0	0
70	251	330	0	0
205	254	328	0	0
70	251	333	0	0
70	255	337	0	0
70	253	337	0	0
38	253	340	0	0
70	252	340	0	0
15	254	345	4	0
1	252	347	1	1
1	250	349	0	1
38	254	417	0	0
38	252	423	0	0
38	252	424	0	0
20	248	440	6	0
1	248	442	0	0
70	261	412	0	0
4	262	413	0	0
205	261	410	0	0
70	263	419	0	0
70	256	412	0	0
4	259	421	0	0
70	257	424	0	0
70	260	423	0	0
70	262	424	0	0
38	263	431	0	0
70	262	429	0	0
70	263	429	0	0
4	257	403	0	0
70	262	406	0	0
4	261	396	0	0
70	261	397	0	0
70	262	387	0	0
70	262	389	0	0
70	261	391	0	0
70	259	387	0	0
70	256	385	0	0
38	258	370	0	0
70	257	368	0	0
70	261	370	0	0
4	257	367	0	0
1	262	357	0	1
1	261	358	1	1
70	267	364	0	0
205	270	356	0	0
70	270	365	0	0
205	267	362	0	0
70	267	361	0	0
70	266	370	0	0
205	269	372	0	0
111	270	381	0	0
70	271	381	0	0
4	264	389	0	0
4	266	373	0	0
111	271	376	0	0
70	265	381	0	0
70	266	381	0	0
111	266	378	0	0
111	270	377	0	0
4	269	379	0	0
1	264	353	1	1
15	259	345	4	0
1	261	347	1	1
38	259	338	0	0
70	263	341	0	0
70	256	332	0	0
70	262	328	0	0
70	261	335	0	0
70	259	335	0	0
205	258	334	0	0
70	261	322	0	0
70	261	317	0	0
4	260	319	0	0
70	257	316	0	0
70	259	317	0	0
205	263	310	0	0
205	263	311	0	0
1	261	304	1	1
70	256	290	0	0
193	259	294	4	0
70	262	291	0	0
70	263	280	0	0
70	260	287	0	0
4	257	273	0	0
4	263	278	0	0
4	259	277	0	0
70	257	276	0	0
205	258	268	0	0
4	264	282	0	0
70	268	280	0	0
70	266	273	0	0
38	267	282	0	0
70	269	267	0	0
97	267	298	6	0
4	270	303	0	0
97	271	296	6	0
70	268	297	0	0
97	266	302	6	0
70	264	289	0	0
70	265	293	0	0
1	267	290	1	1
1	270	293	0	1
1	265	302	1	1
70	259	260	0	0
4	271	262	0	0
70	264	257	0	0
70	257	263	0	0
70	271	257	0	0
70	259	248	0	0
70	261	242	0	0
218	270	247	1	0
70	257	244	0	0
218	271	246	1	0
70	267	244	0	0
4	256	247	0	0
70	268	244	0	0
70	260	233	0	0
70	266	232	0	0
70	260	237	0	0
70	271	237	0	0
70	268	236	0	0
70	256	239	0	0
70	257	225	0	0
4	256	229	0	0
70	268	225	0	0
70	269	224	0	0
4	260	229	0	0
70	261	225	0	0
38	258	227	0	0
70	266	231	0	0
70	260	227	0	0
205	265	225	0	0
70	258	221	0	0
38	262	223	0	0
70	271	219	0	0
70	270	216	0	0
70	264	216	0	0
4	267	223	0	0
4	261	210	0	0
70	262	214	0	0
4	271	210	0	0
70	268	208	0	0
70	260	210	0	0
38	263	211	0	0
4	259	208	0	0
70	257	211	0	0
4	268	201	0	0
70	260	206	0	0
70	261	207	0	0
4	264	202	0	0
70	259	202	0	0
70	257	196	0	0
38	268	193	0	0
70	266	198	0	0
70	260	192	0	0
4	264	192	0	0
70	265	184	0	0
70	268	185	0	0
70	267	187	0	0
70	261	189	0	0
70	273	187	0	0
70	274	186	0	0
70	276	189	0	0
70	272	184	0	0
70	273	186	0	0
4	277	204	0	0
4	276	222	0	0
205	275	203	0	0
38	273	216	0	0
205	278	216	0	0
70	274	215	0	0
70	272	199	0	0
70	277	194	0	0
70	272	201	0	0
70	272	209	0	0
38	272	215	0	0
38	275	222	0	0
70	274	196	0	0
24	273	185	3	1
70	244	183	0	0
70	250	176	0	0
70	248	178	0	0
70	260	183	0	0
6	245	180	2	0
4	258	183	0	0
70	260	180	0	0
319	243	178	2	0
70	275	183	0	0
70	281	177	0	0
70	287	187	0	0
70	281	180	0	0
508	285	185	0	0
38	283	178	0	0
4	284	213	0	0
38	283	194	0	0
70	282	197	0	0
70	284	203	0	0
4	283	195	0	0
70	286	197	0	0
38	282	188	0	0
6	282	185	0	0
70	284	201	0	0
38	282	207	0	0
70	287	218	0	0
205	281	218	0	0
70	281	220	0	0
70	284	216	0	0
70	281	222	0	0
70	277	224	0	0
205	275	231	0	0
4	275	224	0	0
70	283	227	0	0
70	282	224	0	0
4	282	229	0	0
70	282	227	0	0
70	273	237	0	0
70	285	239	0	0
38	273	238	0	0
70	272	233	0	0
70	276	238	0	0
70	284	238	0	0
70	283	233	0	0
70	279	243	0	0
219	278	245	2	0
219	276	245	6	0
38	272	244	0	0
205	275	247	4	0
205	281	244	4	0
70	275	244	0	0
205	277	241	4	0
205	272	243	4	0
220	276	243	4	0
70	280	246	0	0
4	286	245	0	0
219	278	243	0	0
218	276	244	6	0
219	275	243	6	0
205	274	246	4	0
218	278	242	2	0
218	277	245	4	0
70	284	240	0	0
70	284	246	0	0
70	277	246	0	0
4	274	242	0	0
70	282	245	0	0
70	286	244	0	0
219	278	252	0	0
70	277	252	1	0
70	279	251	1	0
70	277	249	0	0
218	278	251	6	0
220	275	249	0	0
220	276	255	2	0
220	275	255	4	0
70	275	251	1	0
70	285	251	0	0
70	273	248	1	0
70	281	248	0	0
70	272	250	1	0
70	287	250	0	0
219	273	249	4	0
70	280	249	0	0
70	272	254	0	0
218	274	249	0	0
38	281	254	0	0
4	278	260	0	0
70	280	260	0	0
218	277	256	4	0
220	278	256	2	0
205	284	262	0	0
220	278	257	2	0
4	276	258	0	0
70	280	257	0	0
220	276	256	0	0
70	283	260	0	0
70	280	259	0	0
70	277	265	0	0
38	275	270	0	0
70	273	269	0	0
70	274	272	0	0
70	278	279	0	0
4	279	279	0	0
70	278	278	0	0
205	274	274	0	0
70	280	268	0	0
205	286	265	0	0
4	280	269	0	0
4	282	267	0	0
70	284	268	0	0
70	280	266	0	0
205	280	274	0	0
70	286	272	0	0
70	281	272	0	0
70	282	278	0	0
38	280	275	0	0
70	281	282	0	0
4	285	282	0	0
70	287	286	0	0
4	272	286	0	0
70	272	282	0	0
70	280	281	0	0
70	282	287	0	0
70	273	286	0	0
11	273	291	4	0
70	276	290	0	0
55	278	293	0	0
3	272	289	0	0
38	276	289	0	0
3	281	293	0	0
3	273	289	0	0
70	277	291	0	0
1	280	293	0	1
1	278	295	1	1
55	278	298	0	0
55	278	297	0	0
97	276	303	6	0
97	274	297	6	0
70	276	299	0	0
70	283	297	0	0
3	281	298	0	0
1	279	299	0	1
70	288	286	0	0
205	291	293	0	0
70	292	285	0	0
205	294	296	0	0
4	292	271	0	0
70	291	273	0	0
4	293	276	0	0
4	294	286	0	0
70	293	281	0	0
38	294	290	0	0
70	291	288	0	0
70	289	291	0	0
205	299	282	0	0
70	297	283	0	0
70	298	287	0	0
70	303	286	0	0
70	300	270	0	0
70	301	282	0	0
70	300	291	0	0
4	303	290	0	0
70	297	269	0	0
70	301	292	0	0
70	299	301	0	0
70	300	274	0	0
70	301	273	0	0
70	301	296	0	0
70	298	272	0	0
70	296	292	0	0
50	303	298	0	0
70	298	273	0	0
117	302	299	2	0
70	303	301	0	0
205	300	298	0	0
97	267	305	6	0
4	276	309	0	0
70	275	306	0	0
4	276	305	0	0
97	271	307	6	0
205	293	308	0	0
38	291	306	0	0
70	301	309	0	0
70	273	310	0	0
205	303	309	0	0
205	295	305	0	0
205	284	310	0	0
4	272	305	0	0
4	285	311	0	0
70	297	308	0	0
4	299	304	0	0
38	296	309	0	0
4	269	313	0	0
70	280	318	0	0
70	283	316	0	0
70	279	314	0	0
70	277	316	0	0
70	267	312	0	0
38	294	319	0	0
70	299	312	0	0
70	278	316	0	0
70	296	312	0	0
38	295	314	0	0
205	275	325	0	0
70	296	327	0	0
70	264	320	0	0
70	299	325	0	0
4	291	322	0	0
70	283	324	0	0
70	272	323	0	0
205	295	320	0	0
70	303	335	0	0
70	302	331	0	0
70	267	329	0	0
70	284	331	0	0
70	287	334	0	0
70	282	329	0	0
70	266	334	0	0
70	293	330	0	0
4	280	329	0	0
70	267	330	0	0
70	292	329	0	0
70	265	330	0	0
70	265	335	0	0
70	292	334	0	0
70	276	329	0	0
70	279	340	0	0
70	287	343	0	0
205	270	343	0	0
70	283	341	0	0
70	272	336	0	0
70	292	336	0	0
70	295	340	0	0
70	289	341	0	0
4	293	337	0	0
4	283	336	0	0
38	281	337	0	0
70	282	338	0	0
70	299	343	0	0
70	270	344	0	0
70	289	350	0	0
70	274	346	0	0
38	288	349	0	0
70	270	345	0	0
70	272	349	0	0
70	281	347	0	0
38	273	348	0	0
4	291	348	0	0
70	284	354	0	0
38	282	352	0	0
70	274	359	0	0
70	291	353	0	0
70	288	356	0	0
70	281	365	0	0
70	295	360	0	0
70	286	360	0	0
4	290	364	0	0
70	282	361	0	0
70	295	366	0	0
205	293	365	0	0
70	293	360	0	0
205	274	366	0	0
4	274	364	0	0
70	290	367	0	0
70	273	367	0	0
70	285	360	0	0
70	298	345	0	0
70	302	358	0	0
70	299	357	0	0
70	298	357	0	0
70	296	355	0	0
70	301	361	0	0
70	302	367	0	0
110	278	375	4	0
111	276	375	0	0
110	282	369	4	0
110	282	373	0	0
70	289	375	0	0
111	274	374	0	0
111	273	371	0	0
70	287	369	0	0
110	276	369	4	0
110	279	373	4	0
70	288	373	0	0
70	288	375	0	0
110	277	377	4	0
111	272	378	0	0
111	273	381	0	0
110	275	377	0	0
110	276	378	4	0
110	280	380	0	0
110	280	377	4	0
70	292	383	0	0
4	290	382	0	0
110	284	382	4	0
110	279	382	4	0
111	272	377	0	0
110	276	382	4	0
110	284	378	0	0
110	286	379	4	0
110	278	379	4	0
70	290	378	0	0
70	285	386	0	0
70	275	390	0	0
70	287	387	0	0
4	292	384	0	0
70	285	387	0	0
70	283	384	0	0
70	291	387	0	0
70	276	399	0	0
70	288	393	0	0
70	279	394	0	0
38	271	396	0	0
205	290	394	0	0
4	292	393	0	0
4	265	394	0	0
70	291	392	0	0
38	273	393	0	0
4	290	395	0	0
38	267	397	0	0
70	295	392	0	0
70	277	395	0	0
70	269	394	0	0
70	273	396	0	0
70	292	392	0	0
70	267	395	0	0
70	273	404	0	0
70	284	400	0	0
70	265	403	0	0
205	295	403	0	0
70	289	405	0	0
205	288	404	0	0
70	270	414	0	0
70	270	409	0	0
205	281	414	0	0
70	268	412	0	0
70	275	415	0	0
4	276	414	0	0
70	280	410	0	0
70	278	411	0	0
205	278	412	0	0
70	265	411	0	0
70	267	412	0	0
70	290	414	0	0
38	264	421	0	0
70	289	420	0	0
205	291	421	0	0
70	292	417	0	0
70	285	417	0	0
70	286	423	0	0
70	268	419	0	0
70	272	420	0	0
70	284	416	0	0
4	293	419	0	0
70	283	421	0	0
70	283	422	0	0
70	296	390	0	0
70	300	408	0	0
70	299	385	0	0
70	300	399	0	0
70	298	385	0	0
70	299	411	0	0
70	299	415	0	0
70	297	411	0	0
70	300	411	0	0
70	303	410	0	0
70	300	394	0	0
70	299	403	0	0
4	299	389	0	0
70	301	422	0	0
70	297	418	0	0
70	303	375	0	0
4	300	371	0	0
70	297	368	0	0
70	296	375	0	0
70	303	373	0	0
4	299	368	0	0
70	300	374	0	0
70	296	261	0	0
70	299	261	0	0
70	299	260	0	0
4	295	258	0	0
38	292	262	0	0
70	292	249	0	0
70	295	254	0	0
4	291	255	0	0
70	288	251	0	0
4	291	254	0	0
70	300	247	0	0
70	300	240	0	0
4	289	242	0	0
70	288	246	0	0
70	292	242	0	0
70	299	240	0	0
4	302	234	0	0
70	299	234	0	0
70	303	232	0	0
70	301	238	0	0
70	289	239	0	0
70	289	238	0	0
70	298	236	0	0
38	301	235	0	0
70	297	239	0	0
70	303	227	0	0
70	290	227	0	0
70	293	229	0	0
4	298	220	0	0
70	297	219	0	0
70	296	219	0	0
70	302	216	0	0
70	293	215	0	0
70	300	211	0	0
70	302	210	0	0
70	300	212	0	0
205	289	210	0	0
70	293	213	0	0
70	294	201	0	0
70	291	204	0	0
70	289	205	0	0
4	297	200	0	0
70	288	206	0	0
4	300	201	0	0
70	295	196	0	0
70	299	192	0	0
70	288	195	0	0
70	300	196	0	0
70	297	196	0	0
4	303	196	0	0
70	293	188	0	0
70	295	184	0	0
70	295	188	0	0
70	295	177	0	0
4	303	182	0	0
70	297	178	0	0
70	289	182	0	0
4	293	179	0	0
4	300	178	0	0
38	302	177	0	0
38	300	180	0	0
70	303	177	0	0
70	292	181	0	0
205	267	169	0	0
70	265	170	0	0
70	279	173	0	0
70	279	172	0	0
70	283	169	0	0
70	282	168	0	0
4	281	173	0	0
70	282	172	0	0
70	288	173	0	0
70	283	168	0	0
205	297	174	0	0
205	303	171	0	0
70	302	175	0	0
70	284	164	0	0
4	266	167	0	0
4	272	166	0	0
4	286	164	0	0
70	265	160	0	0
205	281	167	0	0
70	288	167	0	0
70	300	165	0	0
70	303	162	0	0
38	279	154	0	0
70	268	156	0	0
70	266	159	0	0
205	279	152	0	0
70	277	158	0	0
4	284	157	0	0
205	287	159	0	0
205	275	155	0	0
70	272	158	0	0
70	278	152	0	0
205	276	152	0	0
4	285	157	0	0
70	283	158	0	0
4	289	157	0	0
205	277	150	0	0
70	278	147	0	0
205	265	144	0	0
70	287	144	0	0
205	280	146	0	0
70	301	146	0	0
205	292	144	0	0
4	300	144	0	0
70	291	149	0	0
70	296	144	0	0
70	301	148	0	0
70	288	151	0	0
70	302	149	0	0
70	288	144	0	0
38	279	140	0	0
20	283	142	0	0
20	292	142	0	0
70	271	137	0	0
205	303	140	0	0
70	299	137	0	0
20	302	142	0	0
70	296	139	0	0
20	275	142	0	0
42	268	128	4	0
20	270	132	0	0
70	301	132	0	0
20	267	132	6	0
38	283	129	0	0
205	292	130	0	0
70	274	133	0	0
70	291	133	0	0
703	297	134	6	0
70	282	132	0	0
70	265	120	0	0
20	270	127	2	0
70	286	124	0	0
70	287	127	0	0
20	267	127	4	0
704	297	125	2	0
70	299	123	0	0
70	303	123	0	0
70	294	124	0	0
205	283	120	0	0
205	310	133	0	0
70	311	126	0	0
205	308	134	0	0
205	304	130	0	0
70	306	134	0	0
70	308	149	0	0
70	309	146	0	0
70	308	138	0	0
4	308	162	0	0
4	308	169	0	0
70	304	168	0	0
205	305	170	0	0
70	308	168	0	0
70	306	180	0	0
70	310	180	0	0
38	308	185	0	0
70	306	184	0	0
4	304	187	0	0
70	304	184	0	0
70	311	192	0	0
70	309	194	0	0
70	304	192	0	0
70	307	192	0	0
4	318	174	0	0
70	316	189	0	0
4	317	172	0	0
70	316	170	0	0
70	319	193	0	0
70	314	188	0	0
70	315	171	0	0
70	315	194	0	0
70	318	162	0	0
70	315	170	0	0
205	318	164	0	0
38	317	163	0	0
4	319	173	0	0
38	313	199	0	0
38	318	181	0	0
70	319	177	0	0
38	312	196	0	0
70	319	189	0	0
205	318	200	0	0
70	319	202	0	0
4	318	203	0	0
205	315	203	0	0
4	319	209	0	0
205	308	208	0	0
205	304	215	0	0
70	307	208	0	0
70	308	215	0	0
70	309	210	0	0
38	318	215	0	0
70	309	209	0	0
70	306	213	0	0
70	312	208	0	0
4	311	218	0	0
4	311	219	0	0
4	312	217	0	0
70	307	222	0	0
70	314	218	0	0
70	312	227	0	0
70	316	227	0	0
4	313	228	0	0
70	308	227	0	0
4	304	227	0	0
38	305	236	0	0
4	312	239	0	0
70	314	233	0	0
70	311	237	0	0
70	317	240	0	0
70	306	246	0	0
70	315	245	0	0
70	305	242	0	0
70	306	240	0	0
4	310	243	0	0
70	317	246	0	0
70	306	252	0	0
70	318	248	0	0
4	309	248	0	0
70	307	253	0	0
4	317	253	0	0
4	318	256	0	0
70	313	258	0	0
4	314	261	0	0
205	313	262	0	0
70	304	263	0	0
70	311	258	0	0
4	313	268	0	0
4	318	271	0	0
70	317	269	0	0
70	312	271	0	0
70	312	267	0	0
70	317	271	0	0
4	317	273	0	0
70	313	274	0	0
70	306	273	0	0
70	309	281	0	0
70	308	283	0	0
70	305	280	0	0
4	317	282	0	0
70	319	285	0	0
4	314	283	0	0
70	312	286	0	0
70	309	295	0	0
4	310	290	0	0
4	312	303	0	0
4	313	300	0	0
38	314	301	0	0
70	308	299	0	0
206	319	300	0	0
207	308	301	0	0
117	306	300	0	0
117	304	299	6	0
117	304	303	2	0
70	315	301	0	0
24	304	300	2	1
207	308	304	0	0
207	305	304	3	0
70	307	311	0	0
70	305	308	0	0
4	312	307	0	0
70	310	310	0	0
54	314	307	0	0
38	310	315	0	0
70	310	313	0	0
38	307	326	0	0
4	310	326	0	0
70	311	324	0	0
70	310	323	0	0
4	308	330	0	0
70	310	330	0	0
70	308	336	0	0
205	307	340	0	0
4	306	340	0	0
70	311	342	0	0
70	304	348	0	0
70	304	347	0	0
70	306	345	0	0
70	310	347	0	0
70	311	347	0	0
38	307	357	0	0
70	306	357	0	0
38	306	358	0	0
70	305	359	0	0
205	310	359	0	0
70	304	352	0	0
70	310	363	0	0
70	306	372	0	0
70	310	373	0	0
70	307	371	0	0
70	305	377	0	0
70	311	390	0	0
70	307	385	0	0
70	308	385	0	0
4	305	392	0	0
4	306	394	0	0
70	306	395	0	0
38	307	392	0	0
70	309	397	0	0
4	307	394	0	0
70	309	406	0	0
70	305	406	0	0
70	311	407	0	0
70	313	386	0	0
70	318	368	0	0
205	315	394	0	0
70	319	381	0	0
70	317	385	0	0
70	319	394	0	0
70	312	399	0	0
70	317	374	0	0
4	318	407	0	0
70	316	377	0	0
70	313	407	0	0
70	316	375	0	0
4	317	395	0	0
70	313	379	0	0
70	318	384	0	0
70	313	375	0	0
205	315	403	0	0
205	318	372	0	0
70	312	402	0	0
70	312	368	0	0
70	319	401	0	0
70	316	395	0	0
4	306	413	0	0
70	308	412	0	0
70	304	408	0	0
70	319	414	0	0
70	308	414	0	0
38	304	411	0	0
70	323	377	0	0
205	327	382	0	0
70	322	377	0	0
70	324	381	0	0
70	320	383	0	0
70	327	395	0	0
4	320	382	0	0
70	327	404	0	0
70	326	393	0	0
70	326	398	0	0
70	320	409	0	0
70	320	380	0	0
70	325	411	0	0
4	320	407	0	0
205	321	404	0	0
70	327	388	0	0
70	322	388	0	0
38	325	387	0	0
205	326	402	0	0
70	321	389	0	0
70	326	409	0	0
70	325	404	0	0
70	322	371	0	0
70	319	366	0	0
70	312	365	0	0
205	324	366	0	0
4	323	365	0	0
70	315	356	0	0
4	315	357	0	0
70	317	353	0	0
70	322	356	0	0
70	322	355	0	0
70	320	359	0	0
38	314	358	0	0
38	327	355	0	0
70	319	347	0	0
70	316	345	0	0
70	317	347	0	0
70	316	348	0	0
70	314	351	0	0
38	314	344	0	0
4	326	348	0	0
70	326	349	0	0
38	326	345	0	0
70	320	344	0	0
70	327	351	0	0
205	325	351	0	0
4	325	347	0	0
205	326	344	0	0
70	325	343	0	0
70	319	341	0	0
70	314	340	0	0
4	317	339	0	0
70	319	338	0	0
70	318	334	0	0
70	319	333	0	0
70	315	333	0	0
70	313	328	0	0
205	314	324	0	0
70	317	315	0	0
70	313	315	0	0
4	317	319	0	0
70	319	159	0	0
38	313	159	0	0
70	315	157	0	0
70	318	157	0	0
38	313	154	0	0
4	317	149	0	0
38	313	145	0	0
70	312	147	0	0
70	319	147	0	0
20	315	142	0	0
205	313	139	0	0
70	317	136	0	0
205	322	138	0	0
70	327	155	0	0
4	320	144	0	0
205	321	129	0	0
70	327	167	0	0
70	325	166	0	0
70	323	151	0	0
38	324	159	0	0
70	321	156	0	0
70	327	150	0	0
4	322	164	0	0
20	321	142	0	0
70	320	147	0	0
70	320	160	0	0
70	327	169	0	0
4	324	168	0	0
70	324	172	0	0
70	321	168	0	0
4	321	178	0	0
4	325	181	0	0
70	322	183	0	0
205	325	178	0	0
70	324	186	0	0
70	323	188	0	0
70	327	196	0	0
143	327	199	2	0
70	327	195	0	0
4	323	198	0	0
205	326	196	0	0
70	325	194	0	0
4	320	194	0	0
70	321	195	0	0
70	322	194	0	0
143	327	204	2	0
70	326	206	0	0
70	321	203	0	0
70	320	205	0	0
1	326	201	1	1
38	327	213	0	0
70	327	216	0	0
70	322	228	0	0
4	322	226	0	0
70	323	229	0	0
70	323	233	0	0
70	322	232	0	0
4	321	238	0	0
70	326	232	0	0
4	320	242	0	0
38	325	247	0	0
205	321	242	0	0
70	320	251	0	0
4	324	250	0	0
70	326	248	0	0
70	327	260	0	0
70	323	261	0	0
70	325	259	0	0
70	324	263	0	0
70	320	267	0	0
70	327	265	0	0
70	323	268	0	0
70	326	271	0	0
70	321	278	0	0
70	320	273	0	0
70	325	275	0	0
4	326	274	0	0
70	324	279	0	0
70	324	286	0	0
4	322	286	0	0
70	323	285	0	0
206	327	291	0	0
70	321	287	0	0
206	323	303	0	0
204	321	302	0	0
55	324	300	0	0
55	324	301	0	0
70	323	295	0	0
70	322	288	0	0
70	324	318	0	0
70	323	315	0	0
70	327	315	0	0
70	321	314	0	0
70	321	313	0	0
70	326	326	0	0
70	322	320	0	0
70	326	322	0	0
70	320	326	0	0
70	324	334	0	0
70	322	332	0	0
205	322	330	0	0
70	323	331	0	0
70	330	406	0	0
70	331	377	0	0
70	328	409	0	0
70	328	377	0	0
38	334	387	0	0
4	334	381	0	0
70	334	414	0	0
70	330	398	0	0
4	331	405	0	0
70	329	413	0	0
70	332	373	0	0
38	332	371	0	0
70	328	362	0	0
70	333	365	0	0
205	329	366	0	0
70	334	365	0	0
70	330	353	0	0
70	334	353	0	0
70	333	355	0	0
205	335	355	0	0
38	329	356	0	0
70	334	344	0	0
38	330	349	0	0
4	330	344	0	0
70	329	346	0	0
4	333	336	0	0
70	331	339	0	0
70	328	334	0	0
70	335	330	0	0
205	334	333	0	0
70	329	329	0	0
70	335	326	0	0
4	328	325	0	0
70	328	322	0	0
70	329	327	0	0
70	335	313	0	0
70	328	317	0	0
70	328	314	0	0
4	330	313	0	0
70	333	312	0	0
70	331	296	0	0
205	335	303	0	0
70	335	301	0	0
204	329	291	0	0
205	334	289	0	0
70	334	293	0	0
70	329	283	0	0
70	330	284	0	0
70	331	280	0	0
70	328	286	0	0
70	331	287	0	0
70	328	264	0	0
70	331	257	0	0
70	332	258	0	0
70	332	251	0	0
70	328	249	0	0
70	335	254	0	0
70	331	243	0	0
70	330	242	0	0
70	335	244	0	0
70	329	235	0	0
70	328	232	0	0
70	329	233	0	0
70	329	231	0	0
4	333	225	0	0
70	330	231	0	0
70	328	225	0	0
70	333	221	0	0
4	331	214	0	0
4	328	213	0	0
70	333	208	0	0
70	333	213	0	0
144	333	201	2	0
143	330	205	4	0
143	332	204	6	0
70	335	205	0	0
1	329	206	0	1
143	332	199	6	0
143	330	198	0	0
1	329	198	0	1
70	330	181	0	0
205	328	176	0	0
70	328	177	0	0
70	330	182	0	0
70	331	181	0	0
70	334	181	0	0
70	332	177	0	0
70	331	175	0	0
205	333	174	0	0
70	330	169	0	0
205	333	167	0	0
4	335	163	0	0
70	332	167	0	0
4	333	161	0	0
4	331	156	0	0
347	331	141	2	0
20	328	142	0	0
20	335	142	2	0
205	320	125	0	0
70	329	121	0	0
70	325	122	0	0
70	314	120	0	0
205	319	126	0	0
70	320	421	0	0
70	325	421	0	0
70	326	420	0	0
70	327	421	0	0
205	314	421	0	0
70	317	421	0	0
4	317	417	0	0
70	330	421	0	0
70	334	418	0	0
70	328	418	0	0
4	311	422	0	0
70	305	421	0	0
70	308	416	0	0
205	309	429	0	0
4	310	426	0	0
70	310	427	0	0
38	306	428	0	0
70	304	425	0	0
70	319	425	0	0
70	306	425	0	0
70	313	424	0	0
70	322	431	0	0
38	324	425	0	0
70	314	430	0	0
70	313	427	0	0
4	308	431	0	0
70	335	431	0	0
38	314	428	0	0
205	334	424	0	0
70	318	429	0	0
4	332	430	0	0
144	334	435	2	0
3	333	434	0	0
70	297	430	0	0
70	296	426	0	0
70	297	431	0	0
70	302	428	0	0
70	290	427	0	0
38	290	428	0	0
7	280	434	4	0
7	280	438	4	0
70	283	429	0	0
7	280	436	0	0
7	278	436	0	0
9	277	439	0	0
7	278	434	4	0
7	279	436	0	0
7	277	434	4	0
9	277	435	0	0
4	273	429	0	0
7	277	436	0	0
7	277	438	4	0
7	278	438	4	0
7	279	438	4	0
7	279	434	4	0
70	274	427	0	0
22	273	435	0	1
39	275	439	1	1
70	271	425	0	0
70	270	428	0	0
70	269	427	0	0
70	270	425	0	0
70	264	424	0	0
70	268	426	0	0
5	271	433	0	0
70	266	430	0	0
7	279	440	0	0
7	277	440	0	0
7	278	440	0	0
43	273	444	6	0
5	273	440	0	0
20	272	446	6	0
20	272	443	4	0
7	280	440	0	0
40	278	443	0	1
38	271	441	1	1
70	262	170	0	0
70	259	169	0	0
70	257	168	0	0
38	259	174	0	0
4	248	171	0	0
4	248	173	0	0
38	248	172	0	0
58	241	173	2	0
205	240	171	0	0
4	241	169	0	0
70	247	175	0	0
70	247	170	0	0
70	244	173	0	0
70	255	167	0	0
38	253	161	0	0
4	240	161	0	0
70	252	161	0	0
70	244	160	0	0
70	249	161	0	0
38	260	167	0	0
70	256	167	0	0
70	259	167	0	0
4	263	162	0	0
70	245	153	0	0
205	243	159	0	0
4	246	155	0	0
4	246	150	0	0
70	243	145	0	0
70	246	148	0	0
70	251	144	0	0
38	254	157	0	0
4	248	159	0	0
4	263	153	0	0
210	258	156	2	0
210	256	158	2	0
5	245	3012	2	0
217	263	3011	4	0
215	254	3008	5	0
215	260	3016	6	0
210	268	3000	2	0
215	270	2997	6	0
216	267	3001	5	0
214	267	3000	6	0
216	264	3012	2	0
214	273	3002	6	0
5	282	3017	0	0
487	282	3020	0	0
4	262	147	0	0
20	261	142	0	0
205	262	151	0	0
70	259	149	0	0
70	257	149	0	0
1	257	140	1	1
15	252	140	0	0
1	254	138	0	1
70	260	134	0	0
70	260	133	0	0
38	259	131	0	0
20	247	142	0	0
70	245	135	0	0
70	245	130	0	0
205	246	130	0	0
70	245	138	0	0
5	549	3339	0	0
51	571	3316	6	0
51	565	3315	0	0
51	571	3319	6	0
51	566	3331	4	0
730	581	3329	0	0
730	581	3332	0	0
488	567	3331	0	0
51	568	3331	4	0
730	581	3335	0	0
730	580	3348	0	0
730	582	3337	6	0
730	578	3349	0	0
729	580	3344	5	0
51	563	3319	2	0
51	563	3322	2	0
51	563	3325	2	0
51	569	3315	0	0
51	563	3316	2	0
730	583	3327	0	0
51	571	3322	6	0
730	581	3351	5	0
51	571	3325	6	0
940	581	3342	0	0
730	580	3339	0	0
22	580	3343	0	0
105	546	3328	1	1
367	544	3307	4	0
51	548	3304	6	0
107	545	3307	0	1
51	548	3296	6	0
51	548	3300	6	0
108	546	3302	0	1
70	321	117	0	0
38	319	115	0	0
205	316	111	0	0
70	307	113	0	0
219	304	119	4	0
205	307	115	0	0
709	305	118	4	0
529	300	111	0	0
70	299	105	0	0
70	297	99	0	0
529	299	111	0	0
665	300	110	4	0
665	300	113	0	0
673	296	105	4	0
708	296	111	0	0
665	296	110	4	0
708	297	112	4	0
708	296	112	4	0
708	298	111	0	0
665	298	110	4	0
708	297	111	0	0
665	296	113	0	0
529	300	112	4	0
665	298	113	0	0
708	298	112	4	0
529	299	112	4	0
652	294	113	0	0
705	294	119	4	0
665	292	109	2	0
598	291	110	0	0
707	293	105	0	0
70	289	106	1	0
70	290	107	2	0
706	292	110	0	0
707	295	104	0	0
707	294	104	0	0
598	293	110	0	0
205	258	126	0	0
70	260	121	0	0
70	263	127	0	0
70	262	123	0	0
38	245	120	0	0
70	245	123	0	0
1012	239	126	6	0
1010	239	127	2	0
1010	239	124	2	0
1010	239	125	2	0
1011	239	123	0	0
1013	236	121	0	0
1013	235	120	0	0
1009	236	122	4	0
1011	237	121	0	0
1016	236	120	0	0
1013	238	123	0	0
1016	238	122	0	0
1016	232	127	0	0
1013	237	122	0	0
205	237	112	0	0
70	235	112	0	0
407	234	112	5	0
1011	235	119	0	0
1018	233	119	0	0
205	246	114	0	0
1012	232	119	0	0
70	244	112	0	0
1010	234	119	0	0
70	276	114	0	0
70	279	117	0	0
70	275	114	0	0
70	275	118	0	0
70	286	118	0	0
334	263	106	0	0
273	265	107	0	0
20	262	111	6	0
336	267	110	0	0
38	276	105	0	0
70	280	108	0	0
20	270	111	0	0
70	286	110	0	0
273	267	107	0	0
335	265	110	0	0
273	267	105	0	0
273	265	105	0	0
99	263	104	1	1
99	270	104	1	1
205	277	98	0	0
70	282	101	0	0
70	282	96	0	0
38	291	103	0	0
70	291	100	0	0
6	289	101	0	0
70	289	103	0	0
273	267	102	0	0
273	265	102	0	0
20	270	99	2	0
99	266	100	0	1
38	256	102	0	0
334	263	103	0	0
70	259	103	0	0
20	262	99	4	0
70	256	101	0	0
70	252	102	0	0
205	250	99	0	0
70	242	103	0	0
70	243	97	0	0
70	246	106	0	0
38	234	108	0	0
70	235	101	0	0
205	235	99	0	0
205	234	104	0	0
38	232	111	0	0
407	232	109	6	0
70	224	102	0	0
70	224	100	0	0
1009	226	111	2	0
1009	226	108	2	0
407	230	110	7	0
51	225	111	6	0
1012	225	107	0	0
407	224	113	5	0
1012	230	119	0	0
407	230	112	0	0
1012	226	119	0	0
1016	224	127	2	0
407	229	118	0	0
1027	228	119	0	0
1012	224	119	0	0
51	225	108	6	0
407	225	112	1	0
407	224	117	0	0
1009	225	119	0	0
1009	231	119	0	0
24	227	109	0	1
24	227	107	0	1
1	227	106	0	1
1	226	110	1	1
1188	223	110	6	0
1012	220	107	0	0
51	220	108	2	0
70	219	99	0	0
51	220	111	1	0
407	218	106	0	0
1013	222	120	2	0
1012	217	126	2	0
1016	218	123	2	0
1010	217	125	6	0
1010	217	127	2	0
1016	221	120	2	0
1013	221	121	2	0
1013	220	122	2	0
1011	217	124	2	0
1009	221	122	4	0
1013	219	123	2	0
1013	218	124	2	0
407	220	114	6	0
1010	223	119	0	0
1011	219	122	2	0
407	218	110	0	0
1011	222	119	2	0
1011	220	121	2	0
70	211	99	0	0
205	210	100	0	0
70	212	102	0	0
205	208	109	0	0
38	212	113	0	0
205	208	104	0	0
205	211	107	0	0
70	208	102	0	0
70	209	118	0	0
70	214	115	0	0
70	208	112	0	0
50	205	104	0	0
70	201	99	0	0
106	201	102	0	0
205	206	109	0	0
107	201	101	0	0
205	202	104	0	0
111	200	104	0	0
70	206	113	0	0
70	205	116	0	0
38	199	102	0	0
70	196	102	0	0
103	195	104	0	0
110	196	106	0	0
70	198	113	0	0
106	198	103	0	0
111	198	101	0	0
107	198	105	0	0
70	197	112	0	0
70	192	117	0	0
205	197	119	0	0
70	199	119	0	0
108	199	106	0	0
70	195	101	0	0
70	196	101	0	0
111	197	104	0	0
70	189	98	0	0
38	185	103	0	0
38	188	101	0	0
70	189	109	0	0
70	188	98	0	0
70	184	113	0	0
70	184	118	0	0
70	188	113	0	0
70	185	105	0	0
70	183	98	0	0
70	181	105	0	0
205	182	116	0	0
70	177	107	0	0
38	183	102	0	0
38	183	110	0	0
24	176	107	0	1
70	172	117	0	0
70	171	104	0	0
70	173	116	0	0
70	175	107	0	0
70	171	100	0	0
70	172	112	0	0
70	168	119	0	0
205	168	117	0	0
38	167	110	0	0
70	162	111	0	0
70	163	114	0	0
334	162	107	0	0
70	162	99	0	0
70	163	101	0	0
38	166	102	0	0
205	161	115	0	0
70	162	112	0	0
70	163	100	0	0
100	160	103	0	1
100	160	108	0	1
70	153	99	0	0
70	152	97	0	0
334	158	107	0	0
334	159	107	0	0
70	155	99	0	0
70	155	105	0	0
70	159	119	0	0
70	155	118	0	0
70	146	99	0	0
70	149	96	0	0
70	149	107	0	0
70	149	109	0	0
70	149	97	0	0
205	151	118	0	0
70	149	116	0	0
70	150	116	0	0
70	138	107	0	0
205	139	99	0	0
70	138	103	0	0
70	140	107	0	0
70	141	108	0	0
215	139	113	7	0
70	139	100	0	0
215	136	115	7	0
70	138	115	0	0
205	142	113	0	0
70	135	96	0	0
70	133	107	0	0
216	135	111	7	0
70	135	109	0	0
70	130	102	0	0
70	132	110	0	0
70	128	109	0	0
214	130	113	7	0
70	129	115	0	0
38	132	116	0	0
217	132	114	6	0
205	131	118	0	0
215	128	113	0	0
70	134	119	0	0
70	120	104	0	0
216	124	113	7	0
70	121	98	0	0
70	126	105	0	0
70	122	102	0	0
215	122	114	7	0
70	126	100	0	0
215	127	116	0	0
205	125	101	0	0
70	122	100	0	0
214	120	113	7	0
70	113	98	0	0
70	115	107	0	0
70	115	103	0	0
70	119	102	0	0
70	119	103	0	0
38	118	116	0	0
70	118	114	0	0
70	114	116	0	0
70	116	116	0	0
70	112	117	0	0
70	111	109	0	0
70	106	100	0	0
70	108	111	0	0
70	107	99	0	0
70	104	112	0	0
70	111	106	0	0
70	110	104	0	0
70	106	117	0	0
70	107	110	0	0
70	107	101	0	0
70	103	103	0	0
70	101	103	0	0
70	91	97	0	0
70	88	100	0	0
70	90	106	0	0
70	90	109	0	0
70	91	103	0	0
70	90	115	0	0
205	88	118	0	0
205	88	98	0	0
205	80	103	0	0
205	86	106	0	0
205	84	117	0	0
70	82	102	0	0
205	85	101	0	0
70	82	115	0	0
70	76	102	0	0
274	76	104	4	0
278	79	104	4	0
70	78	102	0	0
70	75	108	0	0
205	73	115	0	0
38	76	119	0	0
205	72	104	0	0
70	74	119	0	0
70	72	108	0	0
205	78	115	0	0
2	77	107	0	1
205	68	119	0	0
70	70	118	0	0
70	61	118	0	0
70	58	112	0	0
205	60	122	0	0
205	61	122	0	0
205	58	120	0	0
70	57	120	0	0
205	65	124	0	0
70	58	127	0	0
70	67	125	0	0
205	67	124	0	0
205	66	121	0	0
205	71	122	0	0
205	53	121	0	0
70	52	124	0	0
70	131	122	0	0
70	128	121	0	0
70	138	121	0	0
70	138	127	0	0
70	140	125	0	0
205	142	120	0	0
70	151	127	0	0
70	148	125	0	0
70	157	123	0	0
70	153	125	0	0
70	164	125	0	0
70	174	123	0	0
70	174	122	0	0
70	182	121	0	0
70	183	126	0	0
70	180	121	0	0
1	180	127	0	1
24	181	122	0	1
70	189	122	0	0
70	195	120	0	0
70	199	122	0	0
80	203	123	4	0
70	211	122	0	0
80	208	122	4	0
70	212	120	0	0
51	426	3375	0	0
205	425	3381	4	0
51	430	3380	6	0
3	442	3372	0	0
730	446	3373	2	0
51	424	3375	0	0
29	441	3375	0	0
278	440	3369	0	0
730	444	3375	6	0
730	452	3378	7	0
29	454	3375	0	0
1012	444	3366	4	0
3	450	3372	0	0
3	452	3370	0	0
3	450	3369	0	0
1187	446	3367	4	0
29	438	3375	0	0
51	426	3384	4	0
3	440	3371	0	0
3	441	3369	0	0
730	442	3367	5	0
278	441	3372	0	0
1012	448	3366	4	0
278	449	3371	0	0
730	454	3373	3	0
278	451	3369	0	0
730	450	3376	0	0
278	452	3371	0	0
5	457	3352	2	0
29	451	3375	0	0
143	270	2947	0	0
143	271	2956	4	0
143	264	2953	2	0
143	266	2956	4	0
143	268	2947	0	0
51	266	2963	2	0
51	278	2967	0	0
143	273	2954	6	0
41	268	2960	4	0
145	278	2962	0	0
58	274	2952	0	0
51	272	2974	6	0
51	267	2971	2	0
20	279	2956	0	0
51	278	2974	4	0
143	273	2950	6	0
51	275	2974	4	0
51	276	2960	2	0
51	266	2965	2	0
51	271	2965	6	0
51	271	2963	6	0
57	272	2972	4	0
143	264	2950	2	0
20	281	2959	0	0
145	276	2959	6	0
20	278	2959	0	0
145	279	2954	4	0
51	276	2957	2	0
145	286	2964	2	0
145	286	2968	2	0
51	281	2954	0	0
51	284	2971	4	0
20	284	2966	0	0
51	286	2966	6	0
145	283	2971	0	0
51	283	2957	6	0
145	283	2958	2	0
51	272	2971	6	0
20	283	2962	0	0
51	284	2960	6	0
51	282	2966	2	0
57	281	2969	4	0
51	280	2967	0	0
666	300	2945	4	0
666	298	2945	0	0
51	301	2947	3	0
51	303	2944	6	0
51	296	2947	4	0
666	296	2945	6	0
51	303	2942	6	0
666	298	2942	2	0
666	300	2942	0	0
51	288	2937	2	0
666	296	2942	1	0
666	292	2941	0	0
51	293	2940	0	0
51	300	2939	0	0
1166	471	3383	7	0
730	469	3385	2	0
730	470	3381	6	0
729	468	3383	7	0
5	289	2933	0	0
51	289	2931	0	0
51	288	2934	2	0
20	115	590	0	0
20	115	587	0	0
20	115	584	0	0
34	113	589	1	0
1	134	590	0	0
38	112	589	1	0
34	112	588	1	0
20	135	591	0	0
20	115	593	0	0
20	114	599	6	0
55	114	597	0	0
20	115	596	0	0
54	116	598	0	0
55	114	598	1	0
34	113	592	1	0
4	115	604	0	0
20	114	603	6	0
0	112	603	0	0
20	112	605	0	0
1	134	596	0	0
55	117	612	0	0
20	112	610	0	0
1	131	612	0	0
59	114	608	4	0
1	135	603	0	0
3	121	609	0	0
55	116	612	3	0
1	131	608	0	0
1	130	614	0	0
20	140	593	0	0
20	143	593	0	0
4	142	604	0	0
1	136	600	0	0
1	136	608	0	0
1	117	622	0	0
0	128	620	0	0
0	128	616	0	0
1	113	623	0	0
1	138	621	0	0
0	122	621	0	0
1	118	623	0	0
1	133	616	0	0
0	120	621	0	0
59	112	616	4	0
0	116	621	0	0
1	127	623	0	0
1	120	623	0	0
0	122	622	0	0
0	131	620	0	0
0	128	621	0	0
20	124	620	0	0
1	115	623	0	0
0	126	617	0	0
0	130	622	0	0
1	132	617	0	0
0	124	621	0	0
1	114	631	0	0
1	112	630	0	0
0	127	629	0	0
45	125	624	0	0
1	114	629	0	0
45	125	626	4	0
1	112	624	4	0
118	130	626	6	0
1	122	628	4	0
0	128	630	0	0
45	126	624	0	0
45	124	626	4	0
38	128	628	0	0
45	126	626	4	0
192	125	631	0	0
45	127	626	4	0
45	127	624	0	0
45	124	624	0	0
59	135	630	4	0
1	121	630	0	0
192	125	629	0	0
1	134	628	1	1
1	131	631	0	1
1	148	584	0	0
59	148	596	2	0
20	150	593	0	0
20	146	593	0	0
0	157	588	0	0
6	118	1599	0	0
6	118	1605	0	0
3	133	1601	0	0
3	131	1606	0	0
14	129	1598	0	0
3	129	1601	0	0
14	129	1604	0	0
66	126	686	1	1
36	120	695	6	0
34	150	646	0	0
1	149	644	0	0
1	151	642	0	0
4	150	641	0	0
34	147	643	1	0
34	145	646	1	0
1	146	644	0	0
34	149	652	0	0
1	149	656	0	0
34	147	648	0	0
1	144	642	0	0
1	145	644	0	0
34	150	654	0	0
4	158	641	0	0
34	158	649	0	0
34	152	652	0	0
34	153	649	0	0
34	155	654	0	0
0	158	662	0	0
34	156	657	0	0
1	153	658	0	0
34	159	645	0	0
1	156	660	0	0
145	125	3511	6	0
17	141	3493	0	0
3	126	3511	6	0
17	141	3492	0	0
145	127	3511	2	0
5	136	3492	0	0
51	139	3493	4	0
51	139	3488	0	0
70	158	687	0	0
36	158	684	2	0
36	157	685	6	0
1	155	685	6	0
36	158	683	6	0
5	138	1612	0	0
6	139	1592	0	0
5	138	1593	0	0
6	139	1610	0	0
6	138	2556	0	0
121	139	2554	2	0
6	138	2537	0	0
7	139	2553	4	0
0	149	629	0	0
1	149	635	0	0
0	151	628	0	0
1	151	636	0	0
0	150	627	0	0
1	144	636	0	0
1	146	635	0	0
59	152	615	2	0
59	158	612	6	0
3	157	617	2	0
4	155	638	0	0
15	157	619	0	0
59	158	614	2	0
1	159	617	0	1
1	167	622	0	0
8	160	615	2	0
36	167	616	0	0
36	166	618	0	0
8	161	615	0	0
34	163	618	0	0
0	163	619	0	0
1	163	622	0	0
0	166	637	0	0
1	165	621	0	0
36	166	617	0	0
1	164	642	0	0
34	167	639	1	0
34	163	647	0	0
34	165	617	0	0
20	161	623	0	0
0	162	632	0	0
4	165	647	0	0
20	161	621	2	0
74	165	604	6	0
54	158	602	0	0
55	163	601	0	0
55	164	602	0	0
1	166	604	0	1
0	171	616	0	0
1	168	615	0	0
1	171	618	0	0
1	175	646	0	0
34	171	620	0	0
1	172	647	5	0
34	169	642	1	0
34	172	621	0	0
0	169	622	0	0
0	174	644	3	0
36	169	621	0	0
72	175	602	3	0
72	173	601	3	0
72	175	600	3	0
72	175	601	3	0
72	173	603	3	0
72	173	602	3	0
72	174	603	3	0
72	173	604	3	0
55	169	600	0	0
72	172	601	3	0
55	169	601	0	0
72	171	603	3	0
72	172	603	3	0
72	172	602	3	0
55	168	602	0	0
72	174	601	3	0
72	175	604	3	0
72	174	604	3	0
72	174	602	3	0
59	172	607	6	0
72	173	605	3	0
72	171	600	3	0
72	172	604	3	0
72	171	601	3	0
72	171	602	3	0
72	172	600	3	0
72	173	600	3	0
72	174	600	3	0
72	172	605	3	0
59	154	593	6	0
1	166	596	0	0
72	175	597	3	0
3	165	599	0	0
5	165	598	0	0
72	175	598	3	0
0	160	597	0	0
72	172	597	3	0
72	175	599	3	0
72	171	599	3	0
72	172	598	3	0
72	171	598	3	0
1	159	596	0	0
20	165	593	0	0
53	166	599	1	0
1	158	598	0	0
1	163	596	0	0
20	173	593	0	0
72	172	599	3	0
0	162	598	0	0
72	173	598	3	0
72	173	597	3	0
20	161	593	0	0
72	174	598	3	0
72	174	597	3	0
72	174	599	3	0
72	173	599	3	0
72	171	597	3	0
20	169	593	0	0
20	177	585	0	0
20	177	577	0	0
20	177	581	0	0
20	177	589	0	0
70	185	569	6	0
70	188	570	6	0
70	184	575	6	0
70	190	574	6	0
70	188	574	6	0
72	176	599	3	0
72	176	597	3	0
72	177	598	3	0
72	176	598	3	0
72	177	597	3	0
72	177	599	3	0
59	177	595	2	0
20	177	593	0	0
72	177	602	3	0
72	176	602	3	0
72	177	600	3	0
72	176	604	3	0
72	178	602	3	0
72	177	601	3	0
72	176	601	3	0
72	176	600	3	0
72	177	603	3	0
72	176	603	3	0
72	178	601	3	0
59	184	604	6	0
59	184	606	2	0
191	181	614	0	0
191	177	612	0	0
0	177	609	0	0
191	181	610	0	0
191	181	612	0	0
191	177	614	0	0
191	181	616	0	0
191	177	616	0	0
1	177	629	0	0
1	181	626	0	0
0	182	628	0	0
191	181	620	0	0
191	177	618	0	0
191	177	620	0	0
191	181	618	0	0
191	189	614	0	0
191	189	612	0	0
191	185	612	0	0
191	185	620	0	0
191	185	618	0	0
191	189	610	0	0
191	189	616	0	0
191	189	618	0	0
191	185	610	0	0
191	189	608	0	0
191	185	608	0	0
191	185	616	0	0
191	185	614	0	0
191	189	620	0	0
28	165	1543	0	0
5	166	1546	0	0
6	165	1542	0	0
6	166	2490	0	0
55	166	2489	0	0
52	166	2487	0	0
1	190	628	0	0
0	188	628	0	0
1	188	625	0	0
1	190	625	0	0
1	190	631	0	0
1	185	625	0	0
1	186	629	0	0
1	187	638	0	0
1	184	636	0	0
1	185	633	0	0
1	189	634	0	0
1	183	632	0	0
1	182	635	0	0
1	179	633	0	0
1	186	633	0	0
1	182	639	0	0
1	179	646	0	0
1	180	641	0	0
1	186	646	0	0
1	186	644	0	0
1	177	646	0	0
1	185	641	0	0
1	178	644	0	0
1	186	640	0	0
1	182	643	0	0
1	182	647	0	0
1	181	642	0	0
37	184	645	0	0
1	182	640	0	0
0	176	644	3	0
0	177	643	3	0
1	191	647	0	0
1	190	644	0	0
72	199	615	3	0
72	199	614	3	0
72	198	612	7	0
72	198	619	3	0
72	198	613	3	0
72	198	620	3	0
72	197	612	3	0
72	197	610	3	0
72	198	614	7	0
72	198	618	3	0
72	197	614	3	0
1	194	645	0	0
72	197	617	3	0
72	198	615	3	0
72	199	616	3	0
72	197	611	3	0
72	197	616	7	0
72	197	618	3	0
0	199	644	0	0
1	199	632	0	0
1	197	631	0	0
34	193	639	1	0
1	192	632	0	0
72	199	619	3	0
72	198	611	3	0
72	198	610	3	0
72	197	613	3	0
72	199	611	7	0
72	197	615	3	0
72	199	610	3	0
72	199	617	3	0
72	199	620	7	0
72	198	617	7	0
72	199	612	3	0
72	199	613	3	0
72	198	616	3	0
72	199	618	7	0
34	194	630	0	0
1	192	628	0	0
34	195	643	1	0
0	197	645	0	0
34	195	647	2	0
34	196	647	1	0
72	197	620	3	0
0	195	627	0	0
72	197	619	7	0
0	197	628	0	0
1	195	637	0	1
45	199	640	1	1
72	207	614	7	0
72	207	615	3	0
72	206	615	3	0
72	206	614	3	0
72	206	613	3	0
72	207	616	3	0
72	206	620	3	0
72	203	615	3	0
72	203	613	3	0
72	203	614	3	0
72	204	611	3	0
72	203	612	3	0
72	200	610	3	0
72	207	617	3	0
72	205	614	3	0
72	200	613	3	0
72	205	615	3	0
72	206	619	7	0
72	200	615	3	0
72	205	619	3	0
34	201	631	0	0
1	207	630	0	0
72	204	617	7	0
72	200	619	7	0
72	203	620	3	0
72	203	619	3	0
72	205	617	7	0
72	205	618	7	0
72	204	616	3	0
72	200	617	3	0
72	205	616	3	0
1	206	637	0	0
34	204	633	0	0
0	207	641	0	0
1	204	635	0	0
72	204	614	7	0
59	203	608	2	0
72	205	612	3	0
72	200	614	3	0
72	204	612	3	0
72	205	620	3	0
72	206	616	7	0
72	207	619	3	0
72	203	611	3	0
72	207	618	3	0
72	200	611	3	0
72	204	613	7	0
72	207	620	3	0
72	206	617	7	0
72	205	613	7	0
72	200	612	3	0
72	204	615	7	0
72	203	617	3	0
72	204	618	3	0
72	204	620	3	0
72	200	618	3	0
0	201	646	0	0
34	206	645	1	0
1	205	642	0	0
72	204	619	3	0
34	204	647	0	0
72	203	616	3	0
0	206	643	0	0
1	203	644	0	0
0	201	645	3	0
1	202	646	0	0
72	206	618	3	0
0	206	639	0	0
1	201	633	0	0
72	200	616	3	0
0	204	631	0	0
1	206	634	0	0
72	200	620	3	0
0	201	627	0	0
1	201	635	0	0
72	203	618	3	0
34	174	650	1	0
1	175	650	0	0
1	173	655	0	0
1	174	652	0	0
1	172	649	0	0
1	179	655	0	0
34	187	648	0	0
1	188	648	0	0
34	178	654	0	0
1	168	653	6	0
1	179	649	0	0
0	181	650	0	0
0	187	655	0	0
1	169	649	5	0
1	169	654	0	0
1	171	648	5	0
1	179	648	0	0
1	185	650	0	0
1	184	648	0	0
1	189	654	0	0
37	177	651	0	0
0	192	655	0	0
0	195	652	0	0
1	184	651	0	0
0	193	648	0	0
0	198	652	3	0
34	197	649	6	0
1	204	652	0	0
1	177	655	0	0
37	183	655	0	0
1	171	655	0	0
0	182	653	0	0
34	204	650	1	0
34	177	653	1	0
0	188	650	3	0
34	190	655	1	0
1	182	655	0	0
1	171	651	0	0
1	176	653	0	0
0	195	655	0	0
34	186	650	0	0
0	194	651	0	0
34	180	652	2	0
37	191	649	0	0
34	190	650	2	0
34	196	648	6	0
34	195	649	6	0
0	194	654	0	0
0	201	655	0	0
0	200	650	3	0
0	200	652	0	0
1	200	648	6	0
4	166	654	7	0
1	161	654	0	0
4	166	652	5	0
38	166	660	0	0
4	166	657	2	0
4	165	656	3	0
38	165	661	0	0
0	164	666	0	0
1	172	664	0	0
0	175	663	0	0
1	174	665	0	0
237	172	662	0	0
0	173	671	0	0
1	170	658	0	0
1	175	666	0	0
1	170	661	0	0
1	169	662	0	0
1	169	667	0	0
1	168	664	0	0
0	168	671	0	0
1	171	670	0	0
1	174	669	0	0
1	173	660	0	0
1	168	658	0	0
34	174	657	1	0
2	183	660	0	0
1	181	658	0	0
4	183	658	0	0
34	181	656	1	0
0	180	659	0	0
1	181	666	0	0
1	177	671	0	0
0	182	664	0	0
1	179	663	0	0
34	179	661	2	0
97	178	668	1	0
1	182	669	0	0
34	177	665	1	0
1	181	671	0	0
1	179	665	0	0
1	179	656	0	0
1	183	666	0	0
0	183	668	0	0
97	178	669	0	0
97	179	669	2	0
97	179	668	0	0
0	178	659	0	0
0	189	656	0	0
34	191	656	1	0
34	190	658	0	0
37	190	663	0	0
34	186	659	1	0
1	185	656	0	0
1	186	658	0	0
34	187	669	1	0
38	187	665	0	0
1	188	667	0	0
1	186	670	0	0
0	189	660	0	0
38	191	670	0	0
0	188	658	0	0
38	186	662	0	0
38	187	661	0	0
1	184	665	0	0
1	190	665	0	0
37	191	666	0	0
0	185	667	0	0
34	191	668	1	0
70	167	677	0	0
36	161	679	6	0
0	178	675	0	0
1	183	677	0	0
1	188	674	0	0
1	180	673	0	0
1	178	677	0	0
0	181	677	0	0
0	183	675	0	0
36	166	686	4	0
36	165	687	4	0
36	161	680	6	0
34	160	687	6	0
34	164	685	4	0
36	160	682	3	0
34	160	685	6	0
34	161	686	1	0
34	163	686	4	0
34	161	681	4	0
34	160	686	0	0
70	173	686	0	0
0	299	537	6	0
37	297	539	4	0
37	297	540	4	0
34	296	540	4	0
34	297	541	4	0
20	311	562	0	0
20	298	545	0	0
145	308	571	0	0
118	306	546	2	0
5	303	563	0	0
37	296	538	4	0
118	310	546	2	0
1	309	543	0	1
1	306	565	2	1
1	306	544	1	1
20	295	529	0	0
0	295	535	6	0
23	292	531	4	0
0	290	535	6	0
1	291	530	4	0
37	295	541	4	0
1	289	530	4	0
34	295	539	4	0
23	289	547	0	0
23	300	530	4	0
37	294	540	4	0
1	289	537	4	0
37	300	533	4	0
4	292	548	6	0
34	299	535	4	0
188	290	548	0	0
188	289	548	0	0
5	321	535	4	0
29	309	531	2	0
5	308	529	0	0
20	298	529	0	0
20	295	548	0	0
0	290	543	6	0
37	295	538	4	0
4	291	545	6	0
37	300	535	4	0
1	302	531	4	0
23	291	539	6	0
37	299	534	4	0
34	301	534	4	0
29	320	531	4	0
34	298	534	4	0
34	299	532	4	0
29	320	532	4	0
89	316	534	4	0
37	298	533	4	0
1	317	533	1	1
1	313	531	1	1
139	287	566	0	0
20	280	554	0	0
37	281	552	2	0
37	281	553	2	0
37	281	551	2	0
3	281	550	2	0
48	311	523	4	0
20	318	521	0	0
5	289	521	4	0
20	313	521	0	0
3	290	524	4	0
15	288	525	4	0
3	306	521	4	0
3	306	522	4	0
7	289	524	6	0
5	308	522	4	0
11	311	521	4	0
1	292	522	1	1
1	309	525	0	1
1	292	526	1	1
37	275	553	2	0
20	276	554	6	0
37	275	552	2	0
3	275	549	2	0
37	275	551	2	0
280	277	547	0	0
1	278	551	0	1
1	269	549	1	1
1	269	557	1	1
1	264	539	0	1
1	269	553	1	1
223	274	566	6	0
1	274	563	0	1
63	287	571	0	0
5	257	538	0	0
1	262	543	0	1
1	261	541	1	1
1	259	549	1	1
1	259	553	1	1
1	259	557	1	1
34	258	532	0	0
34	257	529	0	0
1	253	531	0	0
1	252	548	0	0
45	253	536	0	0
1	250	556	0	0
45	251	539	4	0
34	251	529	0	0
45	252	536	0	0
37	252	550	0	0
1	251	542	0	0
34	251	532	0	0
44	251	537	6	0
37	248	566	0	0
37	249	555	0	0
34	251	553	0	0
37	251	546	0	0
34	250	549	0	0
1	248	559	0	0
1	248	564	0	0
1	250	561	0	0
37	254	529	0	0
45	253	539	4	0
45	252	539	4	0
1	250	565	0	0
1	250	567	0	0
34	250	546	0	0
1	248	552	0	0
1	249	530	0	0
34	251	558	0	0
45	250	539	6	0
1	249	544	0	0
34	248	555	0	0
45	251	536	0	0
1	248	541	0	0
1	252	552	0	0
45	250	536	6	0
1	248	532	0	0
1	252	541	0	1
37	247	529	0	0
37	246	546	0	0
1	247	549	0	0
6	257	1481	0	0
70	235	538	0	0
70	233	536	0	0
70	236	565	6	0
70	233	541	0	0
70	232	531	0	0
70	237	558	6	0
70	234	562	6	0
70	237	553	6	0
70	234	533	0	0
70	235	546	6	0
70	238	541	6	0
70	235	554	6	0
70	236	525	0	0
70	233	524	0	0
70	235	527	0	0
8	233	521	1	0
8	232	520	1	0
1	248	570	0	0
70	232	571	6	0
70	233	568	6	0
70	236	570	6	0
70	233	575	6	0
70	232	573	6	0
1	248	574	0	0
37	249	573	0	0
37	251	570	0	0
34	250	574	0	0
70	237	574	6	0
1	252	573	0	0
70	239	582	6	0
70	238	579	6	0
70	239	576	6	0
70	236	583	6	0
70	235	576	6	0
1	265	583	0	0
1	256	581	0	0
1	263	578	0	0
1	260	583	0	0
70	234	580	6	0
1	253	581	0	0
70	232	581	6	0
20	246	582	0	0
1	267	591	0	0
1	251	590	0	0
1	269	585	0	0
1	257	588	0	0
1	252	584	0	0
1	279	590	0	0
1	277	580	0	0
1	276	584	0	0
1	275	588	0	0
1	283	587	0	0
1	281	585	0	0
20	293	581	0	0
20	289	581	0	0
121	295	579	0	0
7	295	578	4	0
5	299	577	4	0
15	298	579	0	0
29	302	578	2	0
1	303	577	0	1
1	297	577	0	1
1	307	579	1	1
0	310	591	2	0
1	306	586	2	0
1	315	586	2	0
1	305	597	2	0
1	310	598	2	0
0	299	593	2	0
1	285	593	0	0
0	310	594	2	0
34	313	597	2	0
1	316	592	2	0
34	315	596	2	0
1	315	599	2	0
4	286	599	0	0
4	283	599	0	0
1	302	598	2	0
0	304	592	2	0
6	299	1521	4	0
6	318	1505	0	0
5	304	1509	0	0
6	303	1507	0	0
15	295	1523	4	0
5	319	1513	0	0
18	299	1524	4	0
6	317	1512	0	0
1	306	1506	3	1
1	319	1510	0	1
1	321	598	2	0
34	325	595	2	0
0	322	596	2	0
0	324	591	2	0
1	320	591	2	0
1	322	586	2	0
138	343	581	6	0
1	344	579	0	0
1	348	583	0	0
1	349	579	0	0
1	350	586	0	0
1	344	587	0	0
6	333	1513	0	0
15	331	1509	6	0
1	315	1503	1	1
27	311	1473	0	0
6	308	1473	0	0
15	307	1478	0	0
45	318	1492	4	0
25	322	1474	4	0
15	317	1479	0	0
3	312	1479	4	0
45	319	1489	2	0
6	321	1479	0	0
3	326	1483	7	0
3	327	1484	7	0
45	319	1492	2	0
42	316	1490	0	0
45	318	1489	0	0
15	317	1474	4	0
45	316	1489	0	0
47	312	1475	4	0
45	319	1490	2	0
3	322	1489	4	0
48	322	1475	4	0
45	317	1489	0	0
7	322	1488	4	0
45	319	1491	2	0
7	323	1491	4	0
6	328	1482	7	0
7	323	1493	0	0
3	323	1492	4	0
1	321	1477	0	1
2	320	1478	1	1
36	337	522	0	0
36	337	521	0	0
36	338	523	0	0
36	339	524	0	0
36	342	513	0	0
36	341	513	0	0
15	310	1467	0	0
6	308	1466	4	0
47	306	1466	0	0
36	350	521	0	0
36	351	534	0	0
0	350	525	0	0
36	344	530	0	0
36	346	533	0	0
36	350	523	0	0
36	345	531	0	0
36	344	529	0	0
36	348	534	0	0
1	351	524	0	0
36	349	514	0	0
36	351	517	0	0
36	344	514	0	0
36	346	517	0	0
307	345	513	0	0
0	349	519	0	0
1	351	518	0	0
36	348	519	0	0
0	347	517	0	0
307	340	511	0	0
406	346	507	3	0
1	342	508	0	0
36	353	518	0	0
36	357	525	0	0
36	357	523	0	0
0	356	524	0	0
1	355	523	0	0
36	354	527	0	0
1	352	526	0	0
1	359	537	0	0
36	356	522	0	0
36	358	534	0	0
0	358	545	0	0
36	357	535	0	0
1	356	550	0	0
36	356	526	0	0
283	353	521	0	0
307	354	525	0	0
36	353	520	0	0
1	353	525	0	0
36	352	527	0	0
0	358	535	0	0
111	359	547	0	0
0	353	542	0	0
0	352	556	0	0
0	356	561	0	0
0	353	567	0	0
1	353	565	0	0
0	355	571	0	0
2	359	572	3	1
6	319	2457	0	0
175	318	2454	6	0
14	316	2456	0	0
6	304	2453	0	0
15	290	1467	0	0
6	289	1465	2	0
0	305	501	0	0
0	307	497	0	0
0	308	498	0	0
0	303	500	0	0
0	313	488	0	0
0	308	493	0	0
0	306	495	0	0
0	312	491	0	0
0	313	492	0	0
0	311	495	0	0
0	317	489	0	0
0	316	491	0	0
177	327	490	0	0
1	325	493	0	1
1	342	502	0	0
34	340	493	4	0
34	341	494	4	0
34	347	496	4	0
34	347	495	4	0
37	344	489	4	0
1	351	489	4	0
1	349	488	4	0
73	351	492	0	1
0	319	486	0	0
0	317	484	0	0
0	314	486	0	0
0	318	481	0	0
0	317	486	0	0
137	341	487	4	0
0	321	481	0	0
0	321	484	0	0
0	312	481	0	0
0	312	485	0	0
0	344	485	4	0
15	323	487	4	0
37	347	486	4	0
177	327	487	0	0
0	314	482	0	0
1	348	483	4	0
1	325	487	0	1
0	317	478	0	0
0	319	477	0	0
0	314	478	0	0
0	316	476	0	0
0	311	477	0	0
0	313	476	0	0
0	323	479	0	0
0	314	472	0	0
0	317	472	0	0
0	312	472	0	0
0	315	474	0	0
0	313	467	0	0
0	314	469	0	0
0	310	467	0	0
0	311	465	0	0
38	318	458	4	0
0	312	460	0	0
0	313	463	0	0
0	306	457	0	0
70	322	457	4	0
0	308	457	0	0
0	312	458	0	0
0	310	460	0	0
0	307	460	0	0
0	310	462	0	0
34	341	457	0	0
55	318	449	4	0
70	325	451	4	0
70	331	452	4	0
70	327	454	4	0
55	319	449	4	0
1	321	450	3	1
408	328	446	0	0
70	326	443	4	0
70	315	445	4	0
55	321	442	4	0
38	315	447	4	0
1	320	445	0	1
1	327	447	1	1
37	346	459	0	0
1	351	440	0	0
1	344	437	0	0
34	348	449	0	0
1	351	457	0	0
1	346	441	0	0
1	347	445	0	0
1	350	434	0	0
1	345	447	0	0
34	351	453	0	0
1	351	444	0	0
1	350	455	0	0
1	348	450	0	0
1	357	484	4	0
1	352	472	0	0
0	354	481	4	0
1	358	505	0	0
1	352	483	4	0
41	358	492	0	0
20	357	494	4	0
34	355	491	4	0
34	355	488	4	0
1	355	502	0	0
71	356	495	1	1
70	358	495	1	1
72	356	492	1	1
68	357	462	0	0
0	357	459	0	0
34	358	460	0	0
1	354	459	0	0
1	357	457	0	0
34	353	463	0	0
1	356	456	0	0
34	359	453	0	0
1	359	449	0	0
1	358	453	0	0
0	352	449	0	0
37	358	448	0	0
34	355	454	0	0
1	355	451	0	0
34	356	454	0	0
0	353	449	0	0
0	358	440	0	0
34	359	444	0	0
1	353	445	0	0
0	356	446	0	0
1	355	439	0	0
0	286	491	4	0
0	285	499	4	0
0	286	496	4	0
0	283	497	4	0
1	280	491	0	1
142	276	448	6	0
21	277	495	2	0
21	275	497	2	0
6	279	494	0	0
21	272	486	0	0
278	273	494	2	0
278	275	494	2	0
22	276	495	2	0
22	276	496	2	0
21	274	488	0	0
278	274	495	2	0
21	277	491	0	0
3	274	494	2	0
0	279	498	4	0
197	278	493	1	1
21	267	486	0	0
994	271	491	7	0
994	269	488	7	0
0	258	502	0	0
0	263	503	0	0
38	260	499	0	0
70	258	477	6	0
0	258	500	0	0
0	256	508	0	0
0	258	507	0	0
0	257	496	0	0
0	257	504	0	0
0	256	499	0	0
23	263	467	0	0
23	263	465	0	0
23	260	467	0	0
63	261	469	6	0
23	260	465	0	0
19	261	459	4	0
23	263	463	0	0
23	263	461	0	0
23	260	463	0	0
23	260	461	0	0
15	249	467	6	0
198	251	468	0	0
26	255	462	6	0
26	255	466	6	0
15	249	463	6	0
26	255	458	6	0
15	249	459	6	0
1	253	464	1	1
51	270	3327	0	0
51	270	3332	4	0
151	266	3337	2	0
153	278	3335	2	0
153	268	3337	0	0
153	266	3343	2	0
153	273	3337	0	0
151	278	3333	2	0
103	274	3340	0	0
5	279	3326	0	0
153	281	3342	2	0
51	283	3330	6	0
153	266	3342	2	0
29	269	3328	2	0
150	281	3326	6	0
51	268	3336	0	0
153	281	3339	2	0
153	281	3343	2	0
51	283	3335	6	0
100	276	3341	2	0
51	280	3325	0	0
153	281	3340	2	0
153	272	3337	0	0
29	291	3329	0	0
51	290	3327	0	0
51	289	3330	2	0
3	294	3328	6	0
153	269	3337	0	0
153	266	3341	2	0
101	277	3343	0	0
151	277	3336	6	0
153	271	3337	0	0
153	266	3340	2	0
153	275	3337	0	0
153	266	3339	2	0
51	276	3331	2	0
51	267	3328	2	0
153	270	3337	0	0
153	276	3337	0	0
104	274	3342	2	0
153	281	3341	2	0
153	274	3337	0	0
51	272	3336	0	0
152	280	3332	6	0
153	281	3328	6	0
149	281	3330	6	0
149	281	3329	6	0
153	281	3338	2	0
153	281	3334	2	0
153	281	3335	2	0
153	281	3327	6	0
153	281	3331	6	0
29	293	3339	6	0
153	281	3336	2	0
51	290	3337	0	0
153	281	3337	2	0
51	293	3327	0	0
51	293	3337	0	0
102	261	3342	2	0
51	262	3333	0	0
51	261	3340	0	0
51	256	3336	4	0
105	260	3341	0	0
51	260	3333	0	0
57	259	3334	1	1
51	259	3344	2	0
114	261	3346	2	0
51	261	3350	4	0
51	271	3350	6	0
104	260	3349	2	0
150	266	3346	2	0
149	266	3345	2	0
50	283	3349	0	0
153	281	3344	2	0
150	281	3348	2	0
50	286	3346	0	0
3	293	3351	6	0
101	277	3351	0	0
153	281	3345	2	0
51	294	3350	6	0
97	290	3347	6	0
114	277	3349	2	0
50	283	3346	0	0
115	277	3348	0	0
153	281	3347	2	0
51	276	3346	2	0
50	286	3349	0	0
153	266	3344	2	0
97	290	3348	6	0
97	291	3347	6	0
97	291	3348	6	0
51	292	3344	0	0
149	281	3346	2	0
51	254	3335	2	0
51	254	3332	2	0
230	255	3331	6	0
51	272	3355	6	0
100	278	3353	2	0
280	291	3356	0	0
51	279	3354	4	0
111	265	3357	0	0
51	294	3356	6	0
110	264	3358	2	0
51	290	3358	4	0
3	301	3350	6	0
51	301	3346	6	0
51	299	3351	4	0
51	299	3344	0	0
101	270	3366	0	0
101	269	3366	0	0
100	269	3364	0	0
104	264	3366	0	0
51	305	3347	0	0
51	311	3351	4	0
281	305	3358	6	0
51	305	3349	4	0
51	311	3343	0	0
102	272	3369	0	0
103	272	3370	0	0
109	276	3373	0	0
51	313	3345	6	0
51	314	3349	6	0
51	327	3361	0	0
299	325	3367	6	0
51	325	3364	2	0
46	325	3366	4	0
104	263	3367	0	0
104	263	3366	0	0
107	271	3371	0	0
43	251	3369	6	0
102	271	3368	0	0
103	270	3375	0	0
51	266	3380	4	0
51	267	3383	2	0
110	265	3376	0	0
110	266	3376	0	0
111	270	3380	0	0
111	272	3379	0	0
113	264	3380	0	0
110	274	3377	0	0
111	272	3380	0	0
55	268	3381	0	1
111	270	3390	0	0
110	269	3391	0	0
110	267	3389	0	0
110	269	3390	0	0
51	271	3387	6	0
110	266	3387	0	0
110	265	3387	0	0
110	265	3391	0	0
111	272	3391	0	0
111	271	3397	0	0
111	273	3396	0	0
111	275	3396	0	0
111	275	3395	0	0
51	277	3398	6	0
107	267	3399	0	0
51	263	3392	2	0
51	277	3394	6	0
107	268	3399	0	0
111	272	3392	0	0
110	268	3396	0	0
111	272	3394	0	0
111	272	3397	0	0
51	263	3395	2	0
110	266	3394	0	0
110	265	3394	0	0
111	273	3394	0	0
110	267	3396	0	0
107	265	3397	0	0
5	274	3398	6	0
70	254	477	6	0
38	251	498	0	0
0	250	501	0	0
0	251	503	0	0
70	249	475	6	0
70	252	473	6	0
0	252	494	0	0
38	254	502	0	0
70	232	484	0	0
70	233	481	0	0
97	233	495	0	0
97	234	495	0	0
1	245	452	0	0
1	244	449	0	0
1	243	452	0	0
1	242	448	0	0
1	240	450	0	0
45	250	1412	6	0
25	265	1412	6	0
25	262	1406	6	0
45	250	1411	6	0
15	249	1406	6	0
6	251	1412	0	0
15	249	1402	6	0
200	264	1409	6	0
25	262	1412	6	0
25	265	1406	6	0
1	259	1403	1	1
1	260	1406	0	1
1	251	1405	0	1
1	251	1408	0	1
1	253	1403	1	1
45	275	1390	4	0
45	276	1387	0	0
45	276	1390	4	0
44	273	1388	6	0
45	273	1390	4	0
45	277	1388	2	0
45	273	1387	0	0
6	273	1384	2	0
45	277	1389	2	0
45	277	1390	2	0
45	274	1387	0	0
45	277	1387	2	0
45	274	1390	4	0
45	275	1387	0	0
45	272	1390	6	0
45	272	1387	6	0
30	279	1386	0	1
3	232	497	0	0
3	232	500	0	0
63	232	503	6	0
51	235	498	6	0
3	234	500	6	0
3	233	500	0	0
3	233	497	0	0
3	234	497	6	0
51	235	500	6	0
0	248	504	0	0
38	254	506	0	0
0	252	507	0	0
8	232	509	0	0
8	232	510	1	0
3	234	507	2	0
70	237	504	6	0
1	236	511	0	1
3	237	519	1	0
1	235	518	2	1
70	230	486	0	0
0	230	483	0	0
0	229	480	0	0
0	225	481	0	0
3	231	500	0	0
70	230	489	0	0
0	227	486	0	0
70	230	493	6	0
22	230	495	6	0
0	227	484	0	0
0	224	485	0	0
8	228	497	0	0
70	227	501	6	0
51	230	500	2	0
8	228	498	7	0
110	225	504	4	0
104	229	518	4	0
104	229	517	4	0
104	228	516	4	0
7	229	508	4	0
121	229	509	2	0
70	226	493	6	0
110	227	507	4	0
51	230	498	2	0
3	231	497	0	0
110	225	505	4	0
3	231	508	4	0
110	228	505	4	0
34	224	518	1	0
1	230	511	0	1
70	228	533	0	0
70	224	527	0	0
179	227	524	6	0
178	229	523	2	0
179	228	524	6	0
70	231	527	0	0
1	228	521	0	1
34	222	516	1	0
34	219	515	1	0
34	221	509	1	0
70	221	520	0	0
70	219	522	0	0
70	219	526	0	0
70	216	524	0	0
70	223	525	0	0
70	223	522	0	0
70	221	524	0	0
192	208	501	4	0
45	214	511	0	0
45	214	513	4	0
70	208	532	0	0
45	209	511	0	0
70	208	534	0	0
45	213	513	4	0
70	215	531	0	0
45	212	513	4	0
70	208	527	0	0
45	213	511	0	0
45	211	513	4	0
70	212	527	0	0
45	210	511	0	0
45	209	513	4	0
45	212	511	0	0
45	210	513	4	0
192	212	507	4	0
45	211	511	0	0
0	220	481	0	0
20	222	472	0	0
20	220	472	0	0
1	225	464	0	0
1	218	465	0	1
1	219	467	1	1
34	222	458	0	0
203	219	463	0	0
34	222	457	0	0
203	219	461	0	0
20	220	459	6	0
4	225	462	0	0
20	218	459	6	0
70	226	461	6	0
20	216	459	0	0
34	219	457	0	0
1	225	460	0	0
285	229	448	2	0
285	224	448	2	0
15	226	455	2	0
285	228	449	2	0
285	224	449	2	0
15	226	452	2	0
70	224	450	6	0
1	225	452	0	1
29	219	451	2	0
36	211	471	0	0
203	214	465	4	0
6	215	468	0	0
36	213	471	0	0
20	211	472	0	0
36	212	474	0	0
0	209	481	0	0
36	211	467	0	0
36	210	470	0	0
20	209	472	0	0
36	211	470	0	0
36	210	469	0	0
36	208	470	0	0
34	212	457	0	0
203	214	461	4	0
202	214	463	4	0
20	214	459	0	0
34	201	494	0	0
0	206	491	0	0
34	205	492	0	0
1	203	488	0	0
34	193	495	0	0
34	194	488	0	0
34	199	491	0	0
34	196	491	0	0
1	197	488	0	0
0	199	502	0	0
34	195	501	0	0
34	194	498	0	0
0	200	498	0	0
34	204	497	0	0
0	196	509	0	0
0	200	507	0	0
34	202	504	0	0
88	204	540	0	0
70	215	539	0	0
70	200	540	0	0
70	206	538	0	0
70	204	536	0	0
48	220	542	6	0
70	216	536	0	0
70	218	540	0	0
70	225	542	0	0
70	229	538	0	0
70	227	541	0	0
70	224	537	0	0
88	223	541	0	0
70	209	540	0	0
88	210	537	0	0
70	221	538	0	0
1	217	543	1	1
70	199	533	0	0
70	195	536	0	0
70	194	540	0	0
3	197	551	0	0
70	199	544	0	0
25	213	551	4	0
45	206	550	6	0
117	212	549	2	0
45	206	549	4	0
3	219	549	6	0
45	205	549	4	0
42	204	551	0	0
47	220	550	4	0
55	201	551	0	0
3	218	549	6	0
45	203	549	4	0
11	221	544	4	0
70	193	547	0	0
70	193	546	0	0
70	194	549	0	0
25	208	551	4	0
45	206	551	6	0
83	221	549	0	0
43	210	547	4	0
45	204	549	4	0
45	203	550	2	0
45	203	551	2	0
24	208	547	2	1
24	221	550	2	1
1	220	547	1	1
35	212	545	1	1
1	213	545	0	1
37	199	551	0	1
1	215	550	0	1
88	197	545	0	0
70	229	544	0	0
83	220	549	0	0
88	200	547	0	0
11	222	1385	0	0
6	226	1383	0	0
15	226	1386	2	0
51	204	3293	0	0
51	204	3291	2	0
51	207	3282	0	0
51	206	3294	6	0
51	204	3285	2	0
51	204	3313	0	0
51	206	3296	6	0
51	210	3307	6	0
55	211	3286	6	0
51	211	3285	6	0
51	221	3283	4	0
51	218	3286	6	0
5	203	3314	4	0
51	208	3298	0	0
51	210	3301	6	0
51	216	3280	2	0
51	214	3296	2	0
51	216	3286	2	0
51	202	3300	4	0
51	200	3308	4	0
51	208	3292	4	0
51	219	3291	6	0
5	215	3300	0	0
51	200	3303	0	0
51	210	3304	6	0
51	204	3316	4	0
51	210	3312	6	0
51	216	3298	6	0
51	215	3289	0	0
51	214	3298	2	0
51	216	3296	6	0
57	208	3317	2	0
51	210	3315	6	0
22	219	3282	1	1
51	198	3277	4	0
104	194	3294	4	0
110	194	3295	4	0
51	195	3277	4	0
51	194	3272	2	0
195	192	3294	4	0
51	198	3299	4	0
51	198	3297	0	0
110	194	3301	0	0
110	193	3302	4	0
110	192	3302	4	0
51	196	3305	2	0
51	204	3272	0	0
51	214	3276	4	0
51	194	3300	6	0
51	207	3274	4	0
195	192	3298	0	0
51	207	3272	0	0
51	221	3279	0	0
51	194	3296	6	0
51	204	3274	4	0
20	219	3276	0	0
58	211	3272	0	0
51	218	3274	6	0
20	219	3273	0	0
20	213	3270	0	0
305	196	3266	6	0
51	218	3271	6	0
20	216	3270	0	0
51	201	3256	6	0
143	209	3257	4	0
143	212	3257	4	0
209	200	3258	5	0
209	198	3258	6	0
51	196	3262	2	0
209	196	3256	1	0
51	190	3274	0	0
51	186	3274	0	0
51	184	3282	4	0
51	187	3288	2	0
51	187	3294	0	0
51	189	3277	6	0
106	191	3294	4	0
51	189	3291	6	0
51	189	3288	6	0
51	185	3293	4	0
51	186	3291	0	0
51	178	3281	2	0
51	181	3280	0	0
51	182	3274	0	0
51	178	3287	2	0
51	178	3290	2	0
51	181	3293	4	0
55	195	3252	6	0
80	203	3248	4	0
51	195	3251	2	0
80	206	3248	4	0
51	195	3254	2	0
143	210	3250	0	0
55	200	3250	7	0
209	202	3253	4	0
209	201	3255	4	0
51	182	3302	4	0
51	182	3300	0	0
102	189	3297	4	0
108	190	3298	4	0
51	189	3302	4	0
104	187	3302	4	0
51	180	3302	4	0
51	180	3300	0	0
100	186	3296	4	0
102	186	3297	4	0
98	191	3302	0	0
58	186	3300	0	0
51	205	3325	2	0
51	202	3327	0	0
51	205	3322	2	0
51	211	3326	6	0
51	212	3322	6	0
51	199	3332	2	0
51	199	3329	2	0
51	203	3335	4	0
51	212	3329	6	0
51	208	3333	4	0
51	198	3247	6	0
51	198	3244	6	0
57	217	3241	6	0
143	216	3254	6	0
143	219	3249	6	0
143	219	3246	6	0
51	222	3232	0	0
304	230	3248	6	0
51	229	3236	4	0
51	227	3247	2	0
51	227	3244	2	0
51	234	3233	6	0
51	233	3242	6	0
51	232	3236	4	0
51	233	3247	6	0
51	206	3227	2	0
684	207	3224	0	0
598	206	3224	0	0
209	206	3226	6	0
598	205	3224	0	0
5	231	3224	4	0
51	229	3229	0	0
598	208	3224	0	0
51	215	3230	0	0
58	210	3231	4	0
51	208	3227	6	0
209	208	3226	6	0
598	207	3222	0	0
685	206	3222	0	0
598	208	3222	0	0
598	205	3222	0	0
51	210	3218	6	0
51	204	3218	2	0
55	205	3220	0	0
55	209	3217	0	0
51	210	3211	6	0
51	203	3211	2	0
51	203	3208	2	0
51	206	3206	0	0
34	182	501	0	0
34	179	499	0	0
88	213	558	0	0
63	210	556	6	0
70	194	554	6	0
45	206	553	6	0
70	203	558	6	0
45	206	552	6	0
3	215	553	2	0
45	206	554	0	0
7	214	552	6	0
25	208	553	4	0
3	215	552	2	0
116	208	552	2	0
70	199	558	6	0
25	213	553	4	0
3	197	552	0	0
45	203	552	2	0
45	203	554	0	0
3	198	555	0	0
3	197	555	0	0
45	203	553	2	0
55	201	555	0	0
25	213	555	4	0
70	194	556	6	0
25	208	555	4	0
1	202	553	1	1
36	210	553	0	1
70	211	567	6	0
88	205	563	0	0
38	202	565	6	0
70	214	565	0	0
70	198	563	0	0
70	204	560	6	0
88	208	564	0	0
70	211	563	6	0
70	192	566	0	0
38	206	560	6	0
70	203	567	0	0
70	204	565	0	0
38	212	563	6	0
70	201	562	0	0
88	210	561	0	0
70	210	566	6	0
70	194	572	0	0
70	199	571	0	0
70	197	572	0	0
70	207	572	0	0
70	207	575	0	0
70	200	573	0	0
70	210	573	6	0
70	211	574	6	0
70	204	574	0	0
70	196	569	0	0
70	198	568	0	0
70	203	571	0	0
70	206	568	0	0
70	199	569	0	0
70	198	575	0	0
38	209	574	6	0
70	202	575	0	0
38	213	568	6	0
58	210	570	6	0
70	198	580	6	0
70	195	583	6	0
70	196	576	6	0
70	195	578	6	0
70	193	577	6	0
70	192	581	6	0
4	199	601	0	0
70	200	578	6	0
70	206	582	6	0
70	204	576	6	0
70	206	580	6	0
70	202	580	6	0
0	204	600	0	0
4	201	604	0	0
0	202	596	0	0
70	203	589	0	0
70	206	585	0	0
134	230	552	0	0
70	228	552	6	0
70	216	559	6	0
47	220	554	4	0
71	216	555	2	0
70	224	558	6	0
70	228	559	6	0
6	223	553	0	0
24	221	555	3	1
70	219	565	6	0
70	218	561	6	0
38	221	563	6	0
70	223	565	6	0
70	216	567	0	0
86	226	565	0	0
70	227	562	6	0
70	229	567	6	0
38	229	561	6	0
70	229	564	6	0
70	217	573	0	0
70	220	574	0	0
70	228	572	6	0
70	223	572	6	0
70	224	568	6	0
70	228	575	6	0
70	226	573	6	0
70	230	579	6	0
70	230	583	6	0
70	231	577	6	0
70	227	576	6	0
70	228	584	6	0
1	252	613	0	0
1	256	613	0	0
70	223	576	6	0
70	223	582	6	0
55	220	614	0	0
70	219	612	0	0
55	221	615	0	0
34	218	613	0	0
70	217	582	6	0
70	218	578	6	0
70	216	581	6	0
7	216	610	0	0
34	217	611	0	0
3	216	609	4	0
70	218	610	0	0
70	222	579	6	0
5	216	620	0	0
54	222	618	0	0
117	219	620	0	0
70	224	622	0	0
1	249	619	0	0
1	226	617	0	0
1	252	617	0	0
1	219	616	3	1
1	220	620	0	1
70	213	583	6	0
70	213	589	0	0
0	214	597	0	0
70	212	592	0	0
70	209	576	6	0
70	215	577	6	0
70	211	581	6	0
70	210	580	6	0
70	210	581	6	0
70	208	582	6	0
70	212	578	6	0
70	210	576	6	0
70	208	579	6	0
0	215	600	0	0
87	208	604	4	0
15	213	609	4	0
34	212	611	0	0
70	212	615	0	0
62	212	609	0	0
66	208	609	6	0
1	215	613	0	1
72	208	618	3	0
70	210	620	0	0
3	215	620	4	0
7	215	621	0	0
1	215	618	0	1
72	208	619	3	0
72	208	620	3	0
70	211	619	0	0
25	207	3382	0	0
25	207	3379	0	0
51	207	3383	6	0
51	207	3381	6	0
51	206	3379	0	0
51	203	3379	0	0
51	202	3381	2	0
51	202	3383	2	0
25	202	3382	0	0
136	204	3380	2	0
41	204	3383	0	0
125	222	3382	6	0
25	202	3379	0	0
51	207	3385	6	0
51	222	3376	2	0
51	229	3380	4	0
126	222	3378	6	0
129	230	3376	2	0
51	230	3377	6	0
51	230	3379	6	0
128	229	3375	4	0
51	230	3381	0	0
127	223	3375	4	0
51	202	3385	2	0
51	222	3384	2	0
51	230	3383	4	0
130	223	3385	0	0
51	227	3384	6	0
124	225	3386	0	0
51	226	3375	0	0
51	225	3378	0	0
51	227	3386	4	0
32	223	3381	0	1
31	229	3378	0	1
25	225	3376	1	1
28	228	3379	1	1
29	228	3382	1	1
27	225	3379	1	1
30	226	3378	0	1
33	226	3381	0	1
26	228	3376	1	1
25	202	1486	0	0
25	206	1486	0	0
45	209	1493	2	0
25	205	1494	0	0
25	215	1486	0	0
45	212	1491	6	0
117	204	1495	0	0
45	212	1490	4	0
25	210	1486	0	0
51	213	1498	2	0
45	210	1495	0	0
20	222	1498	6	0
45	209	1494	2	0
27	209	1489	0	0
20	221	1496	6	0
117	202	1497	2	0
45	212	1495	0	0
45	209	1492	2	0
25	215	1488	0	0
25	210	1488	0	0
20	217	1495	6	0
44	210	1491	4	0
45	211	1495	0	0
25	202	1489	0	0
45	212	1494	6	0
25	206	1499	2	0
45	209	1491	2	0
45	212	1493	6	0
5	215	1492	0	0
45	209	1490	4	0
24	211	1489	0	0
45	212	1492	6	0
45	209	1495	0	0
25	202	1494	0	0
14	202	1498	6	0
20	219	1497	6	0
51	219	1499	4	0
20	216	1498	6	0
51	213	1496	2	0
51	224	1496	6	0
25	208	1495	2	0
1	207	1495	0	1
24	202	1495	2	1
1	206	1491	1	1
1	209	1497	1	1
24	208	1499	2	1
24	215	1494	2	1
3	215	2438	0	0
6	215	2436	0	0
3	208	2438	0	0
17	202	2441	0	0
120	212	2439	0	0
1	211	2438	0	1
116	202	2438	0	1
1	278	593	0	0
1	272	595	0	0
4	275	594	0	0
4	278	597	0	0
37	268	605	2	0
60	274	603	0	0
37	264	606	2	0
37	268	607	2	0
37	264	605	2	0
37	264	607	2	0
2	271	602	1	1
1	267	605	0	1
1	300	607	2	0
1	295	606	2	0
34	295	602	2	0
1	302	602	2	0
1	309	600	2	0
1	307	602	2	0
1	305	607	2	0
1	296	600	2	0
1	294	615	2	0
1	297	614	2	0
0	301	615	2	0
1	261	617	0	0
1	256	619	0	0
7	257	631	2	0
160	257	630	6	0
3	274	630	0	0
16	256	626	4	0
27	278	630	0	0
48	273	625	6	0
29	266	628	0	0
24	279	630	0	0
7	277	627	0	0
5	277	625	0	0
3	256	631	0	0
16	256	628	4	0
29	264	628	0	0
16	256	627	4	0
8	280	627	0	0
1	262	630	1	1
1	275	628	0	1
3	254	628	4	0
7	255	631	6	0
7	253	629	0	0
3	254	627	4	0
7	253	630	2	0
7	251	630	6	0
7	250	626	2	0
45	249	631	2	0
45	249	630	2	0
5	249	628	2	0
3	253	627	4	0
3	252	630	2	0
3	249	626	4	0
3	253	628	4	0
7	253	626	4	0
21	248	630	4	0
7	249	625	4	0
1	250	630	0	1
1	252	624	0	1
34	303	617	2	0
1	298	619	2	0
34	306	613	2	0
1	310	608	2	0
0	311	618	2	0
1	308	613	2	0
1	307	618	2	0
4	308	631	2	0
70	310	631	1	0
1	306	616	2	0
1	314	603	2	0
1	314	619	2	0
1	313	611	2	0
34	312	615	2	0
0	316	610	2	0
0	319	613	2	0
34	321	610	2	0
0	326	611	2	0
34	323	615	2	0
34	322	613	2	0
1	323	605	2	0
3	343	601	2	0
196	341	602	0	0
3	343	606	0	0
179	343	605	6	0
115	340	607	0	0
196	339	601	0	0
115	338	607	0	0
195	338	603	4	0
196	338	601	0	0
195	337	604	4	0
112	340	613	4	0
114	339	609	4	0
3	343	608	0	0
113	337	612	0	0
20	343	614	0	0
112	339	613	4	0
114	339	608	4	0
112	338	614	4	0
113	337	611	0	0
178	345	601	6	0
3	344	601	2	0
5	349	612	0	0
20	346	614	6	0
3	344	606	0	0
68	347	601	0	1
59	352	600	2	0
47	368	578	6	0
15	370	578	4	0
3	368	580	4	0
1	372	580	1	1
5	360	570	0	0
2	360	568	2	1
2	364	569	3	1
2	363	573	2	1
20	349	625	6	0
6	346	625	0	0
20	345	625	0	0
6	338	626	0	0
5	347	627	0	0
48	349	629	1	1
54	340	630	1	1
0	330	639	0	0
5	350	638	0	0
5	350	633	0	0
201	344	638	0	0
48	349	635	1	1
48	349	632	1	1
1	346	637	0	1
60	338	632	1	1
1	315	633	2	0
1	321	639	2	0
70	319	636	1	0
100	312	636	0	0
30	230	631	6	0
30	226	628	0	0
70	230	627	0	0
70	224	630	2	0
55	217	624	0	0
5	218	630	2	0
3	223	624	0	0
3	222	624	0	0
1	222	625	1	1
1	219	628	0	1
100	311	636	0	0
100	309	635	0	0
4	304	633	2	0
100	310	635	0	0
114	306	639	0	0
4	301	638	2	0
1	302	635	1	0
1	285	639	0	0
1	286	634	0	0
295	280	634	0	0
1	282	637	0	0
1	280	639	0	1
190	276	632	4	0
11	275	638	0	0
7	278	639	4	0
29	278	632	6	0
1	272	637	0	0
3	274	634	0	0
1	276	637	0	1
85	277	632	1	1
21	267	639	4	0
21	271	639	2	0
1	265	635	0	0
1	266	633	0	1
1	271	632	1	1
37	263	635	0	0
37	263	634	0	0
45	259	633	2	0
37	262	634	0	0
45	258	635	0	0
45	259	634	2	0
45	258	632	4	0
45	257	632	4	0
45	256	632	4	0
21	256	633	6	0
45	253	632	4	0
21	252	633	6	0
45	249	632	2	0
21	248	632	4	0
45	255	632	4	0
45	252	632	4	0
45	251	632	4	0
21	249	633	0	0
45	254	632	4	0
55	262	640	4	0
55	261	640	2	0
224	254	643	4	0
21	257	640	6	0
55	264	640	6	0
21	267	643	6	0
21	271	643	0	0
155	265	645	2	0
21	257	642	4	0
225	259	643	4	0
1	254	1572	1	1
1	251	1571	0	1
15	256	1572	2	0
3	252	1568	2	0
1	253	1571	0	1
15	253	1568	4	0
15	249	1568	4	0
6	249	1572	2	0
3	257	1571	2	0
3	257	1568	2	0
29	277	647	4	0
3	278	640	0	0
1	274	643	0	0
3	276	645	4	0
3	279	641	0	0
71	278	645	6	0
7	279	642	0	0
7	279	640	4	0
21	271	652	0	0
21	271	649	2	0
156	265	650	2	0
3	279	655	0	0
157	265	652	2	0
159	276	654	2	0
189	278	651	6	0
22	278	655	0	0
5	277	655	0	0
1	275	655	0	1
1	277	651	0	1
5	282	642	0	0
3	282	640	0	0
1	287	645	0	0
1	287	640	0	0
1	283	644	0	0
1	284	640	0	0
1	286	646	0	0
1	286	642	0	0
0	287	650	0	0
0	282	653	0	0
1	285	648	0	0
0	283	648	0	0
0	283	651	0	0
1	285	655	0	1
34	290	653	0	0
34	294	649	0	0
0	294	654	0	0
0	290	649	0	0
1	289	655	0	0
1	292	644	0	0
37	292	654	0	0
1	289	647	0	0
1	294	646	0	0
406	291	651	4	0
47	279	1586	2	0
47	280	1575	4	0
6	282	1586	0	0
6	277	1569	0	0
6	280	1578	0	0
37	302	650	0	0
1	297	647	0	0
1	301	652	0	0
1	296	652	0	0
38	297	654	0	0
1	300	642	1	0
0	298	648	0	0
70	303	643	1	0
1	301	649	0	0
34	296	650	0	0
3	278	1599	6	0
3	279	1599	6	0
239	264	661	0	0
1	291	662	0	0
238	266	661	0	0
0	294	658	0	0
22	265	658	6	0
1	298	659	0	0
0	301	656	0	0
34	291	659	0	0
0	296	660	0	0
1	293	658	0	0
1	297	656	0	0
1	301	659	0	0
0	288	660	0	0
5	276	659	0	0
55	278	659	0	0
29	275	657	0	0
185	279	658	0	0
22	264	658	6	0
3	279	656	0	0
55	279	659	0	0
30	285	656	1	1
1	273	656	1	1
30	285	660	1	1
47	277	658	1	1
6	276	1603	2	0
45	277	1603	2	0
45	277	1602	2	0
15	273	1602	0	0
1	296	667	0	0
0	298	667	0	0
0	299	671	0	0
1	289	671	0	0
1	296	670	0	0
1	302	670	0	0
1	284	670	0	0
0	289	669	0	0
1	293	670	0	0
1	302	665	0	0
34	300	668	0	0
34	286	670	0	0
34	291	671	0	0
34	298	670	0	0
30	281	666	1	1
30	283	665	0	1
30	281	665	0	1
22	263	658	6	0
239	262	661	0	0
21	260	658	6	0
21	260	660	4	0
240	257	661	0	0
1	294	679	2	0
0	301	684	4	0
1	295	686	4	0
0	292	687	2	0
1	292	684	2	0
1	298	682	2	0
1	299	685	4	0
1	297	687	4	0
102	311	644	0	0
104	308	641	0	0
1	304	656	0	0
4	305	644	2	0
34	307	653	0	0
4	307	647	2	0
1	309	664	0	0
37	310	669	0	0
34	304	668	0	0
0	307	670	0	0
1	305	665	0	0
1	307	666	0	0
0	307	664	0	0
25	303	688	6	0
23	302	690	6	0
0	291	695	2	0
1	295	688	4	0
0	294	689	4	0
63	298	689	4	0
23	300	690	6	0
1	298	695	2	0
0	292	693	2	0
25	301	688	6	0
6	285	711	0	0
0	285	710	0	0
1	283	710	0	0
50	293	711	0	0
1	286	708	0	0
0	286	710	0	0
34	284	713	0	0
34	287	712	0	0
1	286	713	0	0
34	286	715	0	0
1	284	712	0	0
8	292	715	0	0
15	293	713	2	0
1	292	713	1	1
193	303	727	6	0
193	302	727	6	0
1093	294	727	0	0
193	300	724	0	0
193	296	724	0	0
193	297	726	0	0
193	303	726	6	0
1179	289	728	4	0
193	297	732	0	0
193	294	734	0	0
1093	292	728	0	0
1094	292	730	0	0
1180	301	728	4	0
193	301	734	0	0
1101	300	728	6	0
1181	303	728	4	0
22	295	729	0	0
1094	294	731	0	0
193	303	731	6	0
1102	300	730	6	0
193	303	732	0	0
193	302	731	6	0
22	298	729	0	0
22	297	729	0	0
51	275	3525	2	0
5	285	3543	0	0
51	283	3524	6	0
51	282	3528	6	0
51	286	3542	6	0
51	286	3544	6	0
51	278	3532	2	0
51	282	3539	0	0
98	302	3522	2	0
51	275	3521	2	0
51	280	3534	6	0
51	295	3523	6	0
51	279	3541	2	0
51	292	3524	4	0
51	284	3513	0	0
51	294	3515	6	0
51	291	3513	0	0
51	303	3513	0	0
51	280	3519	0	0
98	301	3514	2	0
98	310	3517	6	0
51	308	3515	6	0
98	306	3526	6	0
98	307	3514	2	0
51	309	3528	4	0
98	308	3517	2	0
98	313	3517	6	0
176	312	3517	6	0
51	316	3523	6	0
51	315	3519	6	0
51	315	3526	6	0
99	312	3528	6	0
98	316	3522	6	0
176	314	3527	6	0
19	304	690	2	0
25	305	688	6	0
1	305	680	2	0
0	304	686	4	0
193	311	686	0	0
1	305	684	2	0
3	319	659	0	0
11	314	663	0	0
41	315	659	6	0
20	315	661	0	0
1	318	649	2	0
147	316	666	0	0
34	313	668	0	0
193	313	681	0	0
48	316	664	2	0
1	314	659	0	1
1	319	667	1	1
37	326	657	0	0
37	326	659	0	0
2	323	662	4	0
37	323	655	0	0
37	326	655	0	0
37	323	657	0	0
37	323	659	0	0
1	320	662	1	1
1	324	667	0	1
163	319	710	0	0
162	324	710	0	0
161	326	710	0	0
55	332	666	0	0
29	331	662	2	0
3	330	670	0	0
29	329	668	0	0
89	328	659	4	0
3	332	670	0	0
29	330	658	0	0
1	329	661	1	1
1	328	666	2	1
37	341	660	0	0
47	342	662	0	0
22	341	663	2	0
37	341	662	0	0
37	343	660	2	0
37	340	663	0	0
37	340	660	0	0
306	338	657	0	0
0	340	659	0	0
3	340	665	2	0
0	338	668	0	0
0	337	663	0	0
7	341	665	2	0
37	347	660	1	0
37	347	662	0	0
37	347	663	0	0
37	346	660	2	0
34	350	656	4	0
7	346	662	2	0
0	344	670	0	0
7	348	664	2	0
37	345	660	0	0
0	347	669	0	0
3	349	664	2	0
1	344	662	0	1
1	347	665	1	1
102	318	642	0	0
102	318	640	0	0
70	322	646	1	0
112	315	645	0	0
0	329	651	0	0
0	328	649	0	0
0	331	651	0	0
0	335	643	0	0
0	334	654	0	0
0	332	643	0	0
0	342	655	0	0
306	339	651	0	0
0	337	654	0	0
0	337	649	0	0
34	349	654	4	0
0	347	651	0	0
34	351	651	4	0
306	345	653	0	0
0	348	655	0	0
34	353	653	4	0
34	354	648	4	0
34	355	650	4	0
45	319	1606	0	0
45	317	1606	0	0
45	318	1606	0	0
27	314	1603	0	0
42	315	1603	6	0
45	315	1605	2	0
45	315	1606	0	0
45	316	1606	0	0
1	361	563	0	0
36	361	534	0	0
36	360	535	0	0
36	363	532	0	0
1	363	539	0	0
36	364	531	0	0
103	367	554	0	0
1	367	547	0	0
1	365	543	0	0
0	365	544	0	0
105	363	551	0	0
1	365	552	0	0
103	366	553	0	0
111	360	545	0	0
101	361	549	0	0
103	365	556	0	0
101	361	544	0	0
105	362	550	0	0
105	362	551	0	0
105	362	549	0	0
1	364	516	0	0
1	367	516	0	0
1	367	513	0	0
1	370	555	0	0
1	371	546	0	0
1	369	539	0	0
70	368	529	0	0
70	375	529	0	0
34	372	530	0	0
70	374	533	0	0
34	374	527	4	0
70	374	523	4	0
70	375	520	4	0
34	370	525	4	0
70	381	536	0	0
70	378	534	0	0
70	376	534	0	0
34	380	529	0	0
34	383	529	0	0
70	376	522	4	0
6	376	520	6	0
70	378	522	4	0
70	379	521	4	0
34	372	516	4	0
34	376	512	4	0
70	375	518	4	0
34	380	513	4	0
34	382	514	4	0
70	376	519	4	0
0	397	547	2	0
36	399	554	0	0
36	398	553	0	0
0	398	560	6	0
5	363	1514	0	0
6	360	1514	0	0
47	361	2461	2	0
47	359	2458	0	0
47	361	2456	6	0
6	363	2458	0	0
251	371	505	0	0
29	366	507	2	0
29	366	504	2	0
1	371	506	1	1
34	364	502	4	0
37	362	498	4	0
3	360	502	3	0
248	367	497	0	0
3	381	500	0	0
250	376	502	4	0
5	380	503	0	0
1	361	499	3	1
1	377	501	1	1
1	361	489	4	0
48	362	495	4	0
5	367	494	0	0
7	360	494	0	0
0	360	489	4	0
3	360	493	0	0
6	362	491	4	0
1	368	489	4	0
37	368	488	4	0
20	360	495	0	0
255	363	494	2	0
1	371	495	1	1
69	363	494	1	1
0	360	486	4	0
1	367	482	4	0
1	365	485	4	0
1	360	482	4	0
68	362	467	4	0
68	362	457	4	0
68	367	462	0	0
235	362	462	0	0
0	362	455	0	0
1	363	454	0	0
37	360	448	0	0
1	360	452	0	0
0	362	454	0	0
0	362	442	0	0
0	363	444	0	0
0	363	446	0	0
5	375	438	0	0
278	372	435	0	0
280	372	436	0	0
278	374	436	0	0
278	373	438	0	0
37	373	451	0	0
42	368	438	4	0
1	368	442	0	0
45	370	438	2	0
45	370	439	2	0
278	371	436	2	0
274	369	434	4	0
45	370	440	2	0
74	372	441	0	1
34	375	460	0	0
37	369	458	0	0
34	373	464	0	0
37	380	483	4	0
1	378	482	4	0
1	378	484	4	0
1	380	482	4	0
1	381	483	4	0
0	376	482	4	0
37	381	492	4	0
1	378	493	2	1
1	377	488	1	1
286	382	455	0	0
286	379	449	0	0
286	381	450	0	0
286	380	446	0	0
37	382	436	0	0
37	382	438	0	0
282	381	439	0	0
0	381	435	0	0
0	380	439	0	0
37	379	439	0	0
0	379	435	0	0
1	378	437	1	1
45	384	468	6	0
45	387	467	2	0
45	386	469	0	0
45	387	466	2	0
63	385	465	2	0
45	384	467	6	0
359	385	466	4	0
45	385	469	0	0
45	387	468	2	0
45	384	466	6	0
261	399	503	0	0
34	399	499	6	0
6	399	433	0	0
6	402	435	2	0
164	403	461	0	0
164	404	462	0	0
618	402	463	0	0
70	402	436	2	0
6	402	438	2	0
70	402	434	2	0
6	402	432	2	0
164	401	463	0	0
164	402	464	0	0
164	403	464	0	0
0	403	496	6	0
34	406	499	6	0
261	406	505	0	0
261	398	505	4	0
261	402	507	0	0
1	406	535	0	0
1	404	535	0	0
1	402	539	0	0
1	404	542	0	0
1	407	538	0	0
1	400	542	2	0
1	400	538	2	0
45	360	1439	0	0
45	360	1438	6	0
6	367	1438	0	0
3	362	1437	2	0
22	381	1444	4	0
6	380	1447	0	0
3	377	1446	4	0
257	362	1440	0	0
45	360	1435	6	0
45	360	1436	6	0
45	360	1437	6	0
42	358	1436	0	0
143	347	3317	6	0
217	344	3317	5	0
143	344	3314	0	0
145	350	3326	2	0
51	346	3325	2	0
51	360	3318	0	0
145	364	3321	0	0
51	365	3319	0	0
51	346	3330	2	0
145	350	3332	2	0
51	344	3342	2	0
51	357	3341	4	0
51	369	3329	0	0
145	360	3339	4	0
145	346	3343	2	0
145	364	3339	4	0
51	371	3319	0	0
236	369	3332	0	0
51	374	3340	2	0
117	374	3341	2	0
51	361	3341	4	0
5	362	3323	0	0
143	344	3320	4	0
51	355	3318	0	0
51	353	3331	2	0
256	363	3325	2	0
51	353	3328	2	0
51	365	3341	4	0
259	362	3328	2	0
145	368	3321	0	0
51	374	3337	2	0
145	356	3321	0	0
51	371	3335	4	0
51	378	3340	6	0
51	378	3337	6	0
51	371	3323	2	0
51	378	3332	6	0
51	371	3329	0	0
64	367	3332	1	1
64	369	3336	0	1
63	374	3332	1	1
6	412	443	2	0
6	414	449	2	0
6	412	446	2	0
0	413	492	6	0
0	415	496	6	0
34	410	492	6	0
34	411	490	6	0
34	413	490	6	0
193	414	502	0	0
283	410	508	1	0
6	412	507	1	0
34	412	505	1	0
34	411	510	1	0
34	415	506	1	0
0	410	497	6	0
34	414	511	1	0
300	414	509	0	0
194	409	504	4	0
262	420	489	6	0
34	421	495	0	0
1	418	489	0	1
193	418	500	0	0
34	416	507	1	0
25	428	461	4	0
5	425	487	0	0
359	426	458	4	0
45	428	460	2	0
45	426	461	0	0
45	427	461	0	0
25	425	461	4	0
29	428	483	4	0
29	428	482	4	0
0	425	483	6	0
45	428	459	2	0
45	428	458	2	0
45	425	457	4	0
45	428	457	4	0
45	425	460	6	0
45	425	459	6	0
45	425	458	6	0
1	427	485	1	1
269	426	438	2	0
46	446	436	2	0
46	447	437	2	0
283	445	439	4	0
284	446	432	4	0
46	446	434	2	0
20	444	440	6	0
283	445	432	4	0
284	446	439	4	0
284	447	432	4	0
284	447	439	4	0
46	447	433	2	0
20	444	432	6	0
297	427	492	6	0
29	427	488	0	0
34	427	494	6	0
0	426	494	6	0
0	424	493	6	0
1	428	492	0	1
34	438	484	2	0
0	439	477	0	0
34	436	486	2	0
63	439	497	6	0
5	436	484	2	0
273	435	481	2	0
0	437	469	0	0
48	432	484	0	0
273	434	481	2	0
29	438	493	0	0
11	432	480	0	0
0	433	468	0	0
1	434	494	0	1
1	435	486	0	1
0	438	463	0	0
0	445	459	0	0
0	443	463	0	0
0	445	467	0	0
205	441	470	0	0
0	441	472	0	0
281	440	484	4	0
5	442	484	4	0
273	441	487	4	0
1	444	486	1	1
7	442	488	0	0
8	443	489	0	0
27	440	488	4	0
25	447	493	2	0
25	447	492	2	0
8	442	489	2	0
22	441	503	0	0
139	441	497	6	0
45	439	508	4	0
21	441	509	6	0
21	439	509	0	0
293	436	504	6	0
292	436	509	6	0
291	441	507	0	0
45	441	508	4	0
290	441	506	0	0
45	440	508	4	0
37	454	439	6	0
285	454	437	6	0
284	453	436	0	0
283	454	436	6	0
283	454	432	0	0
283	451	436	6	0
37	452	437	6	0
284	453	432	4	0
284	451	432	4	0
285	452	439	6	0
284	451	437	6	0
284	450	439	4	0
46	449	437	2	0
284	453	446	4	0
285	454	441	6	0
46	449	433	2	0
37	452	445	6	0
284	454	447	2	0
284	448	439	4	0
284	452	436	0	0
283	454	451	6	0
284	451	443	2	0
284	451	445	2	0
284	454	450	2	0
283	451	441	6	0
284	454	452	2	0
284	451	444	2	0
284	454	453	2	0
285	452	443	6	0
1	453	471	0	0
284	454	457	6	0
0	449	457	0	0
283	451	439	4	0
284	452	432	4	0
283	454	456	6	0
283	450	432	4	0
284	451	438	6	0
284	449	439	4	0
46	450	436	2	0
284	449	432	4	0
284	448	432	4	0
285	454	445	6	0
46	450	434	2	0
37	452	441	6	0
284	451	442	2	0
20	454	463	0	0
284	452	446	4	0
284	454	455	6	0
283	451	446	2	0
284	454	448	2	0
284	454	454	2	0
37	454	443	6	0
284	454	449	2	0
284	454	461	6	0
283	454	446	2	0
284	451	440	2	0
284	454	459	6	0
284	454	460	6	0
0	448	464	0	0
284	454	458	6	0
1	455	466	0	0
283	454	462	6	0
0	449	474	0	0
1	453	475	0	0
25	450	492	3	0
25	448	491	0	0
25	449	491	0	0
25	448	494	3	0
1	450	494	2	1
298	454	497	3	0
37	463	492	0	0
43	461	515	4	0
21	456	518	2	0
45	456	519	0	0
143	463	513	0	0
45	457	519	0	0
143	461	513	0	0
143	458	515	2	0
45	456	522	4	0
21	456	523	0	0
45	457	522	4	0
143	461	524	4	0
143	463	524	4	0
283	466	486	0	0
37	466	485	0	0
294	471	489	0	0
285	467	495	0	0
285	467	483	0	0
283	465	485	0	0
285	465	490	0	0
285	470	497	0	0
70	470	512	2	0
37	471	481	0	0
70	466	509	2	0
15	465	501	4	0
285	470	485	0	0
37	469	487	0	0
285	470	492	0	0
285	469	492	0	0
285	467	493	0	0
285	466	490	0	0
143	466	517	6	0
20	468	517	0	0
142	467	518	0	0
1	469	503	1	1
1	463	474	0	0
1	462	477	0	0
1	457	474	0	0
205	463	471	0	0
205	463	467	0	0
1	461	470	0	0
1	462	467	0	0
26	461	457	2	0
20	458	463	0	0
20	462	463	0	0
26	461	454	2	0
5	460	451	0	0
46	460	445	2	0
46	460	442	2	0
23	457	443	2	0
1	462	447	0	1
23	457	438	2	0
26	456	434	2	0
25	461	439	2	0
145	463	433	4	0
25	461	437	2	0
46	460	439	2	0
46	460	436	2	0
25	461	435	2	0
10	469	436	4	0
10	469	437	4	0
10	467	435	6	0
145	471	433	4	0
10	466	437	0	0
285	464	444	6	0
285	464	443	6	0
145	467	433	4	0
285	464	445	6	0
10	466	436	0	0
283	469	454	2	0
280	467	436	6	0
284	466	452	2	0
284	469	453	2	0
284	469	452	2	0
46	471	450	2	0
283	466	450	2	0
284	466	451	2	0
284	466	453	2	0
10	468	438	2	0
145	469	433	4	0
10	468	435	6	0
283	464	446	6	0
10	467	438	2	0
145	465	433	4	0
283	469	461	2	0
283	471	457	0	0
284	466	460	2	0
283	466	457	2	0
283	469	457	2	0
284	469	460	2	0
63	467	442	6	0
283	464	442	6	0
284	470	457	0	0
284	466	459	2	0
145	471	441	0	0
284	466	458	2	0
284	465	457	0	0
283	464	457	0	0
284	470	454	0	0
46	465	450	2	0
283	464	454	0	0
145	465	441	0	0
283	471	454	0	0
283	469	450	2	0
63	467	450	6	0
284	465	454	0	0
283	466	454	2	0
283	466	461	2	0
20	471	463	6	0
284	469	451	2	0
58	467	463	6	0
284	469	459	2	0
284	469	458	2	0
0	470	470	0	0
1	471	466	0	0
0	465	468	0	0
1	471	475	0	0
0	466	476	0	0
20	468	520	2	0
143	466	520	6	0
18	429	1435	0	0
24	425	1434	2	0
23	433	1424	4	0
6	436	1428	6	0
6	425	1431	0	0
15	440	1431	6	0
3	432	1424	4	0
6	442	1428	0	0
47	440	1428	0	0
37	475	498	0	0
285	475	481	0	0
294	472	484	0	0
37	477	490	0	0
285	477	483	0	0
37	472	490	0	0
285	475	484	0	0
283	478	489	0	0
70	472	515	2	0
37	473	497	0	0
285	475	492	0	0
294	475	488	0	0
285	472	481	0	0
285	477	492	0	0
70	477	525	2	0
70	475	521	2	0
70	472	521	2	0
1	478	472	0	0
205	477	475	0	0
1	475	477	0	0
1	474	473	0	0
0	478	466	0	0
205	476	469	0	0
1	475	471	0	0
1	476	467	0	0
20	475	463	6	0
26	473	457	2	0
20	479	463	6	0
284	482	448	2	0
284	482	450	2	0
284	482	454	2	0
205	485	452	0	0
283	482	462	4	0
20	483	463	0	0
283	482	455	4	0
284	482	453	2	0
313	480	481	0	0
1	486	468	0	0
313	480	487	0	0
313	483	483	0	0
313	482	483	0	0
1	484	471	0	0
284	482	457	2	0
284	482	456	2	0
284	482	460	2	0
313	487	482	0	0
284	482	461	2	0
284	482	458	2	0
284	482	459	2	0
306	486	484	0	0
284	482	452	2	0
284	482	449	2	0
5	476	451	0	0
306	486	458	0	0
26	473	454	2	0
284	482	451	2	0
313	485	481	0	0
313	484	486	0	0
70	484	476	0	0
313	485	486	0	0
313	480	484	0	0
285	472	445	6	0
46	476	442	2	0
46	476	445	2	0
284	482	445	2	0
284	482	447	2	0
285	487	440	0	0
284	482	446	2	0
285	472	444	6	0
283	472	446	6	0
285	472	443	6	0
283	472	442	6	0
284	482	442	2	0
284	482	440	2	0
284	482	443	2	0
284	482	444	2	0
283	482	441	4	0
1	474	447	0	1
46	476	439	2	0
25	475	435	2	0
46	476	436	2	0
145	473	433	4	0
43	474	437	4	0
1	492	442	0	0
306	488	448	0	0
1	489	458	0	0
205	488	452	0	0
285	493	451	0	0
41	490	466	0	0
285	488	454	0	0
63	495	463	2	0
306	491	447	2	0
406	491	455	2	0
1	492	475	0	0
284	482	436	2	0
284	482	433	2	0
1	487	436	0	0
284	482	435	2	0
283	482	434	2	0
284	482	439	2	0
284	482	437	2	0
70	491	436	0	0
284	482	438	2	0
284	482	432	2	0
284	482	425	6	0
284	482	430	6	0
0	489	425	0	0
284	482	429	6	0
283	482	426	6	0
284	482	427	6	0
284	482	431	6	0
284	482	424	6	0
284	482	428	6	0
284	482	421	6	0
283	482	420	6	0
284	482	419	6	0
284	482	422	6	0
284	482	423	6	0
284	482	418	6	0
284	482	416	6	0
284	482	417	6	0
70	496	440	0	0
283	503	454	0	0
307	500	437	0	0
0	501	443	0	0
0	503	417	0	0
286	503	444	0	0
36	502	436	0	0
29	499	448	0	0
29	502	448	0	0
0	497	416	0	0
139	499	454	6	0
63	500	454	6	0
284	482	415	6	0
284	484	412	1	0
284	483	413	1	0
283	482	414	6	0
283	485	411	1	0
0	481	411	0	0
0	494	415	0	0
58	490	408	6	0
0	503	410	0	0
285	486	404	0	0
285	487	407	6	0
0	481	405	6	0
37	486	405	0	0
285	486	406	6	0
0	481	404	6	0
285	482	404	0	0
285	483	407	6	0
285	485	405	0	0
37	484	405	0	0
37	482	405	0	0
284	492	405	6	0
284	492	401	6	0
37	487	403	0	0
23	484	404	0	0
191	495	401	6	0
37	488	407	6	0
285	482	406	0	0
284	492	406	6	0
284	492	404	6	0
284	492	400	6	0
37	488	400	0	0
284	489	404	6	0
284	489	405	6	0
285	481	403	0	0
72	499	401	0	0
0	481	401	6	0
37	483	406	0	0
37	484	407	6	0
0	481	402	6	0
285	484	406	6	0
37	486	407	6	0
37	485	406	6	0
1127	500	406	6	0
285	488	401	0	0
284	489	402	6	0
37	487	404	0	0
285	488	403	0	0
285	481	400	0	0
285	487	405	0	0
285	486	403	0	0
72	499	400	0	0
285	485	407	6	0
284	489	400	6	0
37	488	402	0	0
285	483	405	0	0
37	483	404	0	0
37	487	406	6	0
191	495	400	6	0
285	488	404	0	0
283	492	407	6	0
37	488	405	0	0
284	492	402	6	0
284	492	403	6	0
284	489	403	6	0
284	489	406	6	0
283	489	407	6	0
59	496	403	6	0
285	488	406	6	0
284	489	401	6	0
1131	486	394	4	0
1130	484	399	6	0
15	484	392	0	0
0	481	399	6	0
0	481	398	6	0
191	495	398	6	0
191	495	397	6	0
1128	493	394	6	0
284	489	397	6	0
0	481	395	6	0
72	499	397	0	0
284	489	399	6	0
1126	500	392	6	0
72	499	398	0	0
63	490	395	6	0
37	488	396	0	0
284	489	398	6	0
284	492	397	6	0
285	481	397	0	0
0	481	396	6	0
284	492	399	6	0
284	492	398	6	0
285	488	399	0	0
15	497	392	2	0
37	488	398	0	0
283	492	396	6	0
283	489	396	6	0
285	488	397	0	0
1	489	392	0	1
1	487	392	1	1
1	496	392	1	1
0	506	400	0	0
0	507	392	0	0
0	509	409	0	0
0	510	414	0	0
0	505	406	0	0
0	511	394	0	0
24	486	385	0	0
274	484	385	6	0
71	487	389	0	0
47	486	388	4	0
281	484	388	0	0
22	501	386	0	0
27	491	385	0	0
48	497	387	2	0
1138	496	385	0	0
9	488	386	0	0
15	497	388	2	0
11	497	385	2	0
22	501	387	0	0
1140	482	389	6	0
47	491	391	2	0
47	493	386	4	0
1139	501	385	0	0
1132	498	390	6	0
5	494	385	0	0
1133	498	391	0	0
0	511	388	0	0
0	509	386	0	0
1139	500	385	0	0
0	511	387	0	0
1	496	387	1	1
205	484	389	1	1
1	500	388	0	1
1	496	390	1	1
1	485	391	0	1
1	494	390	1	1
47	489	1331	0	0
51	487	1335	0	0
51	486	1329	0	0
15	495	1336	0	0
51	487	1329	0	0
51	486	1334	4	0
51	495	1335	0	0
15	484	1329	0	0
15	484	1336	0	0
15	486	1336	0	0
47	490	1333	2	0
281	492	1330	0	0
15	497	1336	0	0
51	496	1332	0	0
1129	491	1337	2	0
51	497	1335	0	0
51	496	1329	0	0
51	490	1334	0	0
15	488	1336	0	0
47	493	1331	4	0
281	489	1330	0	0
1134	484	1331	2	0
1135	484	1332	2	0
47	491	1329	6	0
51	491	1334	0	0
47	489	1329	6	0
1136	498	1331	2	0
51	485	1335	0	0
15	484	1333	0	0
15	497	1329	2	0
6	494	1329	0	0
15	493	1336	0	0
1137	498	1334	2	0
15	497	1332	2	0
1	487	1333	1	1
1	494	1335	0	1
1	496	1334	1	1
1	489	1334	0	1
1	496	1331	1	1
1	487	1330	1	1
1	492	1334	0	1
1	488	1335	0	1
0	519	390	0	0
0	514	389	0	0
0	516	396	0	0
0	519	385	0	0
0	517	399	0	0
0	517	392	0	0
0	516	385	0	0
0	513	398	0	0
0	517	405	0	0
0	515	402	0	0
0	513	386	0	0
0	518	402	0	0
0	519	392	0	0
0	513	391	0	0
0	517	403	0	0
0	516	387	0	0
36	507	439	0	0
307	509	439	0	0
36	504	439	0	0
36	505	439	0	0
307	513	438	0	0
0	525	397	0	0
0	521	395	0	0
0	524	401	0	0
0	525	405	0	0
0	522	406	0	0
0	524	388	0	0
0	520	388	0	0
0	520	404	0	0
0	523	385	0	0
0	525	385	0	0
0	521	400	0	0
0	525	392	0	0
0	520	405	0	0
0	520	408	0	0
0	523	409	0	0
0	520	409	0	0
192	524	427	0	0
192	524	425	0	0
192	524	424	0	0
192	524	426	0	0
36	522	439	0	0
36	522	436	0	0
36	521	433	0	0
151	535	439	6	0
153	536	438	2	0
153	543	436	0	0
153	541	436	0	0
153	538	436	0	0
151	536	436	2	0
153	542	436	0	0
153	540	436	0	0
20	540	438	0	0
153	539	436	0	0
34	506	446	0	0
306	511	447	0	0
306	508	444	0	0
153	519	443	4	0
34	507	441	0	0
306	515	445	0	0
205	514	446	0	0
34	517	442	0	0
205	514	447	0	0
153	525	443	0	0
283	521	446	0	0
152	522	442	4	0
153	524	443	0	0
383	519	440	4	0
150	518	443	0	0
151	529	442	6	0
153	534	440	0	0
153	533	440	0	0
0	541	440	0	0
153	532	440	0	0
151	530	440	2	0
0	541	446	0	0
153	528	443	0	0
150	518	440	0	0
153	520	440	0	0
153	521	443	0	0
151	521	440	0	0
153	526	443	0	0
383	520	443	4	0
0	539	443	0	0
153	527	443	0	0
0	537	445	0	0
307	512	441	0	0
286	514	441	0	0
1	541	447	0	0
0	537	441	0	0
101	540	445	1	1
1	519	446	0	1
286	502	459	0	0
308	499	461	0	0
283	501	457	0	0
34	503	461	0	0
272	500	467	2	0
274	497	471	0	0
272	500	465	2	0
313	489	486	0	0
306	502	487	0	0
313	489	483	0	0
47	460	1382	0	0
7	465	1380	6	0
6	460	1395	0	0
257	460	1390	0	0
5	459	1393	0	0
45	476	1381	0	0
55	463	1386	0	0
3	460	1385	0	0
24	461	1384	0	0
47	468	1382	4	0
25	476	1386	4	0
45	476	1384	4	0
45	473	1384	4	0
63	464	1384	6	0
45	472	1384	6	0
5	477	1393	0	0
25	473	1386	4	0
25	476	1388	4	0
47	463	1377	6	0
7	463	1381	2	0
47	460	1380	0	0
45	472	1381	6	0
25	463	1390	0	0
47	468	1380	4	0
45	472	1383	6	0
273	466	1380	6	0
45	473	1381	0	0
278	461	1388	0	0
45	472	1382	6	0
44	474	1381	4	0
273	462	1380	6	0
47	467	1377	6	0
3	460	1387	0	0
47	465	1377	6	0
19	474	1390	4	0
3	460	1388	0	0
25	473	1388	4	0
6	476	1395	0	0
45	474	1384	4	0
45	475	1384	4	0
117	463	1384	1	1
1	473	1392	1	1
1	464	1392	1	1
42	490	1410	0	0
6	477	2337	0	0
6	459	2337	0	0
287	461	2335	0	0
313	486	489	0	0
313	484	491	0	0
313	480	491	0	0
313	481	489	0	0
313	483	488	0	0
313	489	492	0	0
313	488	490	0	0
205	507	463	0	0
283	507	459	0	0
306	509	460	0	0
1	511	464	0	0
0	510	485	0	0
0	506	461	0	0
70	507	467	0	0
37	494	500	0	0
1	483	496	0	0
1	491	496	0	0
37	499	503	0	0
1	497	502	0	0
0	503	499	4	0
0	511	498	4	0
2	503	502	4	0
37	490	499	0	0
5	510	507	6	0
1	486	506	0	0
0	490	508	0	0
37	489	506	0	0
3	509	506	6	0
1	505	510	4	0
313	506	507	0	0
313	506	509	0	0
313	506	505	0	0
11	510	505	6	0
1	511	511	0	1
308	485	520	3	0
308	489	522	3	0
283	493	535	0	0
283	491	532	0	0
273	486	542	0	0
273	486	541	0	0
311	494	543	0	0
1	489	543	1	1
272	490	550	0	0
273	492	551	0	0
273	486	546	0	0
273	486	545	0	0
272	490	551	0	0
1	492	549	0	1
0	503	518	0	0
34	501	521	0	0
283	502	532	0	0
1	502	543	6	0
1	497	540	6	0
308	497	549	6	0
598	487	556	0	0
306	503	558	0	0
665	486	557	0	0
37	500	554	0	0
273	492	553	0	0
306	502	553	4	0
1	500	552	6	0
679	497	556	4	0
20	488	556	0	0
598	485	556	0	0
675	486	556	0	0
672	487	553	4	0
671	487	552	0	0
679	496	556	4	0
164	495	558	1	1
164	495	556	1	1
163	497	559	1	1
164	495	557	1	1
164	495	559	1	1
163	497	558	1	1
163	497	557	1	1
163	497	556	1	1
285	510	531	0	0
1	508	543	4	0
145	511	533	4	0
145	508	536	6	0
26	509	544	2	0
284	511	553	0	0
37	507	541	5	0
0	510	542	2	0
1	505	555	0	0
145	511	539	0	0
145	509	533	4	0
0	509	549	2	0
34	505	532	0	0
1	511	547	2	0
41	508	535	4	0
37	509	546	3	0
34	505	547	0	0
306	504	514	0	0
1	508	514	4	0
3	518	461	0	0
286	519	456	0	0
205	519	460	0	0
278	515	461	0	0
27	516	459	0	0
309	515	476	0	0
306	512	460	0	0
13	518	476	6	0
309	519	471	2	0
309	519	476	3	0
12	517	478	6	0
0	515	494	0	0
286	515	490	0	0
12	515	470	4	0
308	513	466	0	0
5	516	467	0	0
13	515	472	4	0
12	515	474	5	0
205	519	466	0	0
310	519	494	0	0
0	517	486	0	0
1	519	465	1	1
1	517	459	0	1
0	516	454	0	0
1	508	455	0	0
50	511	451	0	0
285	509	451	0	0
50	514	451	0	0
286	516	450	0	0
278	519	450	6	0
283	505	453	0	0
308	506	449	0	0
286	517	449	0	0
1	514	455	0	1
274	524	448	4	0
16	525	452	4	0
16	525	453	4	0
16	525	451	4	0
308	523	459	0	0
63	522	455	6	0
280	520	449	6	0
5	525	462	0	0
278	520	448	6	0
12	526	471	6	0
205	526	464	0	0
12	526	469	6	0
280	523	465	0	0
63	523	471	2	0
278	523	464	0	0
13	525	469	7	0
1	522	465	1	1
19	524	477	4	0
23	524	473	4	0
25	522	473	6	0
25	522	475	6	0
12	520	474	6	0
25	522	477	6	0
23	524	475	4	0
0	537	449	0	0
0	542	458	0	0
0	543	461	0	0
0	537	467	0	0
0	543	465	0	0
0	540	457	0	0
0	542	453	0	0
308	520	487	0	0
307	524	484	0	0
310	521	492	0	0
310	524	489	0	0
0	526	492	0	0
286	524	495	0	0
286	522	494	0	0
308	538	488	4	0
1	542	495	0	0
20	543	473	0	0
55	522	1406	4	0
47	522	1408	0	0
121	522	1411	4	0
55	522	1407	4	0
15	524	1410	0	0
6	525	1406	0	0
6	516	1411	0	0
15	516	1408	4	0
286	515	497	0	0
286	519	497	0	0
34	513	499	0	0
0	526	502	0	0
0	520	501	0	0
0	521	497	0	0
34	524	500	0	0
286	525	497	0	0
205	548	468	0	0
0	545	469	0	0
37	548	485	4	0
0	551	473	0	0
285	549	487	4	0
310	548	484	4	0
1	544	502	0	0
37	551	488	4	0
0	549	475	0	0
205	547	460	0	0
0	546	457	0	0
1	547	461	0	0
1	545	463	0	0
0	545	462	0	0
1	547	452	0	0
0	549	455	0	0
1	550	445	0	0
1	550	443	0	0
1	547	445	0	0
1	548	447	0	0
0	551	440	0	0
0	545	441	0	0
0	544	447	0	0
205	547	442	0	0
153	549	436	0	0
153	550	436	0	0
153	548	436	0	0
153	547	436	0	0
153	551	436	0	0
153	546	436	0	0
153	545	436	0	0
153	544	436	0	0
153	558	436	0	0
0	558	438	0	0
153	559	436	0	0
153	557	436	0	0
0	553	454	0	0
153	556	436	0	0
0	556	444	0	0
22	557	453	0	0
0	556	445	0	0
1	558	442	0	0
153	552	436	0	0
0	552	462	0	0
1	559	461	0	0
6	554	451	0	0
1	554	456	0	0
1	553	462	0	0
205	555	442	0	0
205	556	442	0	0
153	555	436	0	0
153	554	436	0	0
0	555	437	0	0
153	553	436	0	0
1	555	446	0	0
1	559	450	0	0
205	557	456	0	0
0	556	463	0	0
205	556	460	0	0
0	555	457	0	0
110	554	449	0	1
0	558	465	0	0
0	556	470	0	0
0	557	465	0	0
0	553	466	0	0
205	559	468	0	0
0	558	468	0	0
1	552	468	0	0
0	553	475	0	0
4	555	485	0	0
3	558	483	0	0
285	552	490	4	0
37	554	491	4	0
4	556	488	0	0
1	553	500	0	0
0	567	463	0	0
205	567	457	0	0
1	566	462	0	0
205	565	461	0	0
0	562	461	0	0
0	561	458	0	0
1	561	462	0	0
4	561	489	0	0
4	562	490	0	0
0	566	474	0	0
358	564	492	4	0
15	560	483	4	0
356	560	472	2	0
82	567	484	1	1
1	560	487	0	1
1	568	461	0	0
0	569	461	0	0
0	570	463	0	0
0	574	469	0	0
1	569	457	0	0
350	569	488	4	0
0	570	468	0	0
1	569	464	0	0
306	568	495	4	0
0	572	475	0	0
0	568	471	0	0
1	571	468	0	0
205	572	466	0	0
1	571	469	0	0
352	570	489	4	0
205	571	465	0	0
350	570	488	4	0
351	571	495	4	0
350	571	488	4	0
5	583	459	2	0
0	577	468	0	0
283	578	466	0	0
1	577	464	0	0
7	583	463	6	0
1	578	460	0	0
283	583	457	0	0
1	579	468	0	0
0	577	462	0	0
55	587	459	2	0
55	587	460	2	0
1	589	471	0	0
3	584	463	2	0
21	588	495	0	0
283	590	470	0	0
45	587	495	0	0
0	585	474	0	0
283	586	473	0	0
21	585	495	6	0
7	585	463	2	0
283	587	472	0	0
283	588	470	0	0
0	589	470	0	0
1	588	471	0	0
45	586	495	0	0
1	588	462	1	1
1	585	459	0	1
286	599	457	0	0
407	598	461	4	0
681	597	458	2	0
681	596	458	2	0
530	595	458	2	0
680	594	458	2	0
290	599	468	0	0
680	593	458	2	0
36	599	474	2	0
290	599	467	0	0
36	598	472	0	0
21	598	471	0	0
193	599	471	3	1
193	597	471	0	1
193	598	471	0	1
0	569	450	0	0
0	563	449	0	0
355	568	454	2	0
1	561	452	0	0
0	564	455	0	0
1	570	448	0	0
1	566	451	0	0
153	590	450	0	0
283	578	452	0	0
153	591	450	0	0
0	577	451	0	0
153	589	450	0	0
45	599	451	4	0
153	584	448	2	0
355	567	450	6	0
1	566	455	0	0
0	571	452	0	0
355	568	453	2	0
45	595	449	0	0
45	599	449	0	0
45	596	449	0	0
153	597	450	0	0
45	593	451	4	0
383	599	450	0	0
45	597	449	0	0
153	596	450	0	0
45	593	449	0	0
45	596	451	4	0
45	595	451	4	0
153	595	450	0	0
355	567	451	6	0
219	567	452	6	0
355	567	449	6	0
153	592	450	0	0
45	592	451	6	0
151	584	449	4	0
1	571	453	0	0
1	578	448	0	0
149	586	450	0	0
149	588	450	0	0
220	568	452	4	0
149	587	450	0	0
283	577	455	0	0
355	569	452	4	0
45	597	451	4	0
45	598	451	4	0
153	598	450	0	0
45	598	449	0	0
149	593	450	0	0
45	594	451	4	0
153	594	450	0	0
45	594	449	0	0
45	592	449	6	0
151	606	454	4	0
153	606	452	6	0
383	606	453	2	0
151	605	450	0	0
153	604	450	0	0
45	600	451	2	0
45	600	449	2	0
153	601	450	0	0
153	600	450	0	0
22	603	470	5	0
407	603	464	2	0
290	600	464	0	0
21	606	474	0	0
407	603	471	6	0
290	603	466	5	0
153	603	450	0	0
0	607	478	0	0
290	601	465	0	0
36	606	485	4	0
36	606	483	4	0
21	601	473	0	0
36	600	475	3	0
286	603	449	0	0
407	602	462	0	0
153	602	450	0	0
22	604	470	0	0
290	605	465	0	0
994	601	468	0	0
23	600	462	0	0
22	602	465	5	0
290	602	464	0	0
22	602	470	5	0
36	601	474	4	0
36	606	481	4	0
36	601	476	4	0
58	604	474	6	0
36	605	484	4	0
0	605	476	0	0
36	606	482	4	0
194	605	468	1	1
193	606	474	0	1
181	603	473	3	1
193	600	472	3	1
193	601	473	0	1
193	607	474	3	1
193	602	473	0	1
0	573	445	0	0
1	569	443	0	0
1	573	440	0	0
0	583	447	0	0
205	568	446	0	0
153	584	441	2	0
153	584	447	2	0
283	588	444	0	0
283	576	443	0	0
283	604	442	0	0
153	584	442	2	0
283	601	446	0	0
1	578	444	0	0
1	582	446	0	0
153	584	445	2	0
153	584	443	2	0
1	578	445	0	0
153	584	446	2	0
153	584	444	2	0
1	580	435	0	0
153	570	436	0	0
0	570	438	0	0
153	569	436	0	0
153	571	436	0	0
153	581	439	0	0
153	573	436	0	0
0	582	437	0	0
151	583	439	0	0
153	575	436	0	0
153	574	436	0	0
0	573	437	0	0
153	577	436	4	0
153	582	439	0	0
153	568	436	0	0
153	572	436	0	0
283	578	438	0	0
153	576	436	4	0
205	568	435	0	0
151	579	438	4	0
151	578	436	0	0
153	561	436	0	0
1	562	447	0	0
1	565	441	0	0
153	564	436	0	0
0	564	447	0	0
0	564	441	0	0
153	563	436	0	0
153	562	436	0	0
153	567	436	0	0
0	567	438	0	0
153	565	436	0	0
1	565	443	0	0
1	562	444	0	0
0	560	445	0	0
153	560	436	0	0
1	565	447	0	0
153	566	436	0	0
1	567	440	0	0
21	329	715	0	0
45	331	714	4	0
21	325	715	0	0
45	326	714	4	0
21	322	714	4	0
36	339	702	0	0
45	327	714	4	0
45	332	714	4	0
34	340	714	0	0
32	337	715	0	0
24	339	704	0	0
45	330	714	4	0
34	343	717	0	0
1	337	717	0	0
3	341	705	4	0
36	338	718	0	0
34	340	716	0	0
34	343	713	0	0
32	338	717	0	0
36	337	702	0	0
45	324	714	4	0
45	323	714	4	0
21	322	713	6	0
45	333	714	4	0
45	328	714	4	0
45	325	714	4	0
45	329	714	4	0
36	341	717	0	0
1	338	710	0	0
182	337	711	0	0
22	341	707	4	0
33	337	714	0	0
15	337	704	4	0
1	339	708	0	1
183	351	701	0	0
36	351	698	0	0
183	351	699	0	0
183	348	703	0	0
183	351	703	0	0
183	348	701	0	0
36	346	701	0	0
36	348	698	0	0
3	349	713	2	0
36	347	700	0	0
183	344	703	0	0
183	348	707	0	0
183	344	707	0	0
36	346	718	0	0
34	346	716	0	0
183	344	705	0	0
7	349	712	4	0
34	347	718	0	0
34	350	713	0	0
3	349	714	2	0
183	351	705	0	0
183	351	707	0	0
7	349	715	0	0
183	348	705	0	0
1	350	718	0	0
36	344	718	0	0
16	344	714	2	0
7	344	713	6	0
36	345	701	0	0
29	345	714	0	0
1	347	712	0	1
183	354	703	0	0
183	357	699	0	0
183	357	697	0	0
183	357	701	0	0
183	357	703	0	0
34	356	715	0	0
29	359	714	0	0
183	354	701	0	0
36	353	717	0	0
36	354	696	0	0
183	354	699	0	0
183	357	707	0	0
34	355	711	0	0
183	357	705	0	0
36	353	697	0	0
36	352	717	0	0
36	354	716	0	0
34	352	714	0	0
36	353	716	0	0
183	354	705	0	0
183	354	707	0	0
34	352	711	0	0
36	355	717	0	0
183	363	695	0	0
183	363	701	0	0
183	363	697	0	0
183	363	703	0	0
183	360	703	0	0
183	360	707	0	0
33	366	716	0	0
89	361	711	2	0
183	363	699	0	0
183	360	697	0	0
183	363	705	0	0
183	360	701	0	0
183	360	705	0	0
183	360	699	0	0
29	361	714	0	0
1	361	712	0	1
32	375	714	0	0
36	375	703	0	0
32	374	706	0	0
32	374	701	0	0
34	373	696	0	0
33	375	716	0	0
36	370	718	0	0
34	373	711	0	0
32	369	692	0	0
32	374	692	0	0
34	374	698	0	0
21	373	689	2	0
21	368	689	4	0
5	147	3330	2	0
33	373	709	0	0
34	371	709	0	0
32	369	711	2	0
36	383	692	0	0
32	381	693	0	0
36	381	718	0	0
36	378	691	0	0
34	376	697	0	0
32	382	716	0	0
34	381	697	0	0
34	380	715	0	0
33	377	695	0	0
34	378	699	0	0
36	376	702	0	0
33	381	712	0	0
34	382	708	0	0
34	379	713	0	0
36	376	717	0	0
32	377	707	0	0
32	381	705	0	0
36	378	702	0	0
33	379	709	0	0
32	381	700	0	0
34	377	692	0	0
36	381	692	0	0
193	373	680	0	0
21	368	685	4	0
193	368	687	0	0
194	368	684	4	0
21	368	681	4	0
21	373	685	2	0
193	373	684	0	0
193	368	680	0	0
21	373	681	2	0
33	389	695	0	0
33	390	687	0	0
34	390	702	0	0
194	370	678	4	0
34	388	696	0	0
33	389	699	0	0
33	386	703	0	0
34	389	690	0	0
32	387	706	0	0
21	369	678	2	0
21	372	678	4	0
34	387	710	0	0
33	390	710	0	0
34	390	705	0	0
32	386	697	0	0
34	394	697	0	0
33	399	676	0	0
34	396	677	0	0
33	398	707	0	0
34	396	682	0	0
33	396	685	0	0
33	396	702	0	0
33	394	688	0	0
33	393	679	0	0
34	392	711	0	0
34	392	700	0	0
33	392	692	0	0
34	392	683	0	0
33	392	706	0	0
34	404	704	0	0
34	401	702	0	0
34	406	708	0	0
34	402	708	0	0
34	401	680	0	0
34	400	681	0	0
34	414	710	0	0
33	410	712	0	0
33	404	675	0	0
33	412	675	0	0
32	411	674	0	0
32	422	674	0	0
34	423	692	0	0
6	422	695	2	0
34	422	698	0	0
34	423	704	0	0
34	419	708	0	0
34	422	708	0	0
32	431	696	0	0
34	427	701	0	0
32	424	694	0	0
34	428	699	0	0
34	428	695	0	0
34	426	710	0	0
34	424	700	0	0
33	424	697	0	0
34	427	708	0	0
55	439	697	2	0
11	438	701	0	0
3	436	691	0	0
27	436	694	2	0
27	439	691	0	0
76	439	694	0	1
47	436	697	2	0
268	444	690	0	0
33	447	694	4	0
33	447	701	4	0
32	447	697	4	0
32	444	697	4	0
33	447	704	4	0
71	441	702	2	0
395	445	726	0	0
5	443	702	0	0
1	444	693	1	1
1	444	700	1	1
25	447	686	2	0
25	447	682	2	0
22	442	680	4	0
25	447	684	2	0
254	434	682	4	0
33	426	674	0	0
32	430	679	0	0
32	429	674	0	0
30	442	673	0	0
278	439	678	4	0
11	447	678	6	0
3	439	677	0	0
15	441	677	4	0
1	440	677	0	1
32	438	668	0	0
33	438	666	0	0
34	443	668	0	0
32	447	664	0	0
33	445	664	0	0
33	436	669	0	0
34	434	669	0	0
143	407	3513	6	0
112	424	3524	0	0
112	425	3526	0	0
112	423	3523	0	0
112	421	3523	0	0
5	422	3527	2	0
58	405	3518	0	1
58	406	3518	0	1
143	403	3505	2	0
143	407	3510	6	0
143	407	3507	6	0
51	422	3482	6	0
207	420	3480	1	0
214	422	3484	0	0
51	419	3485	4	0
51	415	3485	4	0
51	409	3482	2	0
207	418	3484	3	0
207	416	3484	6	0
59	414	3480	1	1
51	417	3475	0	0
51	422	3477	6	0
51	419	3475	0	0
207	418	3476	5	0
51	409	3475	2	0
215	422	3479	0	0
207	416	3476	3	0
51	417	3472	4	0
51	414	3464	2	0
51	413	3467	2	0
51	421	3467	6	0
51	423	3462	6	0
41	419	3460	4	0
113	408	617	0	0
234	405	641	0	0
1	411	636	0	0
45	421	631	0	0
233	410	641	0	0
164	411	643	0	0
32	415	619	0	0
33	409	609	0	0
32	407	633	0	0
164	411	641	0	0
34	422	630	4	0
33	409	608	0	0
45	421	630	6	0
32	405	635	0	0
33	422	629	4	0
164	404	640	0	0
32	416	626	4	0
32	417	629	4	0
45	419	631	0	0
34	416	628	4	0
32	420	626	4	0
32	417	630	4	0
45	418	631	0	0
42	419	628	4	0
107	428	634	0	0
113	407	618	0	0
32	418	632	4	0
33	405	613	0	0
33	423	640	0	0
33	420	632	4	0
107	429	634	0	0
33	428	641	0	0
111	418	636	0	0
164	417	644	0	0
164	403	641	0	0
164	415	641	0	0
1	422	616	0	0
45	421	629	6	0
111	416	611	0	0
33	407	639	0	0
107	415	610	0	0
33	411	608	0	0
1	414	637	0	0
1	418	623	0	0
1	416	617	0	0
45	421	628	6	0
164	412	642	0	0
34	412	639	0	0
111	415	612	0	0
45	418	628	2	0
109	418	638	0	0
45	420	631	0	0
45	418	630	2	0
45	418	629	2	0
45	418	627	4	0
45	421	627	4	0
34	422	644	0	0
34	420	634	4	0
109	418	637	0	0
32	425	645	0	0
34	427	644	0	0
1	424	613	0	0
32	425	618	0	0
107	430	635	0	0
32	414	607	0	0
34	455	670	0	0
3	452	679	6	0
3	451	679	6	0
3	452	678	6	0
3	451	678	6	0
34	450	667	0	0
3	450	679	6	0
7	454	685	0	0
3	450	678	6	0
41	448	673	4	0
274	451	672	4	0
11	454	681	2	0
25	453	687	0	0
3	453	683	0	0
3	453	684	0	0
7	455	683	2	0
34	452	692	6	0
7	450	684	4	0
7	454	682	4	0
3	450	682	6	0
7	451	682	2	0
48	451	681	2	0
25	455	687	0	0
3	454	684	0	0
32	450	693	6	0
3	450	685	0	0
34	449	693	6	0
51	454	702	6	0
34	450	691	6	0
3	454	683	0	0
33	451	693	2	0
267	451	688	6	0
34	451	694	6	0
7	450	686	0	0
33	451	692	4	0
278	454	701	2	0
33	450	692	0	0
278	451	698	0	0
278	452	697	0	0
278	452	703	0	0
51	452	696	0	0
278	454	700	3	0
51	449	702	2	0
278	450	703	0	0
273	453	700	6	0
288	448	700	4	0
278	449	697	0	0
278	451	696	0	0
273	450	697	0	0
1	454	674	1	1
2	449	699	1	1
1	451	688	0	1
82	448	677	1	1
78	448	682	0	1
33	455	663	0	0
16	450	704	2	0
16	453	704	2	0
194	453	710	2	0
51	452	706	4	0
29	451	704	0	0
32	452	725	4	0
396	449	727	1	0
33	455	723	4	0
283	463	686	4	0
283	458	690	4	0
33	463	689	4	0
283	458	684	4	0
430	462	719	0	0
430	463	714	0	0
430	460	717	0	0
34	459	717	0	0
430	463	717	0	0
430	460	715	0	0
75	463	681	0	1
404	463	721	4	0
396	463	727	4	0
33	463	722	1	0
400	462	722	1	0
403	463	724	2	0
404	460	725	4	0
395	459	722	0	0
405	458	727	4	0
401	461	721	0	0
33	457	722	4	0
405	461	723	7	0
32	460	721	1	0
33	456	725	1	0
397	461	727	1	0
15	459	678	6	0
263	461	675	2	0
1	462	679	1	1
1	459	675	1	1
77	463	676	0	1
79	456	679	1	1
80	459	674	0	1
15	461	668	2	0
15	461	665	2	0
1	457	667	1	1
45	466	654	2	0
45	466	652	2	0
21	465	655	6	0
45	466	651	2	0
45	466	656	2	0
21	465	650	6	0
34	466	665	0	0
34	468	687	4	0
48	466	678	4	0
32	469	661	0	0
15	468	678	4	0
5	468	675	0	0
11	467	677	2	0
265	469	672	2	0
45	466	655	2	0
45	466	649	2	0
32	464	661	0	0
45	466	648	2	0
320	468	651	6	0
37	471	684	4	0
15	465	672	4	0
45	466	650	2	0
33	464	658	0	0
37	469	684	4	0
34	469	668	0	0
32	465	669	0	0
34	470	687	4	0
1	470	666	0	0
45	466	653	2	0
1	470	676	1	1
2	467	676	0	1
1	467	680	1	1
37	475	684	4	0
33	474	687	4	0
34	479	668	0	0
34	477	662	0	0
1	477	670	0	0
33	476	668	0	0
34	477	687	4	0
43	474	673	4	0
32	472	669	0	0
37	473	684	4	0
33	473	660	0	0
34	473	665	0	0
142	474	681	6	0
32	476	664	0	0
81	472	674	1	1
285	487	649	1	0
36	484	657	0	0
285	486	653	1	0
395	484	648	0	0
36	487	659	0	0
33	483	677	6	0
394	485	686	6	0
395	484	674	1	0
394	482	684	3	0
395	487	677	1	0
33	485	668	1	0
33	484	682	6	0
395	482	687	3	0
395	484	685	6	0
395	485	683	3	0
321	468	646	6	0
45	466	644	2	0
45	467	644	0	0
21	468	644	0	0
22	467	645	0	0
45	466	645	2	0
45	466	646	2	0
21	465	644	6	0
45	466	647	2	0
399	484	644	1	0
394	493	651	1	0
395	491	654	6	0
285	493	653	1	0
36	494	668	0	0
36	494	667	0	0
32	491	669	1	0
286	494	655	1	0
113	495	660	0	0
36	495	664	0	0
113	495	657	0	0
36	495	666	0	0
286	488	651	1	0
285	488	653	1	0
286	495	673	1	0
33	490	677	6	0
395	493	674	1	0
395	488	647	0	0
286	491	649	1	0
113	494	660	0	0
113	493	659	0	0
395	489	650	6	0
113	493	658	0	0
394	489	654	0	0
286	492	656	1	0
25	490	651	1	0
395	488	675	1	0
36	495	685	3	0
36	495	687	6	0
36	492	687	6	0
36	495	684	6	0
36	494	684	1	0
405	489	680	1	0
395	489	685	1	0
394	490	687	1	0
36	491	687	6	0
395	500	662	1	0
33	500	656	1	0
395	503	661	1	0
33	496	683	6	0
285	500	674	1	0
394	499	665	1	0
396	502	674	1	0
286	500	675	1	0
394	500	682	3	0
32	497	660	1	0
286	497	675	1	0
32	497	657	1	0
394	499	668	1	0
113	496	658	0	0
33	499	678	6	0
113	496	659	0	0
285	503	677	1	0
285	497	673	1	0
32	496	677	6	0
33	465	688	4	0
33	471	688	4	0
32	478	694	2	0
33	479	690	2	0
395	476	695	2	0
394	476	693	1	0
405	494	688	3	0
395	484	689	6	0
394	489	692	3	0
377	483	692	0	0
377	488	690	1	0
285	482	695	2	0
396	488	695	5	0
394	498	694	7	0
399	489	695	5	0
32	496	690	6	0
32	499	688	6	0
377	498	690	7	0
396	487	693	5	0
286	484	694	2	0
394	475	690	1	0
402	487	690	3	0
395	492	693	3	0
396	492	695	5	0
406	490	692	6	0
377	500	690	7	0
32	502	695	0	0
394	500	693	7	0
377	493	690	7	0
402	500	691	3	0
395	486	700	1	0
33	486	698	0	0
395	475	702	1	0
395	484	702	5	0
33	479	696	2	0
394	499	702	6	0
394	498	701	6	0
183	473	699	2	0
394	500	703	3	0
394	501	701	3	0
405	491	700	6	0
32	473	703	2	0
394	481	696	2	0
32	495	703	6	0
394	483	698	2	0
32	493	703	6	0
183	474	700	2	0
401	479	699	6	0
404	497	699	0	0
33	486	702	0	0
34	477	701	2	0
401	477	703	6	0
394	501	702	6	0
205	476	696	2	0
399	481	701	2	0
32	477	699	2	0
34	502	702	6	0
34	491	699	6	0
34	498	700	6	0
394	482	699	2	0
394	476	698	1	0
395	484	700	5	0
394	482	702	5	0
395	488	698	1	0
33	474	698	3	0
394	471	706	3	0
400	471	711	4	0
398	474	707	6	0
183	470	705	2	0
405	469	711	7	0
404	477	705	2	0
404	478	704	2	0
394	472	709	4	0
399	481	704	2	0
394	484	709	0	0
113	477	706	0	0
395	476	708	1	0
398	473	705	6	0
113	487	707	0	0
394	484	710	1	0
394	482	706	4	0
32	494	704	6	0
394	495	711	3	0
394	483	708	2	0
395	485	706	5	0
214	473	708	2	0
113	474	705	0	0
32	487	704	0	0
398	477	710	6	0
113	475	705	0	0
113	486	706	0	0
394	485	709	1	0
394	502	704	5	0
97	501	705	0	0
97	500	706	0	0
402	488	706	0	0
394	502	708	0	0
403	501	710	1	0
32	495	705	6	0
97	500	704	0	0
429	500	705	0	0
113	489	706	0	0
404	481	710	2	0
395	495	708	0	0
399	484	704	2	0
401	479	707	6	0
113	476	706	0	0
394	490	704	0	0
32	490	709	0	0
404	488	709	2	0
402	499	707	0	0
34	496	708	6	0
97	499	705	0	0
113	474	704	0	0
34	495	704	6	0
394	493	709	0	0
45	472	1618	6	0
45	476	1620	4	0
45	472	1619	6	0
45	474	1620	4	0
45	472	1617	6	0
45	473	1617	0	0
45	473	1620	4	0
45	476	1617	0	0
45	472	1620	6	0
45	475	1620	4	0
42	474	1617	4	0
145	440	1646	0	0
15	438	1642	6	0
145	441	1642	4	0
3	442	1646	0	0
6	443	1646	0	0
395	510	674	6	0
395	508	677	2	0
394	510	676	0	0
402	505	702	0	0
286	509	673	1	0
402	505	706	0	0
32	506	698	0	0
33	504	695	0	0
32	504	710	0	0
32	508	701	0	0
394	507	704	0	0
403	507	707	3	0
394	507	680	0	0
395	504	680	4	0
396	504	675	1	0
286	506	675	1	0
395	504	698	0	0
33	510	704	0	0
396	506	673	1	0
395	504	705	0	0
33	507	710	0	0
32	508	705	0	0
396	516	666	2	0
402	508	671	2	0
695	510	668	2	0
394	515	669	1	0
694	509	670	1	0
401	508	668	2	0
402	513	670	2	0
394	513	666	1	0
15	519	665	2	0
401	510	671	2	0
402	517	664	2	0
395	519	669	1	0
396	517	668	2	0
402	507	670	2	0
394	504	665	1	0
395	521	668	1	0
402	522	669	2	0
394	522	667	1	0
394	517	662	1	0
3	519	663	2	0
402	518	659	2	0
395	519	661	1	0
394	515	663	1	0
396	521	661	1	0
395	521	659	1	0
1	520	663	0	1
33	508	644	1	0
33	507	638	1	0
33	510	637	1	0
32	509	636	1	0
455	511	634	5	0
400	470	713	6	0
400	468	714	5	0
32	499	714	0	0
33	499	716	0	0
430	466	719	0	0
394	483	713	5	0
402	498	712	0	0
34	465	714	0	0
397	494	716	0	0
400	468	712	4	0
33	495	715	0	0
34	465	719	0	0
397	502	716	0	0
34	467	718	0	0
33	490	713	0	0
32	493	713	0	0
395	500	712	0	0
400	467	716	6	0
394	486	712	5	0
33	504	714	0	0
33	466	721	3	0
396	464	725	2	0
396	470	727	4	0
32	468	723	6	0
34	465	721	1	0
394	465	723	5	0
403	466	725	4	0
33	467	724	0	0
32	470	723	6	0
401	469	725	7	0
32	467	727	6	0
403	465	720	4	0
404	466	720	4	0
33	472	725	6	0
395	464	722	0	0
32	475	731	0	0
32	472	730	6	0
405	473	728	1	0
394	475	734	0	0
33	477	733	0	0
33	474	733	1	0
399	453	735	1	0
398	454	732	6	0
397	455	734	6	0
33	455	733	6	0
397	452	730	1	0
33	451	732	1	0
404	462	733	7	0
399	458	735	2	0
396	453	728	1	0
396	467	732	1	0
405	451	729	2	0
404	459	729	2	0
399	449	733	1	0
396	458	732	0	0
404	452	729	6	0
404	470	729	7	0
32	456	730	1	0
32	471	733	6	0
405	460	735	4	0
32	463	731	1	0
395	448	731	1	0
399	465	729	6	0
398	463	729	7	0
398	461	731	3	0
399	464	735	1	0
405	465	731	7	0
398	469	733	1	0
397	468	735	1	0
403	468	730	2	0
33	441	733	1	0
33	446	728	4	0
397	446	732	2	0
404	446	734	1	0
398	444	734	1	0
394	444	731	0	0
32	447	735	1	0
397	446	730	0	0
397	455	737	4	0
405	455	742	0	0
21	454	743	0	0
396	452	738	4	0
32	454	740	4	0
404	451	741	4	0
395	468	741	1	0
21	451	743	7	0
33	471	742	6	0
394	471	737	6	0
32	460	742	6	0
397	461	737	4	0
403	463	737	4	0
398	462	739	4	0
396	465	738	4	0
395	474	736	1	0
399	465	740	4	0
33	475	743	6	0
395	470	740	0	0
398	450	736	3	0
399	450	739	3	0
33	460	740	4	0
396	463	741	4	0
399	448	742	2	0
397	456	740	3	0
396	458	739	4	0
405	468	737	6	0
398	458	742	1	0
32	467	742	4	0
396	468	739	6	0
404	473	743	6	0
397	465	742	4	0
32	473	739	6	0
397	467	740	6	0
396	444	743	2	0
401	444	738	2	0
405	446	739	3	0
401	441	737	2	0
397	442	741	2	0
399	473	749	5	0
32	474	750	3	0
32	473	746	6	0
32	475	748	2	0
397	467	746	6	0
398	471	750	5	0
401	465	750	0	0
398	466	744	6	0
395	449	751	2	0
399	456	744	1	0
21	451	745	7	0
21	451	747	7	0
401	458	744	2	0
402	466	751	0	0
394	449	750	2	0
399	471	745	6	0
21	454	747	7	0
399	449	748	3	0
399	464	745	4	0
402	464	749	0	0
21	458	751	0	0
396	467	750	5	0
402	460	747	2	0
405	476	745	1	0
404	458	747	2	0
394	448	746	2	0
395	465	748	1	0
400	468	748	5	0
395	470	749	6	0
21	459	751	0	0
405	463	747	1	0
398	461	744	4	0
50	453	744	0	0
396	469	744	6	0
407	454	750	2	0
21	454	745	7	0
401	447	745	7	0
398	446	747	2	0
700	447	751	6	0
405	445	750	1	0
33	445	751	0	0
402	444	747	7	0
395	445	747	2	0
401	442	745	2	0
32	445	745	2	0
32	441	750	2	0
395	441	751	2	0
21	440	750	2	0
21	440	749	0	0
394	436	731	1	0
33	439	730	2	0
403	437	735	1	0
32	437	733	2	0
395	434	734	7	0
402	435	736	7	0
21	437	747	0	0
32	438	741	0	0
403	438	744	0	0
32	435	739	0	0
21	438	747	0	0
402	432	744	0	0
21	433	750	0	0
397	434	732	1	0
21	436	747	0	0
32	432	749	7	0
394	436	741	0	0
21	433	749	5	0
398	439	742	2	0
396	435	743	1	0
397	433	741	1	0
401	438	738	1	0
398	432	737	1	0
21	435	747	5	0
21	439	748	5	0
21	434	748	5	0
401	435	745	0	0
21	433	751	0	0
397	430	749	0	0
6	426	740	4	0
401	431	742	1	0
402	430	743	1	0
403	430	747	0	0
402	426	750	4	0
397	425	747	4	0
398	427	748	4	0
403	424	749	4	0
394	429	751	4	0
32	427	745	1	0
396	429	748	3	0
33	424	744	1	0
404	429	744	1	0
401	422	751	4	0
394	418	749	4	0
33	419	747	4	0
32	421	745	4	0
403	418	751	1	0
396	422	747	4	0
404	421	749	4	0
405	423	752	4	0
396	421	758	1	0
397	421	756	1	0
398	422	754	1	0
33	419	756	1	0
399	419	753	1	0
401	417	757	1	0
399	427	757	0	0
395	426	753	4	0
399	418	759	1	0
21	439	754	5	0
395	428	754	1	0
32	429	758	1	0
404	417	753	1	0
396	431	752	0	0
404	431	757	0	0
21	438	755	2	0
401	433	756	1	0
399	430	755	0	0
21	437	755	2	0
404	438	756	5	0
21	434	754	5	0
21	436	755	2	0
403	438	759	3	0
404	424	757	3	0
394	439	758	6	0
398	429	753	0	0
395	448	753	2	0
21	433	752	0	0
33	449	752	0	0
395	445	752	2	0
396	448	757	2	0
399	442	756	2	0
32	443	758	7	0
399	452	758	0	0
21	454	756	0	0
403	435	757	7	0
404	434	759	3	0
21	454	755	0	0
403	452	752	2	0
21	440	753	0	0
32	450	756	2	0
21	435	755	5	0
33	449	759	2	0
394	446	753	2	0
398	445	757	2	0
394	432	758	3	0
32	447	753	0	0
21	440	752	2	0
21	433	753	5	0
33	441	752	2	0
430	441	759	7	0
32	451	756	2	0
404	450	755	2	0
397	450	759	2	0
21	457	759	0	0
21	458	759	0	0
21	462	755	0	0
21	460	757	5	0
399	463	759	4	0
21	462	754	0	0
5	458	755	0	0
21	456	753	5	0
397	470	756	4	0
397	469	754	5	0
32	470	752	7	0
397	469	752	5	0
394	470	758	7	0
399	466	755	5	0
394	465	758	4	0
403	465	753	7	0
401	466	752	0	0
401	468	753	5	0
396	468	756	4	0
397	468	759	4	0
402	467	754	5	0
398	467	753	5	0
403	466	754	1	0
32	475	759	5	0
33	474	753	0	0
405	473	755	7	0
32	475	755	6	0
33	476	757	4	0
399	473	759	4	0
398	472	757	4	0
396	472	753	5	0
394	462	766	4	0
398	461	764	5	0
396	471	760	4	0
404	470	766	5	0
396	465	761	5	0
403	468	764	5	0
396	456	762	3	0
33	475	766	5	0
404	464	763	5	0
32	476	765	5	0
394	475	764	5	0
404	472	763	5	0
397	462	761	5	0
395	470	763	4	0
399	460	762	5	0
403	469	761	6	0
32	468	766	2	0
403	465	766	5	0
32	475	763	7	0
33	474	761	6	0
398	466	760	4	0
394	467	762	6	0
33	466	764	6	0
394	472	765	4	0
397	455	760	5	0
33	452	767	2	0
5	454	764	4	0
32	451	761	2	0
398	453	760	7	0
5	449	762	4	0
33	450	762	2	0
1	454	762	3	1
5	444	763	4	0
32	444	767	2	0
32	443	760	2	0
1	447	760	0	1
1	441	761	2	1
395	436	760	3	0
405	434	767	2	0
398	434	764	7	0
404	433	765	2	0
33	437	764	2	0
403	432	763	2	0
396	433	762	7	0
399	432	766	7	0
397	433	760	7	0
394	429	761	1	0
32	427	767	0	0
33	424	764	1	0
396	426	762	1	0
405	426	765	1	0
397	428	763	1	0
32	429	765	1	0
399	424	760	1	0
622	424	762	0	0
398	427	760	1	0
404	420	760	1	0
32	419	766	0	0
397	417	764	1	0
398	417	762	1	0
33	423	766	2	0
395	422	762	1	0
403	416	766	1	0
401	419	763	1	0
402	421	765	1	0
21	410	742	0	0
21	410	745	0	0
402	411	747	0	0
395	415	750	4	0
21	409	738	4	0
402	415	756	1	0
404	414	759	1	0
396	415	754	1	0
55	415	752	7	0
396	408	762	0	0
8	409	757	5	0
395	414	765	1	0
3	408	752	0	0
3	409	755	0	0
394	408	760	1	0
22	408	744	5	0
155	412	753	2	0
394	408	765	0	0
8	410	756	7	0
396	410	765	0	0
8	410	758	2	0
71	410	753	4	0
405	411	763	1	0
32	412	767	0	0
22	401	740	5	0
21	407	741	5	0
8	406	747	5	0
22	401	741	5	0
21	404	741	0	0
155	403	733	2	0
55	406	750	5	0
21	401	738	4	0
55	407	758	7	0
8	403	740	2	0
21	400	746	4	0
22	407	756	5	0
32	400	755	0	0
22	407	754	5	0
21	400	744	6	0
33	407	767	6	0
33	403	757	5	0
394	404	762	4	0
32	401	765	5	0
8	406	740	2	0
155	407	733	2	0
22	403	752	5	0
8	407	747	7	0
33	401	761	1	0
22	400	753	5	0
395	400	765	4	0
32	400	758	1	0
8	404	740	2	0
33	402	753	0	0
21	405	738	4	0
21	403	745	0	0
21	403	749	0	0
21	400	750	4	0
21	400	741	4	0
395	404	759	4	0
8	401	744	7	0
21	407	744	4	0
55	406	759	2	0
394	403	765	6	0
32	406	761	0	0
33	402	759	1	0
22	400	754	5	0
395	405	764	0	0
624	401	762	6	0
553	400	762	0	0
33	402	757	0	0
394	400	760	7	0
33	403	763	5	0
32	404	760	5	0
32	403	767	7	0
161	407	753	1	1
32	407	773	7	0
395	405	772	0	0
401	406	774	7	0
398	406	770	7	0
405	404	774	0	0
399	404	771	7	0
398	404	768	0	0
398	411	768	7	0
396	403	769	0	0
395	413	773	7	0
396	415	769	7	0
396	419	775	1	0
397	408	775	4	0
402	410	769	7	0
396	401	774	0	0
32	419	773	3	0
33	410	772	7	0
394	414	771	7	0
399	421	775	1	0
33	430	772	2	0
32	416	772	4	0
32	430	770	2	0
33	402	771	0	0
397	401	768	0	0
396	415	774	7	0
399	433	769	3	0
33	418	771	6	0
398	433	772	3	0
397	413	769	7	0
32	435	772	3	0
397	423	775	1	0
396	434	774	3	0
404	413	768	7	0
396	422	773	3	0
396	410	770	7	0
404	422	768	3	0
397	408	769	7	0
401	412	770	7	0
396	411	774	4	0
399	400	771	0	0
397	424	770	3	0
395	417	769	3	0
33	425	774	3	0
394	419	768	3	0
5	429	771	0	0
401	432	774	3	0
403	439	774	3	0
396	432	771	3	0
404	437	774	3	0
398	422	771	3	0
405	429	775	3	0
403	424	768	3	0
405	417	774	3	0
399	420	770	3	0
1	425	771	1	1
402	405	776	7	0
394	406	781	7	0
396	406	779	7	0
398	404	779	7	0
395	404	782	7	0
33	400	782	7	0
397	402	779	7	0
32	402	781	7	0
394	402	783	7	0
404	413	783	1	0
32	422	780	4	0
395	431	776	2	0
396	411	778	7	0
33	419	777	1	0
404	417	783	0	0
398	429	782	0	0
395	413	781	7	0
397	413	778	0	0
394	410	781	7	0
396	419	780	0	0
398	410	776	4	0
405	408	783	7	0
395	433	783	7	0
398	414	776	1	0
394	430	780	2	0
33	432	781	5	0
402	432	778	3	0
399	433	780	3	0
397	433	776	3	0
33	438	778	3	0
404	403	777	7	0
396	400	776	7	0
398	439	780	3	0
396	400	780	7	0
403	415	778	7	0
399	438	782	3	0
33	414	783	7	0
396	429	778	0	0
404	424	777	1	0
32	411	783	7	0
399	416	781	0	0
403	412	780	7	0
399	412	776	4	0
399	417	776	1	0
395	408	782	7	0
33	421	783	5	0
395	423	778	1	0
395	431	783	2	0
404	408	777	7	0
398	417	779	0	0
404	412	782	1	0
399	425	781	0	0
394	426	777	1	0
33	427	783	7	0
32	436	776	3	0
32	436	781	3	0
397	428	780	0	0
403	438	776	3	0
398	435	778	3	0
5	441	770	0	0
33	440	775	3	0
8	447	775	2	0
401	442	774	3	0
401	445	775	2	0
5	446	769	0	0
1	442	773	0	1
1	442	769	0	1
404	453	774	3	0
5	451	770	0	0
1	448	771	2	1
1	450	772	3	1
399	463	768	4	0
395	462	772	1	0
394	456	770	7	0
32	461	774	6	0
399	460	772	0	0
398	460	770	0	0
396	458	772	0	0
405	456	775	0	0
401	460	768	4	0
33	459	775	4	0
397	458	769	0	0
402	457	768	4	0
402	471	769	3	0
403	469	768	2	0
396	471	771	4	0
396	464	775	2	0
394	470	774	1	0
404	464	773	7	0
401	466	771	7	0
402	467	769	7	0
399	468	775	0	0
394	466	774	4	0
397	467	772	4	0
398	464	770	4	0
395	468	770	5	0
401	473	775	0	0
402	475	775	0	0
403	474	773	3	0
401	474	769	3	0
33	476	771	4	0
404	472	773	0	0
32	477	768	4	0
32	476	774	4	0
399	469	783	0	0
404	468	781	0	0
402	476	780	5	0
403	466	780	4	0
399	457	778	4	0
396	468	777	0	0
402	456	783	7	0
403	471	780	2	0
401	476	778	0	0
401	476	779	0	0
396	459	779	4	0
398	464	780	0	0
403	472	782	0	0
405	473	778	0	0
401	474	781	0	0
397	463	782	0	0
395	461	781	5	0
33	461	777	2	0
404	458	781	4	0
396	469	779	0	0
398	465	783	0	0
401	461	783	3	0
396	476	782	1	0
397	466	778	0	0
396	466	776	1	0
401	475	781	0	0
401	475	780	0	0
399	471	776	0	0
32	463	778	0	0
399	471	778	0	0
401	474	782	3	0
33	475	777	4	0
402	450	778	2	0
403	451	776	3	0
401	450	781	2	0
401	452	782	7	0
397	453	777	4	0
403	453	780	4	0
398	455	779	4	0
396	450	783	7	0
403	448	781	2	0
395	454	782	7	0
554	446	780	5	0
399	446	782	2	0
554	447	778	7	0
8	446	776	3	0
97	446	778	4	0
97	445	778	4	0
397	441	781	3	0
97	445	779	4	0
554	444	779	3	0
402	441	779	2	0
399	444	782	0	0
401	442	776	2	0
396	440	783	3	0
8	447	776	2	0
401	445	782	0	0
404	443	780	2	0
402	444	777	2	0
97	446	779	4	0
401	443	783	2	0
402	440	777	3	0
403	443	777	2	0
404	447	785	0	0
33	447	787	5	0
405	445	786	0	0
394	448	789	2	0
397	444	787	0	0
394	454	791	5	0
396	442	787	0	0
396	446	791	2	0
405	455	789	3	0
402	440	790	0	0
402	440	785	0	0
399	442	789	0	0
402	448	784	2	0
396	458	789	0	0
396	463	784	0	0
398	457	786	7	0
399	458	784	7	0
401	461	791	2	0
396	463	790	2	0
402	461	788	2	0
397	452	785	7	0
404	456	790	2	0
397	462	786	2	0
396	452	787	7	0
396	450	791	1	0
396	474	787	0	0
395	444	789	2	0
397	468	789	2	0
405	465	785	4	0
36	473	788	1	0
36	472	789	1	0
397	454	787	7	0
401	469	791	7	0
33	467	785	4	0
404	455	785	0	0
403	476	787	1	0
32	443	785	7	0
404	464	786	0	0
32	475	784	5	0
404	470	786	0	0
404	473	791	4	0
401	458	791	0	0
403	440	788	0	0
401	442	791	0	0
397	451	784	7	0
396	449	787	4	0
404	450	789	4	0
398	460	789	2	0
394	460	785	3	0
403	463	788	0	0
403	453	784	7	0
399	459	787	2	0
32	452	789	6	0
394	465	791	2	0
396	472	786	0	0
397	474	789	5	0
36	472	788	1	0
404	471	784	0	0
403	468	787	2	0
401	471	789	2	0
401	470	788	3	0
394	470	790	6	0
401	475	786	7	0
401	474	785	7	0
401	472	791	5	0
405	474	790	0	0
404	466	789	3	0
401	438	787	0	0
398	436	787	0	0
399	436	790	0	0
404	438	790	3	0
397	433	788	0	0
394	437	784	7	0
405	435	785	0	0
397	434	790	0	0
405	429	785	4	0
32	429	791	3	0
398	429	789	7	0
396	431	786	2	0
404	431	790	2	0
397	427	790	1	0
394	428	787	1	0
33	425	784	0	0
32	424	787	1	0
404	425	791	3	0
395	426	788	3	0
396	423	789	1	0
397	420	790	7	0
405	419	785	1	0
394	417	789	4	0
395	419	788	1	0
32	422	785	0	0
394	421	788	3	0
397	417	785	4	0
398	417	787	4	0
33	423	791	3	0
398	418	791	7	0
32	412	789	1	0
398	412	791	2	0
399	413	785	4	0
403	410	788	4	0
403	412	787	2	0
402	409	789	1	0
32	410	791	4	0
398	411	786	4	0
404	414	791	4	0
396	415	786	4	0
397	409	785	4	0
395	408	790	4	0
396	408	787	4	0
399	415	789	4	0
404	401	786	3	0
395	400	784	7	0
395	400	788	4	0
394	401	790	4	0
401	407	789	1	0
395	404	789	4	0
404	405	787	4	0
404	403	791	3	0
395	406	784	7	0
394	402	787	4	0
395	403	785	4	0
397	405	791	3	0
33	399	759	1	0
32	397	758	2	0
553	397	759	0	0
32	397	757	7	0
553	397	756	0	0
405	394	774	0	0
33	395	759	1	0
33	395	758	1	0
396	392	770	0	0
32	398	760	1	0
394	395	772	0	0
33	394	764	7	0
32	396	755	0	0
395	396	770	0	0
32	394	755	4	0
33	393	767	5	0
32	397	772	0	0
32	395	766	5	0
405	394	756	5	0
405	398	783	7	0
399	392	773	0	0
401	397	781	7	0
396	398	788	4	0
553	398	762	0	0
396	395	778	5	0
401	392	768	0	0
193	395	754	0	0
33	396	785	4	0
398	395	776	5	0
396	399	786	4	0
395	397	790	4	0
32	394	783	5	0
33	397	763	1	0
32	394	757	2	0
33	398	765	7	0
397	393	775	5	0
396	393	791	4	0
33	399	767	5	0
394	393	761	5	0
394	395	791	4	0
402	399	779	7	0
394	393	787	4	0
32	392	762	1	0
397	393	789	4	0
401	397	778	7	0
404	394	770	7	0
553	394	760	0	0
553	395	763	0	0
397	399	777	7	0
394	396	767	4	0
404	397	776	7	0
404	399	769	0	0
398	397	774	0	0
394	396	783	5	0
395	392	782	5	0
395	393	777	5	0
396	392	779	5	0
397	399	774	0	0
32	394	780	5	0
403	393	785	4	0
33	395	788	4	0
32	397	787	4	0
397	399	790	3	0
395	391	766	4	0
618	389	752	4	0
32	390	761	7	0
32	390	765	5	0
33	389	763	7	0
398	388	764	0	0
397	386	766	0	0
394	389	773	0	0
398	387	763	0	0
404	388	781	4	0
397	385	772	3	0
660	387	760	2	0
394	388	783	5	0
398	387	771	1	0
399	390	770	0	0
397	384	765	0	0
405	384	774	0	0
404	391	785	4	0
32	388	786	4	0
394	388	776	5	0
399	387	788	4	0
398	391	789	4	0
397	386	790	5	0
32	386	762	5	0
32	384	776	7	0
32	385	770	5	0
396	388	779	5	0
32	384	786	6	0
395	384	781	5	0
398	388	791	3	0
397	386	784	2	0
398	390	780	4	0
399	389	784	4	0
395	391	771	0	0
33	387	765	5	0
402	389	768	0	0
404	386	768	0	0
395	389	789	4	0
33	384	768	0	0
396	385	764	0	0
32	390	782	5	0
404	391	775	5	0
405	384	783	4	0
33	390	778	5	0
395	386	782	5	0
33	390	787	4	0
399	388	770	7	0
404	386	786	2	0
33	385	779	5	0
32	386	777	5	0
396	385	788	2	0
399	390	791	1	0
33	386	773	0	0
33	384	790	4	0
398	376	769	5	0
403	381	784	4	0
397	376	779	5	0
404	382	772	0	0
402	383	783	0	0
403	383	777	0	0
397	376	786	4	0
404	383	787	0	0
396	383	769	0	0
405	381	781	7	0
395	381	790	4	0
394	376	771	5	0
405	381	769	4	0
404	378	772	7	0
396	378	780	5	0
33	379	775	7	0
33	380	788	4	0
394	377	789	4	0
395	379	778	7	0
399	380	771	5	0
394	376	783	4	0
398	378	769	7	0
32	376	775	7	0
404	379	782	4	0
403	379	791	4	0
399	380	773	7	0
396	378	785	4	0
164	372	763	6	0
405	371	768	7	0
164	373	764	7	0
397	373	773	7	0
396	374	770	7	0
398	370	772	1	0
32	368	771	0	0
399	369	775	1	0
405	371	781	3	0
399	375	788	4	0
399	371	778	5	0
395	369	783	1	0
396	369	786	1	0
398	374	780	5	0
402	374	782	4	0
32	373	784	4	0
394	374	777	7	0
692	369	781	6	0
404	372	776	6	0
403	369	778	1	0
398	372	786	4	0
32	369	789	1	0
32	374	791	4	0
692	368	781	6	0
404	373	789	4	0
404	367	774	0	0
396	365	771	1	0
401	367	769	0	0
33	364	773	0	0
218	366	770	1	0
399	362	779	1	0
395	360	771	0	0
394	361	779	2	0
692	367	781	6	0
32	360	775	0	0
32	361	768	0	0
397	367	787	1	0
397	364	776	1	0
396	366	783	1	0
399	362	771	1	0
403	363	769	1	0
398	366	789	1	0
395	360	789	0	0
401	361	773	0	0
398	364	781	1	0
395	367	776	1	0
405	363	785	1	0
396	362	791	1	0
404	361	786	0	0
399	364	790	1	0
394	363	788	1	0
401	365	778	1	0
394	366	785	1	0
33	360	784	4	0
397	352	773	6	0
403	358	789	6	0
405	357	774	0	0
397	357	780	4	0
394	357	778	0	0
397	359	787	0	0
404	352	768	0	0
399	359	791	6	0
396	357	783	4	0
394	353	775	0	0
396	358	769	0	0
399	359	779	0	0
218	356	780	4	0
32	355	786	0	0
403	353	788	6	0
404	356	790	7	0
405	353	785	4	0
395	359	781	4	0
397	353	781	4	0
401	357	771	0	0
404	356	777	4	0
398	355	770	0	0
394	356	784	4	0
403	355	774	0	0
397	359	777	0	0
397	358	776	0	0
397	356	768	0	0
32	355	781	4	0
399	353	771	0	0
403	355	776	0	0
403	359	773	0	0
396	354	778	0	0
399	357	787	6	0
33	355	772	0	0
405	352	790	4	0
396	358	778	0	0
398	355	788	6	0
397	358	785	0	0
21	396	744	6	0
157	397	745	2	0
21	398	742	4	0
21	395	739	6	0
21	395	741	6	0
155	397	738	2	0
33	359	793	0	0
32	366	799	4	0
398	357	797	6	0
98	361	796	4	0
398	357	792	6	0
404	358	795	6	0
398	353	794	6	0
395	361	794	6	0
398	367	795	1	0
399	366	797	1	0
405	381	794	0	0
98	360	797	3	0
396	383	793	0	0
396	368	792	1	0
395	375	793	0	0
394	369	795	1	0
395	391	797	4	0
397	376	796	7	0
395	377	798	7	0
99	359	799	3	0
395	386	794	6	0
404	386	799	4	0
396	386	792	4	0
32	384	793	4	0
394	390	794	6	0
397	364	794	1	0
402	369	798	4	0
33	365	792	1	0
405	353	797	0	0
405	363	797	1	0
397	355	798	6	0
394	354	792	0	0
403	362	793	6	0
33	378	793	4	0
396	355	795	6	0
98	361	798	4	0
396	373	797	7	0
404	379	797	0	0
33	389	799	4	0
33	384	797	4	0
405	371	794	4	0
396	389	796	4	0
394	374	799	7	0
401	381	798	7	0
32	388	793	6	0
394	387	797	4	0
33	392	793	6	0
32	393	798	3	0
395	394	794	6	0
395	394	796	3	0
32	397	795	3	0
394	396	797	3	0
404	397	799	1	0
398	397	792	3	0
33	398	793	3	0
395	395	799	3	0
399	399	797	1	0
405	395	793	4	0
394	404	793	3	0
404	404	797	1	0
32	407	795	0	0
395	406	794	3	0
399	402	798	1	0
399	400	792	3	0
396	400	795	3	0
399	406	798	3	0
32	402	794	3	0
399	402	796	1	0
401	415	799	0	0
395	413	799	0	0
395	414	793	1	0
401	414	795	1	0
399	414	797	1	0
106	412	796	1	0
404	412	799	1	0
108	410	795	6	0
398	410	799	2	0
110	410	796	7	0
109	410	794	5	0
555	412	794	0	0
195	412	795	2	0
394	408	797	0	0
111	411	794	4	0
107	410	797	2	0
32	408	793	4	0
405	408	792	4	0
401	417	799	1	0
397	420	799	6	0
402	417	795	1	0
399	420	792	7	0
394	422	794	3	0
394	423	799	0	0
403	416	794	0	0
396	419	797	6	0
33	416	792	4	0
33	421	796	3	0
405	419	795	6	0
403	416	797	1	0
396	418	793	7	0
32	453	798	6	0
397	454	798	6	0
397	449	795	0	0
401	452	795	0	0
399	450	793	2	0
404	457	795	2	0
404	448	798	0	0
394	448	792	4	0
405	461	796	0	0
398	463	793	2	0
404	452	792	0	0
32	461	794	6	0
397	455	799	6	0
405	450	799	0	0
396	454	799	0	0
402	459	793	2	0
403	455	794	0	0
564	471	795	1	0
396	466	799	0	0
564	471	794	0	0
402	450	796	0	0
401	467	796	0	0
403	458	797	0	0
405	451	798	3	0
396	461	799	0	0
33	460	792	6	0
430	470	793	1	0
402	468	794	6	0
396	457	799	0	0
401	466	797	6	0
403	462	798	0	0
402	475	799	2	0
405	463	797	2	0
564	470	794	7	0
218	476	793	6	0
397	474	795	5	0
397	474	793	5	0
396	476	795	7	0
401	476	794	6	0
32	473	796	2	0
396	469	795	6	0
396	466	792	2	0
396	472	799	4	0
36	473	797	1	0
36	472	795	1	0
402	472	797	6	0
36	472	796	1	0
399	459	799	0	0
394	459	797	5	0
401	471	796	6	0
398	457	798	0	0
395	456	797	2	0
32	452	798	6	0
402	456	792	0	0
404	468	797	0	0
403	459	795	0	0
403	464	799	0	0
401	470	795	1	0
402	465	797	6	0
402	470	796	0	0
32	476	792	0	0
402	474	799	2	0
401	475	793	5	0
402	464	797	6	0
397	470	799	0	0
396	475	797	6	0
396	472	793	1	0
401	472	794	5	0
401	473	793	0	0
396	447	796	0	0
398	446	794	2	0
397	444	793	2	0
403	444	798	0	0
395	440	799	0	0
397	444	796	6	0
33	442	794	0	0
404	442	798	6	0
34	441	796	0	0
399	439	793	0	0
405	439	796	3	0
396	437	798	6	0
402	439	798	6	0
398	437	794	0	0
32	437	792	0	0
397	435	798	6	0
396	433	792	0	0
403	436	796	0	0
394	434	799	0	0
397	435	793	0	0
33	433	794	0	0
404	433	796	0	0
395	431	792	2	0
394	431	796	2	0
394	430	799	2	0
396	431	794	2	0
396	429	798	6	0
405	429	795	0	0
395	427	799	2	0
399	429	793	7	0
400	424	794	7	0
404	424	796	7	0
32	427	796	3	0
32	425	797	7	0
395	426	794	3	0
403	427	792	1	0
33	354	801	6	0
32	352	800	6	0
397	368	804	4	0
394	373	805	7	0
398	366	806	4	0
32	356	800	6	0
395	369	807	4	0
399	373	802	7	0
396	360	803	4	0
33	370	800	4	0
396	371	803	4	0
32	376	806	4	0
399	366	802	4	0
394	380	804	7	0
398	379	801	7	0
397	377	803	7	0
394	390	803	4	0
403	364	800	4	0
396	387	806	4	0
404	388	802	2	0
32	384	801	4	0
32	371	806	4	0
398	361	800	4	0
396	365	804	4	0
404	362	806	4	0
33	360	807	4	0
396	383	800	0	0
404	368	801	4	0
33	372	800	7	0
397	363	803	4	0
396	388	804	2	0
32	381	802	7	0
395	381	806	0	0
397	385	807	3	0
405	384	804	4	0
397	390	805	2	0
33	386	803	2	0
405	376	800	0	0
32	390	807	2	0
401	385	805	2	0
32	391	801	4	0
33	399	802	3	0
32	397	806	3	0
396	399	806	3	0
33	397	804	3	0
403	399	800	1	0
394	397	802	3	0
405	395	803	3	0
394	395	805	3	0
32	394	807	4	0
404	392	806	4	0
394	394	801	3	0
403	392	803	2	0
404	407	805	1	0
399	406	800	1	0
394	407	807	2	0
398	405	807	1	0
397	402	801	1	0
33	405	803	3	0
395	401	803	3	0
404	405	805	3	0
398	404	800	1	0
33	403	804	3	0
396	403	807	1	0
396	401	800	1	0
32	400	805	3	0
404	401	807	3	0
395	407	803	2	0
404	411	805	7	0
32	411	801	2	0
395	409	806	2	0
394	411	807	2	0
397	415	802	2	0
395	413	806	2	0
394	415	807	2	0
397	409	801	2	0
396	412	801	2	0
33	414	801	0	0
394	415	804	2	0
405	408	801	1	0
394	410	803	2	0
395	413	803	2	0
403	422	804	5	0
398	420	803	6	0
404	421	806	0	0
33	422	801	0	0
396	420	801	6	0
402	416	801	1	0
405	419	805	3	0
403	417	804	1	0
403	418	802	1	0
396	431	807	2	0
398	429	803	7	0
403	431	804	2	0
404	427	804	5	0
396	426	806	0	0
395	425	802	7	0
32	428	801	7	0
404	431	801	2	0
404	424	807	7	0
404	425	800	6	0
396	424	804	6	0
397	429	800	6	0
405	429	805	1	0
397	429	807	7	0
33	437	801	0	0
394	439	804	0	0
396	436	803	6	0
397	436	805	3	0
401	435	806	1	0
399	433	805	6	0
398	433	803	6	0
402	432	807	1	0
397	434	801	6	0
404	446	800	0	0
396	444	801	6	0
399	445	805	0	0
398	444	803	0	0
34	441	802	0	0
33	440	802	0	0
396	442	801	6	0
399	447	807	3	0
397	446	803	0	0
394	442	805	6	0
397	455	800	6	0
394	453	804	3	0
405	452	802	3	0
396	448	801	0	0
394	455	804	3	0
396	452	806	3	0
404	453	802	3	0
405	451	800	5	0
397	450	807	3	0
32	451	804	3	0
404	454	807	3	0
398	449	806	3	0
396	454	800	6	0
401	454	803	3	0
194	453	800	0	0
394	458	801	0	0
32	462	805	0	0
33	457	804	0	0
397	463	801	0	0
396	461	807	0	0
398	459	803	0	0
397	458	800	0	0
404	460	801	0	0
32	456	802	3	0
395	461	802	0	0
405	459	807	0	0
399	460	805	0	0
402	467	807	0	0
398	471	806	0	0
398	467	802	0	0
33	464	805	0	0
405	471	804	0	0
33	471	800	0	0
397	466	806	0	0
396	468	804	0	0
32	468	800	0	0
32	466	803	5	0
399	464	802	0	0
33	469	807	3	0
397	470	802	0	0
399	473	806	0	0
402	477	805	2	0
404	475	800	2	0
401	477	804	2	0
402	475	802	5	0
402	477	806	2	0
401	476	800	7	0
402	476	803	2	0
401	476	801	7	0
401	477	807	0	0
394	475	805	0	0
398	475	807	0	0
396	473	801	5	0
32	473	804	0	0
398	436	810	1	0
33	439	812	3	0
394	436	814	1	0
401	445	808	3	0
395	433	814	1	0
33	434	808	1	0
397	452	811	3	0
396	434	810	1	0
397	432	810	1	0
32	434	812	1	0
396	462	813	3	0
395	443	814	2	0
396	454	810	3	0
192	443	810	0	0
404	461	809	3	0
399	461	815	3	0
405	453	813	3	0
32	460	813	3	0
32	451	809	3	0
405	469	813	0	0
33	456	811	3	0
32	471	813	0	0
404	468	809	0	0
398	471	814	0	0
404	468	815	0	0
404	470	810	4	0
401	471	808	0	0
401	438	809	1	0
404	432	813	1	0
395	438	810	3	0
395	446	812	3	0
32	446	808	3	0
32	437	808	3	0
33	437	812	1	0
395	464	808	0	0
32	466	812	0	0
396	440	811	1	0
394	445	810	3	0
396	455	813	1	0
399	459	814	3	0
398	449	810	3	0
401	460	811	3	0
394	441	812	3	0
394	448	813	3	0
399	450	812	3	0
399	463	810	7	0
394	440	814	1	0
32	466	814	0	0
399	468	811	0	0
403	466	810	0	0
402	458	812	3	0
397	464	814	3	0
396	459	809	3	0
395	353	814	2	0
401	372	811	4	0
33	355	812	2	0
404	373	808	0	0
397	373	814	0	0
394	364	809	4	0
395	378	811	0	0
396	381	812	0	0
32	388	813	4	0
394	380	813	0	0
399	389	808	2	0
397	375	809	0	0
396	382	809	0	0
32	358	813	2	0
394	374	811	0	0
404	388	810	2	0
405	384	813	4	0
403	385	810	2	0
32	366	812	4	0
398	391	808	2	0
98	359	809	1	0
404	369	811	0	0
398	377	808	0	0
403	370	809	0	0
396	370	813	0	0
405	362	814	0	0
405	375	813	0	0
395	390	812	4	0
399	379	809	0	0
33	361	811	4	0
397	382	814	0	0
398	387	808	4	0
394	386	812	4	0
33	386	809	2	0
396	390	810	2	0
398	350	782	2	0
400	351	783	2	0
395	351	777	6	0
573	351	781	4	0
399	350	783	2	0
394	347	778	4	0
573	350	781	0	0
32	349	777	6	0
394	349	779	6	0
394	350	784	0	0
397	346	779	4	0
396	347	781	4	0
403	349	787	4	0
33	344	778	4	0
399	351	787	4	0
399	350	791	4	0
397	346	791	4	0
32	350	795	4	0
398	347	784	4	0
405	349	805	4	0
396	349	799	4	0
33	348	796	4	0
403	345	787	4	0
403	351	804	4	0
98	349	807	1	0
398	348	792	4	0
33	350	802	6	0
403	344	801	0	0
404	346	803	4	0
701	346	807	0	0
98	344	806	0	0
404	345	811	4	0
399	348	810	4	0
396	350	789	4	0
403	347	812	4	0
396	345	776	4	0
397	345	783	4	0
405	344	781	0	0
398	347	787	4	0
395	351	798	4	0
399	349	785	4	0
399	344	789	2	0
395	345	785	4	0
404	348	790	4	0
401	346	789	4	0
401	345	795	2	0
396	347	794	4	0
394	351	793	4	0
399	344	792	2	0
397	347	800	4	0
399	348	804	4	0
405	345	798	4	0
398	348	814	4	0
33	349	808	4	0
701	347	806	0	0
32	351	808	4	0
397	348	802	4	0
394	351	812	2	0
398	398	808	3	0
394	397	811	7	0
395	395	809	3	0
394	395	813	3	0
33	398	813	3	0
399	392	810	2	0
398	399	811	0	0
33	394	811	4	0
33	392	813	4	0
397	405	809	1	0
405	406	813	3	0
394	404	813	3	0
395	405	811	7	0
396	407	809	1	0
399	400	809	3	0
395	400	813	3	0
404	403	810	1	0
32	402	813	7	0
398	403	812	6	0
398	401	811	6	0
394	412	812	3	0
32	408	812	3	0
394	415	810	3	0
395	408	810	3	0
395	413	809	3	0
395	410	813	3	0
32	410	811	3	0
405	414	813	0	0
404	411	809	1	0
395	420	808	0	0
395	423	812	0	0
398	421	811	7	0
399	423	809	7	0
394	421	813	0	0
395	419	812	0	0
401	418	808	7	0
394	417	812	0	0
398	418	810	7	0
33	416	813	0	0
396	416	809	1	0
403	428	811	7	0
405	429	813	0	0
399	430	810	2	0
398	429	809	0	0
32	426	813	0	0
398	425	809	7	0
394	431	812	2	0
394	427	808	0	0
399	426	811	7	0
399	476	813	0	0
33	476	811	3	0
401	477	808	0	0
397	474	809	0	0
397	473	813	0	0
398	474	811	0	0
403	472	811	6	0
401	476	809	0	0
396	475	814	0	0
396	472	808	0	0
32	475	812	0	0
401	476	808	0	0
404	410	817	4	0
402	411	817	4	0
32	423	816	6	0
395	410	821	0	0
403	412	816	4	0
33	408	819	0	0
395	431	816	2	0
33	431	819	2	0
395	420	817	6	0
32	421	817	6	0
395	419	820	4	0
32	420	823	2	0
396	414	821	7	0
401	408	817	4	0
99	408	816	2	0
395	416	816	6	0
396	416	817	6	0
33	415	816	6	0
588	428	820	6	0
400	435	816	1	0
395	424	816	6	0
33	434	818	0	0
588	427	819	6	0
32	408	821	0	0
33	412	819	1	0
401	409	823	0	0
97	412	822	0	0
394	414	816	6	0
396	423	817	6	0
404	435	820	0	0
395	409	816	4	0
403	439	819	0	0
397	439	821	0	0
401	412	817	4	0
32	439	817	0	0
405	413	816	4	0
588	428	819	6	0
33	435	819	1	0
394	427	816	2	0
404	445	821	1	0
33	432	817	6	0
395	447	821	1	0
32	447	819	1	0
588	428	821	6	0
398	444	819	1	0
396	422	821	3	0
33	419	817	6	0
394	422	817	2	0
395	415	819	6	0
98	430	822	1	0
32	416	820	5	0
395	429	816	2	0
32	417	817	6	0
32	430	817	2	0
394	418	817	6	0
32	410	816	4	0
394	441	817	0	0
395	442	819	0	0
588	426	819	6	0
32	426	817	6	0
32	425	816	6	0
401	432	823	0	0
395	437	820	0	0
396	433	816	7	0
396	436	818	1	0
399	433	821	0	0
396	425	817	6	0
397	437	816	1	0
398	435	822	0	0
397	443	817	1	0
32	432	820	6	0
33	441	822	0	0
32	443	823	1	0
396	445	817	1	0
33	405	823	3	0
99	407	817	4	0
99	406	816	4	0
402	406	817	4	0
32	407	816	4	0
98	402	817	4	0
396	401	821	6	0
99	401	817	5	0
33	405	816	4	0
98	402	816	3	0
395	404	816	4	0
32	403	820	5	0
395	404	822	4	0
99	403	816	4	0
403	400	817	3	0
33	398	819	3	0
395	398	822	2	0
395	399	817	3	0
32	398	817	3	0
399	395	817	3	0
32	395	820	4	0
395	395	816	2	0
401	393	816	3	0
395	392	819	5	0
396	394	816	3	0
21	397	823	0	0
394	397	816	3	0
33	396	816	3	0
401	397	817	3	0
21	396	822	0	0
32	393	818	6	0
21	394	822	0	0
33	392	816	3	0
21	393	823	0	0
32	386	817	0	0
98	388	818	2	0
98	389	818	3	0
403	389	817	5	0
33	385	820	0	0
395	384	818	0	0
395	388	822	0	0
404	385	818	0	0
401	390	821	2	0
99	390	818	5	0
34	391	817	3	0
32	387	820	1	0
403	388	816	3	0
395	391	816	3	0
32	390	817	3	0
32	384	816	0	0
405	387	818	5	0
33	386	822	7	0
394	381	816	0	0
396	376	817	0	0
398	376	821	0	0
32	379	818	0	0
403	381	820	0	0
397	378	820	0	0
404	369	822	1	0
394	373	819	0	0
401	375	819	0	0
399	371	821	0	0
395	368	816	1	0
404	371	816	1	0
398	374	821	0	0
33	368	819	1	0
395	366	823	1	0
404	366	819	1	0
395	366	821	1	0
32	365	817	1	0
32	363	820	4	0
395	362	817	4	0
405	364	822	4	0
403	349	817	7	0
33	348	820	1	0
403	345	821	7	0
398	350	819	7	0
396	355	818	4	0
403	352	822	4	0
394	357	817	4	0
651	353	816	0	0
405	346	817	1	0
395	347	823	1	0
397	344	819	6	0
401	350	823	7	0
395	359	817	4	0
399	352	817	0	0
395	344	817	4	0
397	354	820	4	0
401	352	816	0	0
33	353	829	0	0
405	355	827	1	0
98	364	827	4	0
33	371	827	0	0
396	351	825	1	0
33	350	827	1	0
405	362	828	1	0
100	363	828	2	0
32	357	828	4	0
394	351	830	0	0
394	367	828	1	0
402	348	825	7	0
397	375	829	0	0
404	374	826	4	0
403	373	828	4	0
402	368	826	4	0
32	381	827	0	0
403	378	826	4	0
399	381	830	0	0
401	379	831	0	0
395	355	831	0	0
395	363	831	0	0
402	370	826	4	0
394	358	830	0	0
404	354	825	4	0
404	346	827	1	0
99	365	827	2	0
98	364	825	1	0
405	346	830	0	0
396	373	824	0	0
396	377	830	0	0
394	372	830	0	0
403	382	824	4	0
32	366	825	1	0
395	379	824	0	0
398	370	824	0	0
403	376	824	4	0
405	376	827	0	0
98	361	828	1	0
399	379	829	0	0
395	389	826	0	0
396	390	831	7	0
21	391	827	0	0
21	391	825	0	0
32	389	824	0	0
396	385	825	5	0
397	385	828	4	0
395	385	829	3	0
396	387	824	3	0
394	384	825	3	0
32	385	824	6	0
395	388	830	0	0
395	384	827	3	0
32	387	828	2	0
397	384	828	3	0
397	384	829	3	0
21	399	827	0	0
21	398	824	0	0
5	396	826	0	0
21	394	830	0	0
21	392	824	0	0
21	399	825	0	0
21	392	828	0	0
21	398	828	0	0
21	396	830	0	0
21	397	829	0	0
33	397	831	1	0
21	393	829	0	0
33	406	830	6	0
21	407	826	0	0
395	403	826	2	0
401	402	829	0	0
33	401	824	7	0
403	400	829	0	0
775	400	826	0	0
394	411	824	3	0
21	411	830	5	0
32	411	826	4	0
395	413	824	2	0
401	422	828	3	0
32	420	831	5	0
395	419	826	1	0
395	422	830	0	0
33	416	824	0	0
2	417	827	3	1
32	431	831	6	0
399	430	824	6	0
33	424	828	2	0
394	426	828	1	0
399	431	826	6	0
398	430	828	6	0
403	426	830	0	0
396	424	831	7	0
395	429	831	1	0
32	429	826	6	0
394	439	829	1	0
33	438	827	0	0
32	434	826	0	0
401	434	830	1	0
33	434	824	0	0
395	435	829	1	0
33	432	825	6	0
394	436	827	0	0
32	437	830	1	0
32	438	824	0	0
32	433	829	5	0
394	436	825	0	0
404	432	827	1	0
33	447	824	1	0
405	441	826	0	0
32	440	824	0	0
402	440	831	1	0
33	446	830	1	0
32	442	830	1	0
404	449	829	6	0
710	451	827	5	0
402	444	829	1	0
710	451	829	1	0
710	451	828	0	0
396	452	827	6	0
394	455	818	1	0
394	444	825	1	0
396	446	826	1	0
33	451	818	1	0
397	444	827	1	0
398	444	831	1	0
394	447	828	1	0
32	452	816	1	0
398	452	820	1	0
403	450	816	1	0
394	449	819	1	0
395	454	817	1	0
396	450	822	1	0
33	455	816	1	0
33	449	818	1	0
395	453	818	1	0
395	441	828	1	0
401	448	830	1	0
405	449	826	1	0
403	452	828	2	0
399	448	816	1	0
395	448	819	1	0
405	448	818	1	0
401	452	823	1	0
710	450	827	6	0
710	450	829	2	0
32	449	827	6	0
710	450	828	7	0
32	457	819	1	0
691	459	828	2	0
33	463	819	1	0
396	459	816	1	0
397	462	817	1	0
397	457	817	1	0
397	459	818	1	0
394	461	818	1	0
33	461	816	1	0
32	460	819	1	0
32	469	816	1	0
399	470	818	1	0
394	467	818	1	0
396	471	830	5	0
396	468	818	1	0
403	471	825	5	0
164	465	831	7	0
395	471	828	5	0
164	465	826	2	0
32	464	816	1	0
164	465	825	1	0
164	464	831	0	0
33	471	817	1	0
404	469	829	5	0
398	465	818	1	0
397	467	816	1	0
33	466	819	1	0
395	465	817	1	0
32	467	831	5	0
33	467	824	5	0
394	468	826	5	0
164	464	830	6	0
164	464	826	4	0
396	473	827	5	0
405	474	825	4	0
688	473	825	4	0
32	472	820	1	0
395	474	817	1	0
397	473	817	1	0
397	474	819	1	0
398	472	818	1	0
405	472	824	0	0
32	446	836	1	0
397	447	833	1	0
397	446	839	1	0
396	444	839	1	0
396	445	833	1	0
33	442	837	1	0
397	441	839	1	0
398	471	833	5	0
399	442	832	1	0
404	443	834	1	0
405	448	838	1	0
399	471	839	5	0
403	474	836	5	0
403	449	832	1	0
394	473	834	5	0
403	467	835	5	0
33	474	832	5	0
395	466	833	5	0
394	441	835	1	0
398	449	834	1	0
33	469	838	5	0
395	449	836	1	0
20	470	835	5	0
32	451	838	1	0
394	467	839	5	0
397	469	832	5	0
395	439	836	1	0
32	432	838	1	0
32	438	835	1	0
397	438	833	1	0
33	432	836	6	0
404	436	837	1	0
33	434	836	1	0
399	433	832	1	0
32	432	833	5	0
33	434	838	1	0
405	439	838	1	0
398	436	835	1	0
396	436	832	1	0
394	433	834	1	0
396	429	834	6	0
396	430	837	6	0
396	429	832	1	0
397	431	835	6	0
98	428	838	6	0
395	427	834	1	0
395	426	833	5	0
32	424	834	4	0
33	427	835	1	0
403	427	836	1	0
32	428	833	1	0
396	428	836	1	0
405	428	837	6	0
399	351	770	7	0
395	350	775	7	0
398	349	773	6	0
397	349	769	6	0
404	349	771	6	0
98	346	769	1	0
98	346	768	7	0
98	347	768	1	0
397	344	773	4	0
98	345	769	1	0
397	347	771	5	0
98	344	770	3	0
32	346	773	4	0
98	344	769	6	0
396	348	775	6	0
98	345	770	2	0
396	340	778	4	0
33	336	774	5	0
396	337	780	4	0
404	339	790	4	0
32	342	782	4	0
399	342	772	4	0
397	341	786	4	0
398	340	780	4	0
32	338	772	5	0
394	338	782	4	0
32	338	788	4	0
398	339	785	4	0
33	342	794	4	0
399	337	784	4	0
404	340	792	2	0
396	343	799	0	0
405	339	795	2	0
396	342	788	2	0
405	336	786	4	0
396	340	799	0	0
402	343	803	0	0
32	341	801	0	0
98	343	770	1	0
395	343	775	4	0
404	339	801	0	0
399	340	803	0	0
397	342	779	4	0
396	339	774	4	0
395	340	783	4	0
404	338	778	4	0
403	337	776	4	0
397	341	776	4	0
397	336	778	4	0
401	343	786	4	0
395	342	790	4	0
396	342	784	4	0
399	343	796	2	0
403	337	790	2	0
33	337	799	0	0
397	339	797	0	0
394	341	797	2	0
396	338	805	0	0
395	337	803	0	0
394	337	793	4	0
399	342	792	2	0
395	337	796	2	0
396	343	809	4	0
398	339	810	4	0
397	341	811	4	0
405	342	813	4	0
396	338	812	4	0
32	336	809	4	0
404	337	814	4	0
33	336	811	4	0
32	336	818	4	0
401	342	821	7	0
396	343	823	1	0
403	336	822	4	0
32	340	820	1	0
403	340	817	7	0
397	343	826	7	0
33	336	831	4	0
32	341	828	1	0
32	337	828	4	0
394	339	825	1	0
33	341	839	0	0
403	338	836	4	0
404	340	834	4	0
32	357	838	0	0
401	354	838	4	0
401	359	839	0	0
396	356	834	0	0
395	342	832	0	0
403	349	837	4	0
401	349	832	4	0
32	336	837	4	0
32	347	835	0	0
401	344	835	4	0
32	352	833	0	0
403	358	836	0	0
404	352	836	4	0
395	342	842	4	0
403	358	845	0	0
396	342	845	4	0
401	355	845	4	0
401	356	840	0	0
404	354	847	4	0
398	348	841	4	0
397	341	847	4	0
405	346	843	0	0
404	347	846	4	0
396	351	844	4	0
33	345	840	0	0
394	351	840	0	0
32	337	841	4	0
399	354	843	4	0
398	349	843	4	0
403	349	847	4	0
404	339	846	4	0
405	339	843	4	0
396	357	843	0	0
403	353	841	4	0
405	338	850	4	0
404	349	855	4	0
403	341	854	4	0
402	342	849	4	0
395	354	854	4	0
401	340	850	4	0
397	341	852	4	0
32	354	850	4	0
402	348	850	4	0
396	358	851	0	0
395	343	855	4	0
405	350	849	4	0
403	346	853	4	0
397	337	848	4	0
401	351	854	4	0
396	338	853	4	0
32	336	855	4	0
401	350	852	4	0
32	346	851	4	0
401	344	852	4	0
397	345	848	4	0
398	358	854	0	0
394	342	857	4	0
404	340	856	4	0
403	342	860	4	0
405	337	861	4	0
399	339	859	4	0
394	347	856	4	0
397	349	859	4	0
33	354	862	4	0
399	355	856	4	0
397	337	857	4	0
398	346	858	4	0
398	351	858	4	0
33	345	861	4	0
396	354	858	4	0
33	336	859	4	0
394	358	859	0	0
405	350	861	0	0
401	352	856	4	0
401	359	856	0	0
399	345	856	4	0
32	358	862	0	0
402	357	857	0	0
1091	359	869	0	0
1086	358	869	1	0
1092	356	871	4	0
1092	357	869	2	0
1086	355	869	3	0
1091	356	869	3	0
1091	354	871	4	0
1092	358	867	2	0
397	357	865	5	0
1091	356	867	3	0
396	351	865	5	0
1086	358	871	4	0
204	358	868	2	1
204	356	868	3	1
204	357	867	0	1
204	355	871	0	1
204	359	867	0	1
204	359	871	0	1
204	357	871	0	1
394	357	876	7	0
395	359	875	7	0
397	344	876	5	0
396	355	878	5	0
403	349	872	5	0
32	345	873	5	0
395	348	877	5	0
399	340	877	5	0
34	353	875	5	0
395	355	874	5	0
394	352	879	7	0
204	354	872	3	1
204	358	872	3	1
204	356	872	2	1
397	360	852	0	0
394	363	841	0	0
1091	360	871	4	0
404	363	853	0	0
395	361	844	0	0
397	360	847	0	0
397	360	841	0	0
395	360	862	0	0
1176	363	877	0	0
33	361	857	0	0
399	362	855	0	0
395	363	873	7	0
405	363	850	0	0
403	361	871	5	0
397	361	875	7	0
34	367	875	7	0
404	365	843	0	0
397	365	847	0	0
403	365	849	0	0
397	364	857	0	0
403	363	847	0	0
396	367	850	0	0
398	360	864	4	0
395	366	859	0	0
1092	360	869	7	0
398	363	859	0	0
1091	360	867	1	0
396	365	855	0	0
32	363	861	0	0
397	366	852	0	0
204	360	868	3	1
204	360	872	2	1
404	360	832	0	0
397	360	837	0	0
399	362	835	0	0
398	366	839	0	0
401	367	832	4	0
405	363	838	0	0
403	366	835	0	0
396	365	832	0	0
398	341	883	5	0
405	356	882	5	0
395	351	883	5	0
394	350	884	5	0
396	355	886	7	0
32	352	886	5	0
32	367	885	0	0
32	360	886	7	0
396	347	882	5	0
430	343	884	0	0
394	357	886	7	0
398	345	882	5	0
397	344	887	5	0
397	359	880	5	0
1176	363	882	0	0
395	364	885	7	0
395	366	881	7	0
394	343	880	5	0
399	373	855	0	0
401	372	853	0	0
404	373	850	0	0
398	368	854	0	0
395	369	852	0	0
405	375	861	4	0
397	375	858	0	0
32	370	884	5	0
399	372	881	4	0
405	375	884	1	0
403	375	872	5	0
33	374	857	0	0
404	373	878	4	0
32	369	859	0	0
33	368	862	0	0
430	369	876	0	0
396	370	878	4	0
1176	369	872	0	0
33	370	861	0	0
397	370	857	0	0
397	372	875	4	0
403	372	862	4	0
396	372	858	0	0
405	375	844	4	0
32	373	842	0	0
401	375	847	0	0
33	369	843	4	0
33	372	847	0	0
33	371	841	0	0
402	368	847	0	0
396	373	838	3	0
32	370	839	0	0
402	368	834	4	0
402	373	835	4	0
403	370	837	4	0
394	370	833	4	0
401	371	835	4	0
395	376	832	0	0
33	381	845	0	0
32	381	843	4	0
401	376	838	4	0
32	379	847	0	0
399	376	841	4	0
401	381	833	0	0
397	379	839	4	0
403	378	845	0	0
404	377	835	4	0
395	381	837	4	0
396	378	850	0	0
398	379	854	0	0
399	376	852	0	0
33	382	856	0	0
401	381	862	4	0
33	378	858	0	0
397	379	861	4	0
396	379	857	0	0
32	381	859	0	0
403	376	856	0	0
402	377	862	4	0
403	383	876	5	0
394	379	874	5	0
395	380	878	5	0
398	377	877	4	0
396	390	841	0	0
404	390	844	7	0
33	391	845	0	0
164	388	843	2	0
32	389	844	0	0
401	385	842	6	0
402	385	843	6	0
33	390	861	2	0
405	385	841	6	0
164	391	858	6	0
403	385	844	6	0
400	391	861	2	0
32	391	857	4	0
32	388	860	2	0
98	386	858	2	0
98	387	859	2	0
394	387	841	0	0
98	386	857	1	0
395	384	865	0	0
405	387	874	0	0
395	391	873	1	0
32	389	875	5	0
164	386	847	7	0
99	385	845	6	0
404	386	846	0	0
397	385	840	2	0
164	390	857	3	0
32	384	848	3	0
611	388	851	0	0
613	384	851	0	0
98	384	845	5	0
394	389	860	2	0
396	390	860	2	0
402	386	848	0	0
100	384	846	4	0
99	384	847	3	0
402	385	858	2	0
98	385	857	2	0
395	388	859	2	0
399	385	873	7	0
403	385	860	2	0
394	387	866	0	0
396	385	877	6	0
404	389	878	5	0
192	391	836	0	0
98	384	833	3	0
98	384	836	2	0
33	387	832	1	0
98	384	835	0	0
400	385	838	2	0
32	388	839	1	0
98	385	837	6	0
395	385	839	2	0
98	384	832	2	0
32	385	836	0	0
99	386	837	6	0
403	390	838	2	0
192	389	834	0	0
395	393	843	5	0
3	399	842	4	0
21	394	836	4	0
192	399	836	0	0
21	394	835	6	0
29	399	852	4	0
23	398	842	6	0
192	396	833	0	0
5	395	838	0	0
712	394	851	0	0
32	395	844	6	0
403	393	853	0	0
118	399	840	0	0
404	393	850	0	0
164	393	855	0	0
33	393	832	0	0
405	393	848	2	0
33	392	838	3	0
3	399	843	4	0
395	399	847	0	0
401	394	854	0	0
21	396	836	4	0
21	396	835	6	0
396	394	846	5	0
394	393	841	4	0
164	392	846	0	0
33	392	844	6	0
1	395	840	0	1
1	399	851	1	1
98	399	861	7	0
99	398	861	6	0
32	396	860	5	0
396	394	860	5	0
98	397	861	5	0
164	392	857	5	0
402	392	861	2	0
164	392	856	0	0
401	392	860	0	0
395	395	860	5	0
401	394	861	5	0
405	393	861	0	0
395	397	860	5	0
402	393	860	0	0
397	399	865	0	0
1091	399	870	0	0
398	393	865	0	0
1092	398	867	1	0
1100	397	867	0	0
1099	396	870	0	0
1091	396	869	2	0
1091	395	870	4	0
396	397	864	0	0
1100	398	869	2	0
1086	395	867	1	0
1099	394	869	2	0
1086	398	870	1	0
1092	394	870	5	0
1092	397	870	2	0
394	395	865	0	0
204	399	867	3	1
204	397	869	2	1
204	393	871	3	1
204	395	869	3	1
204	399	869	3	1
204	396	867	2	1
204	394	867	3	1
404	397	877	0	0
403	399	879	0	0
1091	398	872	5	0
403	393	877	5	0
1099	394	872	3	0
1092	396	872	4	0
399	399	874	2	0
1176	394	874	0	0
397	396	875	1	0
204	395	873	0	1
204	399	873	0	1
204	397	873	0	1
33	407	860	5	0
401	401	860	5	0
395	402	861	5	0
33	403	861	5	0
1086	400	867	5	0
394	404	861	5	0
395	405	849	7	0
404	405	861	5	0
395	401	864	0	0
290	402	842	6	0
29	403	850	4	0
1099	400	869	2	0
1086	400	870	7	0
32	407	853	1	0
32	406	861	5	0
394	406	856	2	0
33	402	832	7	0
32	400	861	5	0
403	402	856	5	0
396	401	861	5	0
33	405	857	3	0
1100	401	870	0	0
395	406	860	5	0
779	404	846	5	0
404	407	843	3	0
1	405	851	1	1
2	402	855	0	1
1	407	847	2	1
165	400	845	0	1
1	401	848	0	1
204	401	867	2	1
204	402	871	2	1
204	401	869	2	1
29	414	834	2	0
394	414	854	7	0
32	415	862	0	0
403	409	833	0	0
98	411	861	5	0
33	415	841	0	0
396	414	844	1	0
396	411	858	2	0
32	415	861	0	0
33	409	861	5	0
403	409	855	4	0
394	408	864	0	0
33	412	860	3	0
403	414	837	0	0
33	414	858	2	0
395	414	861	0	0
394	409	860	5	0
32	413	859	1	0
396	414	851	5	0
98	412	861	5	0
401	413	861	3	0
403	412	862	4	0
396	412	845	0	0
396	412	854	6	0
98	413	860	5	0
394	413	847	1	0
99	410	861	5	0
33	408	858	3	0
404	411	841	3	0
396	408	860	5	0
2	412	850	2	1
1	411	834	1	1
2	413	833	0	1
1	415	847	1	1
32	418	836	7	0
394	419	833	6	0
33	422	833	6	0
396	421	835	3	0
98	429	840	1	0
32	434	843	6	0
98	428	843	1	0
401	424	843	1	0
396	435	847	2	0
33	438	843	0	0
33	418	841	7	0
32	421	842	6	0
403	445	847	5	0
32	428	847	2	0
397	432	843	5	0
32	428	845	1	0
395	422	844	2	0
394	428	846	2	0
32	445	845	3	0
397	443	841	1	0
98	430	841	1	0
403	429	843	1	0
397	436	841	7	0
33	426	845	1	0
395	429	842	1	0
99	428	844	1	0
395	430	847	6	0
399	439	847	5	0
404	437	846	1	0
33	429	841	1	0
32	430	845	6	0
395	433	845	4	0
401	442	846	2	0
394	446	841	1	0
397	441	843	7	0
394	440	845	6	0
1	418	844	0	1
2	422	847	1	1
32	423	851	5	0
395	429	853	2	0
33	417	852	6	0
394	422	855	3	0
33	423	854	2	0
768	419	853	2	0
395	430	852	6	0
394	427	848	2	0
394	428	851	2	0
395	429	848	2	0
396	431	855	6	0
395	437	851	1	0
396	428	848	2	0
394	426	854	0	0
396	416	855	7	0
396	428	852	2	0
33	430	850	6	0
98	428	850	2	0
32	431	853	6	0
394	431	849	6	0
396	424	849	4	0
32	437	848	4	0
398	443	848	2	0
402	441	848	2	0
397	440	852	2	0
397	436	853	2	0
398	432	848	3	0
394	432	852	0	0
403	442	850	2	0
33	441	854	2	0
99	428	849	2	0
394	446	853	0	0
399	442	853	2	0
405	439	850	3	0
397	437	855	2	0
395	425	852	6	0
32	427	854	7	0
32	433	854	2	0
33	432	850	6	0
33	425	848	3	0
33	434	851	1	0
399	444	851	2	0
395	446	849	0	0
33	432	854	6	0
397	447	851	0	0
32	444	854	2	0
2	418	851	0	1
399	454	841	1	0
397	451	843	1	0
33	450	840	1	0
395	448	842	1	0
396	452	841	1	0
404	452	853	6	0
394	449	849	0	0
395	450	853	0	0
398	450	851	0	0
395	454	854	0	0
404	454	849	5	0
396	448	853	6	0
399	453	851	0	0
405	450	855	0	0
32	455	852	0	0
32	452	849	5	0
33	457	842	1	0
398	461	854	0	0
32	456	840	1	0
397	456	855	4	0
398	457	853	5	0
394	463	854	0	0
33	456	849	6	0
401	460	849	0	0
404	459	850	5	0
405	459	842	1	0
396	458	855	2	0
402	458	849	0	0
33	462	852	0	0
780	471	850	0	0
21	470	855	4	0
21	470	853	4	0
781	466	850	0	0
404	469	855	5	0
401	438	860	2	0
399	433	859	2	0
33	434	862	2	0
396	437	857	2	0
397	446	856	0	0
399	447	858	0	0
395	443	860	2	0
32	445	861	2	0
398	434	856	2	0
33	432	858	6	0
394	432	861	2	0
402	436	861	2	0
33	443	858	2	0
404	435	858	2	0
403	439	862	2	0
397	439	857	2	0
32	444	856	2	0
395	449	862	0	0
397	445	858	0	0
396	447	861	0	0
394	432	856	2	0
32	444	862	2	0
395	442	856	2	0
399	471	860	5	0
396	455	857	0	0
394	465	856	5	0
399	462	856	0	0
403	470	859	5	0
405	441	861	2	0
395	453	859	0	0
395	463	859	0	0
405	470	861	5	0
32	461	859	0	0
396	469	858	5	0
32	454	861	2	0
405	463	861	0	0
395	448	856	0	0
395	465	861	5	0
397	451	858	0	0
32	467	861	5	0
394	458	861	0	0
394	452	861	0	0
403	467	859	5	0
397	460	856	3	0
33	456	860	0	0
398	449	859	0	0
32	440	859	2	0
404	458	858	1	0
33	467	856	5	0
1100	439	868	5	0
1091	438	866	7	0
1099	438	868	6	0
1099	438	870	4	0
1086	437	866	0	0
1099	435	866	1	0
1092	441	868	4	0
1100	434	866	2	0
164	443	870	2	0
1091	434	870	4	0
1092	433	868	1	0
400	463	869	1	0
1176	452	870	0	0
1092	432	866	3	0
164	442	870	4	0
1100	440	870	4	0
1100	461	866	7	0
1099	464	868	5	0
1099	463	870	3	0
1092	436	870	4	0
164	459	869	0	0
395	457	871	0	0
399	434	864	1	0
164	460	868	7	0
1086	432	870	4	0
398	440	864	1	0
1086	436	868	7	0
1086	462	868	6	0
1091	435	868	0	0
1099	441	866	5	0
164	443	869	3	0
1100	461	870	4	0
32	446	870	1	0
397	455	864	1	0
1099	463	866	0	0
397	463	864	1	0
1092	440	866	6	0
398	468	871	2	0
1086	465	870	2	0
1092	465	866	1	0
396	471	864	1	0
204	433	866	3	1
204	460	870	3	1
204	439	870	2	1
204	437	871	0	1
204	439	866	2	1
204	432	868	3	1
204	467	869	2	1
204	442	868	2	1
204	461	868	3	1
204	462	867	0	1
204	442	866	2	1
204	435	870	2	1
204	464	867	0	1
204	462	870	3	1
204	440	868	0	1
204	437	868	0	1
204	434	868	0	1
204	464	870	2	1
204	466	870	2	1
204	436	866	3	1
204	465	868	2	1
204	433	870	3	1
204	463	868	3	1
204	441	870	3	1
395	427	859	2	0
394	428	858	0	0
396	426	860	0	0
99	429	857	2	0
395	431	858	6	0
98	425	860	2	0
99	429	858	2	0
32	425	859	2	0
1100	430	870	6	0
1086	431	865	4	0
1092	431	866	0	0
1100	429	868	0	0
164	427	868	1	0
164	427	869	3	0
1099	428	865	0	0
1092	430	867	0	0
99	426	859	2	0
32	430	861	2	0
395	429	860	6	0
33	428	859	2	0
397	431	860	6	0
403	425	856	1	0
1099	431	868	0	0
32	424	860	1	0
1099	429	866	0	0
1091	430	865	0	0
1092	428	867	0	0
1099	428	870	5	0
164	427	867	2	0
1176	425	870	0	0
204	430	868	2	1
204	430	866	2	1
204	429	867	3	1
204	431	867	2	1
204	431	866	0	1
204	428	868	3	1
204	428	866	3	1
204	429	865	3	1
33	423	858	0	0
395	423	860	2	0
32	422	861	4	0
395	422	858	0	0
401	422	860	3	0
33	418	858	3	0
394	421	861	5	0
33	419	860	7	0
32	417	860	2	0
403	417	861	1	0
32	420	860	7	0
397	420	861	6	0
396	419	861	1	0
32	420	856	0	0
395	418	861	0	0
394	416	861	0	0
395	419	864	0	0
1176	417	870	0	0
397	423	875	7	0
32	450	873	6	0
398	451	878	1	0
1151	452	872	1	0
399	441	878	0	0
403	445	879	6	0
1151	451	873	7	0
404	448	877	7	0
32	416	873	2	0
405	442	873	3	0
404	437	877	2	0
395	419	873	3	0
397	434	878	5	0
397	448	872	1	0
396	453	876	2	0
1151	453	873	0	0
1113	424	878	0	0
395	445	875	0	0
32	431	879	6	0
399	427	875	1	0
33	440	875	7	0
394	439	872	2	0
430	436	875	0	0
404	430	872	5	0
398	428	878	2	0
32	454	873	5	0
395	425	874	0	0
398	432	876	4	0
395	433	872	1	0
394	430	876	4	0
403	427	872	0	0
32	436	873	0	0
403	434	875	3	0
396	438	875	6	0
397	417	876	2	0
32	423	872	1	0
32	417	879	2	0
33	420	875	6	0
394	420	878	5	0
399	439	878	1	0
400	459	872	7	0
404	460	879	7	0
396	463	877	7	0
394	461	875	7	0
404	463	873	7	0
403	457	876	7	0
394	465	874	2	0
394	465	872	1	0
404	471	872	0	0
1113	468	878	2	0
32	467	873	3	0
395	471	874	7	0
21	472	853	0	0
21	472	855	0	0
32	473	840	5	0
32	476	858	5	0
33	475	861	5	0
399	474	859	5	0
394	473	861	5	0
1099	476	866	1	0
1086	476	867	1	0
401	475	879	1	0
1099	476	868	1	0
1100	477	868	1	0
398	473	858	5	0
32	472	878	2	0
400	477	867	1	0
405	474	872	1	0
1092	477	866	1	0
204	475	869	3	1
33	475	874	3	0
204	476	870	0	1
394	461	885	7	0
396	452	883	2	0
430	454	882	0	0
398	460	882	7	0
33	468	887	0	0
397	464	883	7	0
399	475	884	0	0
394	449	887	2	0
396	468	883	0	0
399	458	884	7	0
396	458	886	7	0
404	471	885	0	0
395	456	880	7	0
405	453	886	5	0
32	474	882	0	0
1176	453	880	0	0
395	470	883	0	0
394	473	887	0	0
33	456	883	7	0
394	449	880	1	0
397	462	880	7	0
32	445	887	2	0
398	447	884	2	0
1176	444	882	0	0
1113	442	887	2	0
403	440	881	7	0
32	438	886	4	0
404	437	883	2	0
405	432	883	4	0
396	433	886	4	0
401	434	881	3	0
398	437	880	6	0
33	415	876	2	0
398	413	874	2	0
1113	411	878	2	0
394	408	877	2	0
395	410	875	2	0
33	405	876	1	0
404	407	874	2	0
396	405	873	2	0
32	403	879	2	0
395	402	873	0	0
398	401	876	2	0
1099	401	872	7	0
1086	400	872	6	0
405	429	881	7	0
405	419	881	4	0
397	430	885	7	0
401	427	884	0	0
400	423	885	2	0
404	430	883	2	0
32	424	883	0	0
395	426	886	1	0
32	418	884	2	0
33	422	881	5	0
395	420	884	3	0
397	425	881	6	0
404	414	3573	4	0
399	423	3587	0	0
397	413	3582	0	0
404	414	3576	5	0
34	409	3578	0	0
403	413	3571	0	0
401	417	3577	5	0
398	418	3577	5	0
398	416	3582	0	0
399	416	3580	0	0
403	415	3577	5	0
396	427	3573	0	0
5	426	3572	0	0
399	426	3576	7	0
396	413	3584	0	0
403	422	3579	5	0
404	426	3573	0	0
401	411	3571	0	0
401	412	3575	0	0
401	423	3574	7	0
396	413	3579	0	0
404	419	3574	7	0
398	421	3587	0	0
396	419	3588	0	0
397	420	3589	0	0
401	423	3585	0	0
402	424	3583	0	0
404	424	3579	5	0
396	424	3571	7	0
403	425	3573	0	0
151	410	3576	2	1
151	411	3575	2	1
151	409	3577	2	1
1	198	668	0	0
0	192	665	2	0
1	198	658	0	0
0	194	660	0	0
0	196	657	0	0
0	199	669	0	0
0	193	667	0	0
0	199	666	0	0
1	195	659	0	0
0	198	663	0	0
0	193	662	2	0
34	195	667	1	0
0	194	664	3	0
1	197	661	0	0
1	199	659	0	0
1	193	657	0	0
0	196	668	0	0
0	196	665	0	0
1	198	665	0	0
34	153	694	4	0
70	159	688	0	0
34	153	692	4	0
34	152	694	4	0
70	161	693	0	0
34	152	695	4	0
36	164	688	4	0
34	139	695	0	0
36	136	692	0	0
36	136	691	0	0
36	137	693	0	0
36	138	690	0	0
34	137	691	0	0
34	137	694	0	0
36	135	691	0	0
34	135	695	0	0
0	120	702	0	0
36	127	699	6	0
36	127	696	6	0
34	120	697	6	0
36	134	703	6	0
36	135	701	6	0
34	134	696	0	0
34	135	696	0	0
36	134	701	6	0
36	133	702	6	0
34	137	702	0	0
34	138	702	0	0
4	121	704	0	0
36	153	697	4	0
36	152	698	4	0
36	153	696	4	0
36	152	697	4	0
0	204	656	0	0
34	206	656	0	0
34	204	660	1	0
0	200	662	3	0
0	201	663	0	0
0	206	661	0	0
0	201	657	3	0
0	206	658	0	0
0	201	660	2	0
34	203	665	1	0
1	202	665	0	0
34	208	632	0	0
34	210	636	0	0
1	210	639	0	0
45	214	670	6	0
0	212	643	0	0
34	209	633	0	0
1	208	645	0	0
34	215	649	0	0
21	215	671	0	0
0	208	656	0	0
1	213	659	0	0
34	208	662	0	0
45	212	667	2	0
45	212	669	2	0
45	212	668	2	0
34	211	645	0	0
45	212	670	2	0
1	213	636	0	0
34	209	641	0	0
1	212	638	0	0
34	212	641	0	0
1	212	634	0	0
1	214	652	0	0
34	213	655	0	0
45	212	671	2	0
34	208	647	0	0
45	214	671	6	0
45	214	669	6	0
45	214	667	6	0
1	213	663	0	0
45	214	668	6	0
21	211	670	4	0
45	214	679	4	0
21	215	675	0	0
45	214	675	6	0
45	214	676	6	0
45	212	672	2	0
45	214	672	6	0
45	212	675	2	0
45	212	674	2	0
45	212	673	2	0
45	212	679	2	0
45	212	677	2	0
45	212	676	2	0
45	215	679	4	0
45	213	679	4	0
45	215	677	0	0
45	214	674	6	0
21	211	679	4	0
21	211	674	4	0
45	214	673	6	0
45	212	678	2	0
21	219	677	0	0
45	218	679	6	0
45	216	677	0	0
45	218	678	6	0
34	223	647	0	0
34	223	652	0	0
193	221	664	4	0
21	216	676	4	0
45	218	677	6	0
45	217	677	0	0
193	224	661	4	0
193	224	659	4	0
34	227	651	0	0
34	227	646	0	0
34	230	645	0	0
139	216	633	2	0
3	216	638	2	0
70	228	636	0	0
63	219	633	2	0
3	223	638	2	0
30	226	633	4	0
70	225	635	0	0
70	214	629	2	0
55	214	627	0	0
0	209	631	0	0
1	214	625	1	1
6	216	1564	0	0
15	214	1564	0	0
6	218	1574	2	0
15	221	1574	2	0
141	216	1562	6	0
21	214	680	0	0
0	215	686	0	0
0	215	684	0	0
4	219	685	0	0
4	219	686	5	0
21	219	681	0	0
45	218	683	6	0
45	216	680	2	0
45	216	682	2	0
45	216	681	2	0
45	218	682	6	0
45	218	680	6	0
45	216	683	2	0
45	218	681	6	0
3	214	695	0	0
3	213	695	0	0
27	212	694	0	0
7	211	693	0	0
3	211	692	0	0
5	218	692	0	0
0	221	690	6	0
6	218	694	4	0
47	210	694	0	0
3	217	695	0	0
25	219	695	4	0
47	216	694	4	0
26	219	688	0	0
9	215	689	0	1
1	213	691	3	1
1	215	692	2	1
1	217	690	2	1
47	215	697	0	0
47	214	697	4	0
25	217	697	4	0
3	214	696	0	0
3	213	696	0	0
3	213	1633	0	0
20	215	1642	0	0
25	214	1634	0	0
3	213	1634	0	0
25	214	1636	0	0
6	218	1636	0	0
25	216	1642	0	0
20	213	1642	0	0
15	218	1638	0	0
20	215	1640	0	0
14	210	1636	0	0
3	210	1639	0	0
7	212	1634	2	0
3	213	1635	0	0
17	218	1640	0	0
5	216	1636	0	0
20	213	1640	0	0
25	216	1639	0	0
1	214	1639	0	1
1	216	1638	1	1
1	214	1637	1	1
1	216	1635	1	1
25	213	2583	0	0
51	215	2582	0	0
15	214	2583	6	0
25	211	2583	0	0
51	215	2585	4	0
6	216	2580	0	0
47	218	2582	4	0
1	216	2581	1	1
1	216	2579	1	1
1	217	2582	0	1
25	213	3519	0	0
25	213	3524	0	0
15	213	3520	0	0
20	223	3514	0	0
14	223	3516	0	0
51	220	3519	6	0
25	216	3519	0	0
20	223	3520	0	0
23	216	3522	0	0
3	218	3520	0	0
25	218	3522	0	0
25	218	3519	0	0
27	217	3525	0	0
20	221	3517	0	0
20	226	3517	0	0
25	218	3524	0	0
5	218	3526	0	0
20	221	3520	0	0
51	220	3521	6	0
20	226	3514	0	0
20	226	3520	0	0
1	215	3525	0	1
1	221	3514	1	1
283	582	625	4	0
284	582	624	2	0
0	571	625	0	0
286	568	630	0	0
286	571	637	0	0
35	580	623	2	0
284	582	629	2	0
284	580	629	2	0
284	579	630	0	0
20	581	631	4	0
70	571	618	2	0
0	589	625	0	0
0	589	626	0	0
283	580	630	0	0
23	580	631	6	0
283	580	628	2	0
79	569	605	0	0
283	582	628	2	0
344	569	607	0	0
362	580	638	0	0
1	571	630	0	0
70	571	615	2	0
365	581	639	2	0
284	578	625	4	0
71	580	633	0	0
34	591	617	2	0
283	584	628	2	0
0	591	626	0	0
35	577	619	2	0
283	582	630	0	0
363	582	638	4	0
0	584	625	0	0
6	578	637	2	0
0	590	636	4	0
7	583	634	6	0
364	581	637	6	0
283	578	630	0	0
20	587	612	0	0
1	584	624	0	0
0	591	628	4	0
58	598	603	4	0
6	578	626	0	0
283	580	625	2	0
283	604	612	4	0
0	590	633	4	0
284	581	630	0	0
0	593	625	4	0
284	577	626	2	0
284	580	624	2	0
34	592	620	2	0
286	596	605	4	0
1	586	626	0	0
6	581	624	0	0
6	581	629	0	0
284	579	625	4	0
6	579	629	0	0
1	605	631	0	0
6	583	629	0	0
6	579	624	0	0
45	593	604	4	0
45	592	602	0	0
1	594	638	4	0
3	584	634	4	0
7	584	633	4	0
0	604	638	3	0
1	607	634	0	0
20	603	610	0	0
22	583	639	0	0
51	579	636	0	0
20	589	611	2	0
45	593	602	0	0
0	585	628	0	0
366	583	637	6	0
1	588	629	4	0
0	577	633	0	0
5	578	636	2	0
45	594	602	0	0
283	577	625	2	0
0	590	624	0	0
286	596	606	4	0
284	584	629	2	0
283	577	627	4	0
45	594	604	4	0
0	594	624	0	0
45	592	604	4	0
63	607	603	4	0
34	596	622	2	0
45	595	604	4	0
45	595	602	0	0
2	600	633	0	0
283	599	609	4	0
34	594	618	2	0
1	607	625	0	0
1	585	632	1	1
1	581	635	0	1
1	603	601	1	1
1	582	603	0	1
1	581	636	0	1
1	603	606	1	1
281	563	600	0	0
58	567	607	2	0
330	560	612	2	0
1	560	602	0	0
1	565	604	0	1
323	566	594	0	0
46	567	592	0	0
1	562	597	0	0
46	564	592	0	0
0	565	598	0	0
63	587	592	0	0
23	585	594	2	0
23	583	594	2	0
23	581	594	2	0
46	562	592	0	0
1	573	599	0	0
1	572	595	0	0
1	590	595	0	1
328	558	617	0	0
46	559	592	0	0
281	554	600	6	0
281	554	605	6	0
325	555	593	0	0
329	559	617	0	0
27	556	600	0	0
290	557	615	0	0
278	557	605	6	0
290	557	614	0	0
29	553	611	6	0
1	555	617	0	1
1	559	613	0	1
1	554	602	1	1
327	551	599	0	0
63	550	612	4	0
139	550	611	4	0
322	544	599	0	0
285	557	591	4	0
142	558	587	4	0
25	562	586	0	0
25	562	591	0	0
331	563	587	0	0
6	560	591	0	0
500	581	587	2	0
285	558	591	4	0
3	583	584	2	0
3	583	587	2	0
25	559	586	0	0
23	570	591	0	0
285	571	586	0	0
26	568	589	6	0
23	581	589	2	0
281	576	584	2	0
285	569	588	0	0
19	577	591	2	0
281	576	586	2	0
326	544	590	0	0
25	559	591	0	0
58	556	586	2	0
26	571	589	6	0
285	555	591	4	0
37	571	588	0	0
286	568	586	0	0
328	560	588	0	0
285	556	591	4	0
37	568	587	0	0
285	568	584	0	0
47	581	584	6	0
286	571	587	0	0
7	579	587	2	0
48	579	584	6	0
3	578	587	6	0
23	583	589	2	0
98	561	586	0	1
94	565	586	0	1
1	578	585	1	1
1	581	586	1	1
273	536	601	0	0
45	536	614	0	0
45	537	614	0	0
21	540	613	2	0
5	542	599	0	0
45	538	614	0	0
21	537	613	2	0
45	540	616	4	0
45	542	614	0	0
45	542	616	4	0
45	539	614	0	0
45	541	614	0	0
45	541	616	4	0
157	536	617	4	0
45	540	614	0	0
93	538	592	0	1
94	539	599	0	1
55	534	616	4	0
45	534	614	0	0
20	529	614	6	0
55	532	616	4	0
21	534	613	2	0
155	531	617	4	0
283	533	605	2	0
45	533	614	0	0
20	533	599	0	0
20	533	604	2	0
45	535	614	0	0
20	533	594	0	0
45	532	614	0	0
45	530	616	4	0
45	530	614	0	0
21	531	613	2	0
20	529	616	4	0
45	531	614	0	0
46	525	604	2	0
34	521	619	6	0
34	522	618	6	0
34	527	619	6	0
36	522	621	6	0
0	526	608	0	0
34	527	622	6	0
46	520	609	2	0
299	522	603	0	0
283	522	608	1	0
283	524	608	1	0
34	524	621	6	0
36	521	616	6	0
45	524	607	2	0
46	525	609	2	0
286	520	602	0	0
285	525	603	0	0
36	525	619	6	0
36	524	616	6	0
45	523	608	0	0
286	523	611	0	0
286	524	613	0	0
299	522	610	0	0
45	522	606	6	0
286	524	618	0	0
45	522	607	6	0
6	523	607	4	0
0	521	596	0	0
45	524	606	2	0
0	526	600	0	0
299	526	606	2	0
1	523	598	0	0
46	520	604	2	0
285	525	595	0	0
1	522	592	0	0
1	518	613	2	0
0	518	592	0	0
283	514	611	2	0
283	517	610	2	0
299	519	606	2	0
454	513	614	0	0
454	514	636	1	0
0	509	584	0	0
1	508	587	0	0
286	516	586	0	0
0	505	585	0	0
334	540	591	0	0
1	510	590	0	0
286	522	585	0	0
0	513	587	0	0
286	523	590	0	0
285	516	589	0	0
1	516	590	0	0
51	508	3444	2	0
51	508	3426	2	0
51	509	3432	2	0
51	508	3429	2	0
51	510	3445	4	0
207	505	3430	7	0
51	509	3417	2	0
317	508	3443	6	0
51	508	3439	2	0
315	504	3424	0	0
51	509	3423	2	0
51	513	3417	6	0
207	505	3425	6	0
51	515	3429	4	0
51	512	3423	6	0
51	517	3427	0	0
51	517	3429	4	0
46	524	3432	6	0
5	523	3439	0	0
46	524	3436	6	0
46	520	3436	6	0
51	515	3427	0	0
46	520	3432	6	0
51	513	3439	6	0
51	513	3445	4	0
51	512	3432	6	0
90	512	3441	0	1
91	508	3427	1	1
88	509	3441	0	1
51	513	3412	6	0
318	511	3410	4	0
51	508	3412	2	0
316	511	3415	4	0
51	509	3410	0	0
51	512	3410	0	0
92	510	3415	0	1
207	503	3427	1	0
315	501	3427	6	0
315	502	3431	0	0
51	501	3429	2	0
51	503	3431	4	0
207	503	3428	0	0
51	503	3424	0	0
51	501	3426	2	0
1	505	579	0	0
1	505	581	0	0
286	511	580	0	0
286	525	580	0	0
286	519	579	0	0
306	512	582	0	0
273	542	577	6	0
193	496	613	0	0
286	501	578	0	0
307	497	580	0	0
21	498	613	4	0
307	497	584	0	0
193	491	613	0	0
193	494	613	0	0
193	489	608	0	0
21	492	613	4	0
307	494	577	0	0
193	491	611	0	0
99	488	611	3	0
193	488	608	0	0
193	491	610	0	0
21	495	613	4	0
1	500	572	0	0
0	495	574	0	0
306	495	569	0	0
308	503	571	0	0
1	489	569	0	0
110	525	572	0	0
111	519	568	0	0
111	519	569	0	0
111	519	570	0	0
110	515	570	0	0
309	509	571	0	0
110	514	574	0	0
110	525	569	0	0
283	499	570	0	0
103	523	570	0	0
103	517	573	0	0
103	517	572	0	0
102	517	574	0	0
110	524	572	0	0
110	523	574	0	0
110	518	572	0	0
20	541	571	0	0
20	540	575	0	0
20	541	573	0	0
15	551	568	4	0
7	551	571	4	0
3	551	572	4	0
45	546	577	2	0
45	544	577	6	0
324	551	583	0	0
20	544	568	0	0
273	548	577	6	0
5	548	580	6	0
1	545	581	0	1
1	541	565	0	0
1	545	562	0	0
34	522	561	0	0
34	546	566	0	0
0	515	561	0	0
72	559	563	3	0
72	559	565	3	0
72	559	561	3	0
72	559	567	3	0
72	558	567	3	0
72	555	565	3	0
72	555	563	3	0
72	556	563	3	0
72	554	563	3	0
72	552	561	3	0
72	553	567	3	0
72	554	561	3	0
72	553	565	3	0
72	556	561	3	0
15	554	568	4	0
306	558	578	0	0
72	552	563	3	0
72	557	563	3	0
72	557	561	3	0
72	558	565	3	0
46	559	581	0	0
72	556	567	3	0
72	557	567	3	0
72	558	561	3	0
72	556	565	3	0
72	558	563	3	0
72	552	567	3	0
72	553	563	3	0
72	552	565	3	0
72	554	565	3	0
72	555	561	3	0
72	553	561	3	0
72	557	565	3	0
72	554	567	3	0
283	555	582	4	0
58	554	578	4	0
283	558	582	4	0
1	553	573	0	1
309	507	567	0	0
306	508	560	4	0
1	508	564	0	0
0	506	561	0	0
21	500	562	6	0
1	502	567	0	0
5	501	564	0	0
37	500	567	0	0
666	497	563	0	0
283	496	561	0	0
21	496	562	0	0
20	499	560	0	0
21	498	562	0	0
21	496	564	2	0
666	499	563	0	0
677	495	563	2	0
676	488	563	2	0
20	488	560	0	0
676	489	563	2	0
676	490	563	2	0
1	491	561	0	0
530	491	563	2	0
530	492	563	2	0
20	493	560	0	0
676	487	563	2	0
306	512	558	0	0
306	518	556	0	0
283	525	553	0	0
0	522	556	0	0
34	539	557	0	0
59	535	555	4	0
1	540	552	0	0
72	558	559	3	0
72	559	555	3	0
72	558	557	3	0
72	559	557	3	0
72	559	559	3	0
72	555	557	3	0
72	556	555	3	0
72	558	555	3	0
72	557	559	3	0
72	555	555	3	0
72	554	557	3	0
72	552	555	3	0
72	554	555	3	0
51	558	554	4	0
72	552	559	3	0
72	553	555	3	0
72	556	559	3	0
72	557	557	3	0
72	557	555	3	0
72	553	559	3	0
72	554	559	3	0
72	553	557	3	0
72	556	557	3	0
72	555	559	3	0
72	552	557	3	0
24	518	548	6	0
0	514	549	2	0
1177	517	546	2	0
0	522	548	0	0
1079	512	550	2	0
7	518	547	2	0
0	532	544	6	0
1	515	545	1	1
1	543	549	0	0
34	548	551	0	0
34	547	548	0	0
3	517	543	0	0
0	535	539	6	0
145	517	536	2	0
1080	512	540	6	0
1	521	536	3	0
47	518	543	4	0
34	524	541	0	0
308	544	536	6	0
0	539	540	6	0
1170	512	536	2	0
42	515	538	6	0
145	514	539	0	0
0	543	541	6	0
1149	515	543	6	0
0	547	537	6	0
308	536	543	6	0
10	517	534	5	0
280	515	535	2	0
34	519	531	0	0
145	515	533	4	0
145	513	533	4	0
678	499	1506	4	0
6	501	1508	0	0
529	500	1506	4	0
665	497	1507	0	0
678	498	1506	4	0
665	499	1507	0	0
309	513	525	0	0
0	516	522	0	0
34	526	530	0	0
37	524	531	0	0
0	529	535	6	0
0	545	533	6	0
1	555	539	0	0
1	553	546	0	0
60	558	542	0	0
15	557	550	4	0
3	559	534	0	0
15	561	533	0	0
55	563	533	0	0
191	566	554	2	0
72	560	559	3	0
74	565	537	6	0
306	566	548	0	0
72	560	555	3	0
72	560	557	3	0
5	566	530	0	0
191	566	545	2	0
55	563	532	0	0
2	560	545	4	0
72	560	567	3	0
53	565	532	0	0
3	561	549	4	0
59	564	552	4	0
72	560	561	3	0
191	566	551	2	0
51	560	554	4	0
72	560	565	3	0
72	560	563	3	0
0	554	532	6	0
1	562	551	1	1
1	560	532	0	1
1	566	537	0	1
20	566	568	0	0
20	562	568	0	0
284	561	577	6	0
46	564	581	0	0
285	567	580	0	0
308	563	578	0	0
37	566	581	0	0
284	561	578	6	0
283	561	581	6	0
46	567	581	0	0
46	562	581	0	0
284	561	580	6	0
284	561	579	6	0
285	560	581	6	0
37	560	577	0	0
285	560	578	6	0
1	563	582	0	1
34	525	523	0	0
46	533	525	0	0
46	530	525	0	0
46	530	522	0	0
46	536	525	0	0
46	536	522	0	0
309	554	524	4	0
306	523	515	0	0
46	536	513	0	0
46	533	513	0	0
46	530	519	0	0
45	532	516	6	0
46	530	513	0	0
45	533	515	4	0
6	533	516	4	0
46	536	519	0	0
45	534	516	2	0
46	536	516	0	0
46	530	516	0	0
1	547	513	0	0
55	569	533	0	0
55	569	532	0	0
191	568	548	2	0
191	570	551	2	0
191	570	545	2	0
191	572	551	2	0
191	572	548	2	0
54	572	528	0	0
191	568	551	2	0
191	570	548	2	0
191	572	545	2	0
191	568	545	2	0
50	565	508	7	0
379	565	504	0	0
1	550	508	0	0
1	548	506	0	0
1	543	509	0	0
1	563	505	1	1
1	563	510	3	1
306	522	509	3	0
1	515	506	4	0
1	514	514	4	0
5	566	502	4	0
353	572	500	4	0
354	572	503	4	0
11	566	498	4	0
280	563	500	4	0
2	562	498	0	1
143	530	3339	2	0
217	533	3338	6	0
51	532	3334	2	0
51	532	3342	2	0
51	530	3331	4	0
143	530	3337	2	0
51	540	3331	4	0
51	544	3352	4	0
51	531	3346	2	0
51	537	3347	0	0
51	540	3352	4	0
20	541	3329	6	0
51	533	3353	4	0
51	535	3346	6	0
51	535	3351	6	0
51	536	3339	6	0
51	542	3346	0	0
370	537	3337	6	0
5	533	3348	0	0
51	536	3331	4	0
51	531	3351	2	0
104	533	3342	0	1
106	536	3349	1	1
104	533	3335	0	1
51	554	3355	0	0
51	554	3357	4	0
369	534	3371	6	0
24	543	3373	2	1
24	542	3374	2	1
51	521	3371	6	0
117	521	3374	6	0
51	520	3368	0	0
51	533	3321	6	0
51	540	3323	0	0
51	536	3323	0	0
51	530	3324	0	0
51	531	3321	2	0
20	541	3326	4	0
21	590	501	4	0
261	588	496	4	0
45	591	503	4	0
45	586	502	2	0
261	588	498	4	0
45	586	503	2	0
376	589	501	4	0
21	585	499	4	0
261	585	498	4	0
261	585	496	4	0
45	590	503	4	0
45	589	503	4	0
21	585	503	4	0
45	586	501	2	0
45	588	503	4	0
21	585	502	6	0
45	586	500	2	0
376	588	500	0	0
21	588	499	0	0
6	583	1403	0	0
0	582	518	2	0
20	581	519	0	0
0	583	518	4	0
0	579	519	0	0
1	581	525	0	0
20	581	523	2	0
11	583	520	0	0
0	577	520	0	0
0	580	523	0	0
0	578	526	0	0
982	576	523	5	0
1	582	532	2	0
34	583	535	6	0
1	580	530	2	0
3	588	519	0	0
45	587	514	6	0
45	587	513	6	0
194	588	505	4	0
21	588	509	0	0
5	584	523	4	0
45	587	512	6	0
1	587	531	2	0
194	585	505	4	0
0	588	526	0	0
21	585	509	4	0
21	590	504	0	0
285	591	517	7	0
45	586	504	2	0
3	585	519	0	0
0	591	524	0	0
1	591	533	2	0
1	586	519	0	1
112	586	524	0	1
1	583	537	2	0
1	577	542	2	0
443	588	540	4	0
1	579	536	2	0
22	594	543	0	0
3	599	519	7	0
285	593	517	7	0
285	592	515	7	0
1	598	533	6	0
1	596	525	0	0
285	593	521	7	0
285	594	514	7	0
34	598	535	6	0
0	599	526	0	0
21	594	504	0	0
1	599	542	6	0
285	599	508	7	0
285	598	511	7	0
34	593	534	6	0
285	596	512	7	0
22	595	543	0	0
217	596	540	1	0
1	598	516	2	1
0	605	521	0	0
0	601	521	0	0
285	602	509	7	0
1	603	524	0	0
285	603	511	7	0
0	607	505	0	0
1	607	518	0	0
3	601	517	7	0
0	607	543	6	0
34	603	538	6	0
1164	606	516	6	0
34	605	533	6	0
1	603	506	3	1
6	581	1465	6	0
6	584	1467	4	0
45	598	503	4	0
45	599	503	4	0
45	597	503	4	0
194	596	501	4	0
0	607	497	0	0
261	593	501	4	0
21	594	501	4	0
0	612	501	0	0
0	614	519	0	0
0	614	497	0	0
0	608	503	0	0
0	609	501	0	0
0	613	509	0	0
0	615	514	0	0
0	608	498	0	0
0	614	512	0	0
0	613	525	0	0
0	613	520	0	0
0	611	521	0	0
1	615	517	0	0
1	611	523	0	0
0	610	512	0	0
0	608	508	0	0
0	609	515	0	0
0	610	519	0	0
1164	610	499	6	0
0	611	506	0	0
1164	609	506	6	0
205	613	499	0	0
1	611	510	0	0
205	613	516	0	0
0	612	503	0	0
1	611	534	6	0
1	610	496	0	0
0	608	523	0	0
205	608	525	0	0
1164	615	495	6	0
0	614	493	0	0
205	611	492	0	0
1	612	489	0	0
0	610	494	0	0
0	609	490	0	0
0	622	500	0	0
0	622	490	0	0
0	618	507	0	0
0	616	501	0	0
0	621	492	0	0
981	616	492	0	0
0	621	497	0	0
0	617	504	0	0
0	622	496	0	0
0	616	495	0	0
1	616	507	0	0
205	619	501	0	0
0	623	518	0	0
0	621	503	0	0
0	616	521	0	0
0	619	520	0	0
1	620	489	0	0
0	619	495	0	0
205	620	519	0	0
205	619	491	0	0
0	619	517	0	0
1	617	524	0	0
0	620	524	0	0
0	622	525	0	0
980	617	490	6	0
0	617	512	0	0
1164	617	493	6	0
0	622	515	0	0
0	620	505	0	0
0	617	497	0	0
1	620	510	0	0
0	618	498	0	0
205	615	483	0	0
0	621	485	0	0
0	619	487	0	0
1	619	483	0	0
0	611	487	0	0
0	617	485	0	0
1	578	551	0	0
4	581	551	0	0
1	579	547	2	0
55	597	545	6	0
4	582	549	0	0
34	593	550	6	0
1	583	544	6	0
0	603	545	6	0
444	593	547	2	0
0	596	550	6	0
55	596	548	0	0
1	593	545	1	1
191	570	554	2	0
191	568	554	2	0
191	572	554	2	0
4	583	559	6	0
4	583	556	0	0
1	580	552	0	0
4	585	559	6	0
4	579	556	0	0
4	578	557	0	0
1	577	554	0	0
1	593	556	6	0
1	606	558	6	0
0	606	554	6	0
4	580	554	0	0
4	581	559	6	0
4	582	558	6	0
34	584	557	6	0
406	585	556	1	0
1	587	554	0	0
4	589	553	0	0
4	580	563	6	0
4	579	561	6	0
36	600	564	6	0
0	593	560	6	0
1	604	562	6	0
0	593	563	6	0
15	584	561	4	0
4	584	567	6	0
4	577	560	0	0
7	582	564	0	0
7	582	562	4	0
34	590	561	6	0
8	583	566	0	0
4	580	566	6	0
8	582	566	7	0
3	582	563	4	0
36	596	564	6	0
34	581	561	6	0
36	599	566	6	0
1	586	563	1	1
34	590	572	6	0
139	583	571	2	0
20	579	572	6	0
20	582	572	0	0
20	589	568	0	0
36	599	571	6	0
36	598	569	6	0
36	595	569	6	0
3	603	574	2	0
20	586	568	6	0
36	595	573	6	0
48	605	575	2	0
37	607	571	0	0
15	603	571	4	0
48	581	581	0	0
278	607	583	4	0
11	581	578	0	0
5	585	578	0	0
306	569	578	0	0
3	606	579	0	0
274	607	579	4	0
286	569	582	0	0
281	607	581	0	0
7	607	582	0	0
37	568	582	0	0
94	586	581	1	1
94	581	580	1	1
34	612	564	6	0
34	611	538	6	0
0	612	539	6	0
1	613	557	6	0
0	609	554	6	0
37	608	572	0	0
37	608	571	0	0
37	608	574	0	0
20	608	568	6	0
3	611	568	0	0
1	609	562	6	0
20	610	568	0	0
37	608	575	0	0
1	612	573	0	1
152	608	573	1	1
118	591	590	4	0
7	584	584	2	0
23	585	589	2	0
71	584	587	2	0
1	589	591	1	1
1	586	586	1	1
178	604	590	0	0
20	603	597	0	0
286	600	594	4	0
335	582	1527	0	0
6	585	1522	0	0
15	606	1525	6	0
334	584	1527	0	0
15	581	1522	4	0
51	558	3421	0	0
51	557	3424	2	0
51	563	3424	6	0
80	567	3434	4	0
97	580	3418	0	0
278	578	3418	5	0
51	562	3418	0	0
278	581	3416	5	0
278	579	3416	5	0
434	580	3421	0	0
51	552	3426	4	0
51	559	3437	4	0
51	554	3426	4	0
80	564	3434	4	0
51	561	3437	4	0
278	582	3418	5	0
5	569	3437	6	0
80	561	3430	2	0
95	557	3425	1	1
51	565	3420	6	0
103	549	3429	0	0
103	549	3430	0	0
80	550	3423	4	0
111	549	3431	0	0
80	544	3430	6	0
51	547	3436	4	0
111	550	3432	0	0
51	550	3433	6	0
111	548	3436	0	0
111	545	3435	0	0
103	550	3434	0	0
51	544	3432	2	0
80	544	3428	6	0
111	545	3434	0	0
80	577	3409	4	0
290	577	3411	0	0
80	577	3414	6	0
80	583	3408	4	0
51	580	3410	0	0
290	578	3412	5	0
335	539	1547	0	0
6	542	1543	0	0
336	542	1546	0	0
6	548	1524	0	0
14	544	1520	4	0
145	559	1532	6	0
290	559	1530	0	0
290	559	1533	0	0
49	559	1534	6	0
145	559	1535	6	0
34	619	532	6	0
0	620	539	6	0
34	621	537	6	0
1	620	534	6	0
192	619	543	0	0
34	619	557	6	0
1	621	559	6	0
6	618	551	6	0
192	616	543	0	0
34	616	541	6	0
96	617	556	0	1
0	631	551	0	0
0	626	553	0	0
0	631	556	0	0
0	629	522	2	0
0	631	520	2	0
37	628	513	2	0
98	631	512	2	0
0	639	513	2	0
98	639	512	2	0
0	633	544	0	0
1	635	546	0	0
0	637	515	2	0
98	636	514	2	0
37	634	523	2	0
0	632	514	2	0
0	637	523	2	0
0	645	513	6	0
192	645	524	6	0
0	641	518	2	0
98	640	512	2	0
0	644	550	0	0
0	645	519	6	0
1	646	543	0	0
0	646	517	6	0
1	647	535	0	0
0	640	549	0	0
37	643	516	2	0
192	651	518	6	0
34	651	513	2	0
98	652	514	2	0
34	652	512	2	0
192	654	518	4	0
34	650	515	2	0
22	653	540	2	0
273	651	537	2	0
306	653	541	0	0
22	652	540	2	0
377	649	537	2	0
0	648	514	6	0
0	648	518	6	0
0	648	540	0	0
273	651	534	2	0
1	654	549	0	0
15	655	530	4	0
1	650	546	0	0
1	655	533	0	1
1	654	536	1	1
38	633	510	2	0
38	633	508	2	0
0	628	507	2	0
98	638	506	2	0
0	637	511	2	0
307	654	510	2	0
306	635	505	2	0
38	635	509	2	0
38	632	507	2	0
0	647	505	6	0
98	640	510	2	0
4	649	505	6	0
34	654	509	6	0
0	641	510	2	0
0	626	509	2	0
98	628	510	2	0
0	646	511	6	0
34	643	507	2	0
38	631	508	2	0
0	628	509	2	0
34	649	507	0	0
34	653	507	6	0
34	646	507	2	0
23	651	509	6	0
309	643	504	2	0
34	653	508	0	0
0	627	502	2	0
34	630	496	2	0
0	627	497	2	0
0	630	499	2	0
34	629	500	2	0
38	638	498	2	0
0	637	499	2	0
37	633	496	2	0
37	632	503	2	0
0	634	497	6	0
34	634	502	2	0
37	634	500	2	0
98	637	502	2	0
0	645	497	6	0
38	641	499	2	0
0	643	500	6	0
38	641	497	2	0
0	646	500	6	0
0	640	503	2	0
38	640	500	2	0
23	651	503	6	0
34	650	500	0	0
192	654	503	4	0
34	648	503	0	0
0	648	498	6	0
0	662	511	6	0
192	659	499	4	0
307	660	498	2	0
4	662	509	6	0
192	661	502	4	0
192	661	504	4	0
309	661	520	2	0
0	662	517	6	0
98	659	510	2	0
58	659	533	0	0
34	661	516	0	0
98	660	509	2	0
34	659	514	0	0
34	658	517	0	0
0	667	498	6	0
22	667	530	2	0
0	669	505	6	0
0	664	504	6	0
0	664	498	6	0
0	670	516	6	0
34	664	501	0	0
34	664	508	0	0
0	668	502	6	0
0	670	500	6	0
22	668	530	2	0
0	667	508	6	0
0	667	512	6	0
0	671	509	6	0
0	636	490	6	0
0	639	494	6	0
0	632	494	6	0
0	638	488	6	0
37	647	495	6	0
24	651	490	0	0
0	669	495	6	0
37	647	490	0	0
37	647	488	0	0
306	640	490	2	0
71	652	489	6	0
37	647	494	0	0
37	644	494	0	0
25	650	493	0	0
34	667	489	5	0
23	658	491	6	0
37	647	493	2	0
59	647	491	0	0
0	666	495	6	0
37	647	489	6	0
3	652	494	6	0
3	652	492	6	0
0	666	491	6	0
5	654	491	0	0
0	669	492	6	0
37	645	490	6	0
0	671	489	6	0
63	649	491	4	0
495	645	486	1	0
34	643	483	5	0
34	646	484	0	0
34	669	486	5	0
0	639	481	6	0
34	656	487	5	0
0	644	487	6	0
34	670	481	5	0
0	634	482	6	0
100	648	486	6	0
0	637	483	6	0
100	649	483	6	0
98	651	485	6	0
34	654	486	0	0
0	637	485	6	0
34	644	482	2	0
34	647	484	0	0
34	668	486	5	0
0	633	487	6	0
0	629	487	6	0
0	626	485	6	0
0	628	491	6	0
0	628	483	6	0
0	625	488	6	0
21	615	475	0	0
58	617	474	6	0
21	612	477	0	0
21	619	474	0	0
58	613	475	6	0
0	617	478	0	0
0	622	479	0	0
407	614	473	0	0
21	609	478	0	0
0	619	477	0	0
407	609	473	1	0
0	621	477	0	0
21	621	475	0	0
183	616	474	2	1
193	621	476	3	1
193	621	475	1	1
193	610	478	0	1
193	612	476	2	1
193	620	474	3	1
193	611	477	2	1
193	613	475	1	1
193	619	474	0	1
193	608	475	1	1
193	615	475	0	1
193	608	476	1	1
182	609	478	0	1
184	622	477	0	1
193	608	477	3	1
193	623	477	0	1
68	631	478	2	0
479	631	476	4	0
37	630	476	4	0
68	631	473	2	0
37	630	474	3	0
37	629	476	2	0
34	629	472	0	0
37	627	477	0	0
37	628	474	2	0
37	630	478	2	0
21	626	475	0	0
21	626	473	0	0
37	627	475	2	0
193	626	473	1	1
193	626	475	1	1
193	626	472	2	1
193	626	474	1	1
185	625	476	2	1
193	624	477	0	1
37	636	478	0	0
34	635	475	0	0
37	634	477	2	0
37	633	475	2	0
0	639	477	2	0
37	635	473	2	0
407	611	471	0	0
407	619	469	0	0
407	608	468	0	0
407	619	466	3	0
21	628	470	0	0
0	640	472	5	0
36	632	464	2	0
21	625	465	0	0
0	638	468	5	0
0	638	465	5	0
0	642	473	5	0
0	636	469	5	0
0	646	474	2	0
98	641	468	2	0
0	646	465	2	0
0	641	470	5	0
34	645	467	2	0
193	625	465	1	1
193	626	467	3	1
193	625	466	3	1
193	625	464	1	1
193	628	469	1	1
193	628	470	1	1
193	627	468	3	1
193	627	471	2	1
0	654	465	0	0
98	651	470	2	0
0	654	469	0	0
0	652	465	0	0
98	648	473	2	0
0	648	476	2	0
0	651	468	0	0
98	648	471	2	0
0	655	471	0	0
0	658	465	0	0
0	656	468	0	0
482	659	471	0	0
0	656	464	0	0
463	662	467	1	0
192	660	473	0	0
192	657	473	0	0
192	663	473	0	0
192	658	473	0	0
192	661	473	0	0
192	662	473	0	0
461	659	472	4	0
192	659	473	0	0
0	669	466	0	0
0	668	464	0	0
98	667	471	2	0
0	670	470	0	0
98	666	469	2	0
0	668	468	0	0
0	666	466	0	0
0	665	471	0	0
34	647	459	2	0
98	649	461	2	0
36	636	461	2	0
34	651	461	0	0
489	661	462	1	0
0	650	463	0	0
36	635	460	2	0
69	652	458	2	0
36	661	456	2	0
34	645	461	2	0
0	639	459	2	0
34	653	460	0	0
34	648	456	2	0
0	642	459	2	0
36	662	456	2	0
0	642	462	2	0
462	662	463	0	0
34	653	462	2	0
0	642	456	2	0
193	624	463	3	1
186	624	462	1	1
0	667	457	5	0
0	669	462	0	0
0	665	463	0	0
0	665	460	0	0
0	639	454	5	0
0	636	454	5	0
0	636	451	5	0
0	639	452	5	0
34	652	451	0	0
37	654	449	0	0
34	648	452	0	0
0	650	452	5	0
34	652	450	0	0
36	663	454	2	0
0	651	450	5	0
34	648	454	0	0
3	655	450	0	0
34	640	455	2	0
0	645	452	2	0
36	663	449	0	0
5	656	450	0	0
37	654	450	2	0
464	660	449	0	0
36	659	451	0	0
1	658	448	1	1
1	655	448	1	1
286	616	451	0	0
111	617	461	0	0
407	620	460	5	0
111	614	459	0	0
111	612	463	0	0
111	614	460	0	0
110	610	460	2	0
110	609	460	2	0
110	608	461	2	0
110	609	458	2	0
985	616	1436	0	0
383	612	455	0	0
150	615	455	4	0
383	613	455	0	0
153	614	455	0	0
153	611	455	0	0
153	608	455	0	0
153	609	455	0	0
153	610	455	0	0
286	609	447	0	0
283	609	444	0	0
283	616	443	0	0
0	669	448	5	0
36	664	452	0	0
0	670	453	5	0
36	664	453	0	0
34	654	442	0	0
34	654	447	0	0
36	660	445	0	0
34	653	441	0	0
34	651	442	0	0
23	646	447	1	0
0	651	440	5	0
100	649	440	0	0
36	663	443	0	0
34	652	447	0	0
36	665	447	2	0
100	670	441	4	0
100	668	445	4	0
36	659	446	0	0
36	662	442	2	0
55	657	447	2	0
0	670	445	5	0
55	656	445	3	0
100	656	441	2	0
34	651	444	0	0
100	670	443	4	0
34	652	446	0	0
55	655	446	2	0
34	652	443	0	0
34	652	442	0	0
0	654	444	5	0
36	655	434	2	0
0	647	433	5	0
0	644	434	5	0
0	643	433	5	0
36	653	437	2	0
0	650	433	5	0
98	653	434	0	0
100	651	436	0	0
98	650	435	0	0
0	648	436	5	0
36	663	432	2	0
100	649	438	4	0
36	658	433	2	0
36	667	432	2	0
36	671	436	2	0
0	679	438	3	0
597	672	435	6	0
0	676	436	3	0
597	672	437	6	0
306	675	432	3	0
597	672	456	6	0
597	672	440	6	0
306	677	452	4	0
306	674	440	0	0
597	672	469	6	0
597	672	441	6	0
0	678	458	3	0
306	673	461	4	0
313	678	443	3	0
597	672	448	6	0
597	672	457	6	0
597	672	459	6	0
34	678	462	3	0
597	672	468	6	0
597	672	447	6	0
597	672	458	6	0
597	672	442	6	0
34	677	445	3	0
597	672	467	6	0
597	672	436	6	0
597	672	434	6	0
34	675	437	3	0
597	672	439	6	0
597	672	432	6	0
597	672	466	6	0
597	672	438	6	0
597	672	433	6	0
313	679	435	3	0
597	672	460	6	0
597	672	461	6	0
597	672	451	6	0
597	672	455	6	0
597	672	454	6	0
34	677	467	3	0
597	672	463	6	0
597	672	449	6	0
597	672	453	6	0
597	672	471	6	0
597	672	470	6	0
597	672	443	6	0
597	672	446	6	0
0	674	468	3	0
597	672	444	6	0
597	672	462	6	0
34	675	454	3	0
597	672	445	6	0
597	672	450	6	0
597	672	452	6	0
597	672	465	6	0
597	672	464	6	0
15	656	1391	2	0
6	656	1394	0	0
470	650	1435	0	0
7	651	1436	6	0
3	652	1435	4	0
7	652	1433	6	0
3	653	1433	4	0
6	654	1435	0	0
3	653	1438	2	0
3	652	1436	4	0
7	652	1438	6	0
47	657	1435	4	0
3	655	1434	4	0
377	663	538	0	0
377	658	541	2	0
562	662	542	0	0
562	661	540	0	0
0	632	554	0	0
1	636	559	0	0
0	645	557	0	0
1	637	553	0	0
34	620	566	6	0
445	623	567	0	0
445	623	564	0	0
446	623	563	2	0
281	628	567	0	0
445	627	563	6	0
0	616	562	6	0
445	630	563	6	0
281	629	565	0	0
445	624	563	6	0
1	610	580	1	1
1	610	578	0	0
5	608	582	0	0
37	614	576	0	0
452	614	579	0	0
278	615	582	0	0
3	615	577	0	0
281	615	583	0	0
278	614	583	0	0
37	615	576	0	0
179	609	589	0	0
345	613	586	4	0
278	614	588	0	0
281	614	587	0	0
278	614	586	0	0
20	610	596	2	0
278	608	586	4	0
25	614	621	0	0
25	614	613	0	0
7	614	620	2	0
20	614	614	4	0
46	614	608	2	0
37	613	613	4	0
3	613	620	0	0
7	614	619	2	0
23	610	618	2	0
20	608	621	4	0
7	612	619	6	0
7	612	620	6	0
342	611	601	4	0
63	615	610	2	0
411	615	614	2	0
3	613	619	0	0
25	612	621	0	0
58	615	608	2	0
20	610	614	4	0
37	613	612	6	0
20	608	616	4	0
6	613	618	0	0
20	610	610	0	0
43	615	618	0	0
20	610	611	4	0
1	614	616	3	1
1	610	584	1	1
1	614	585	1	1
445	623	590	0	0
445	623	599	0	0
445	623	593	0	0
445	623	585	0	0
4	617	595	0	0
4	616	591	0	0
445	623	596	0	0
413	622	603	4	0
446	623	614	4	0
20	620	614	2	0
25	617	613	0	0
445	623	605	0	0
445	623	602	0	0
46	617	608	2	0
16	617	586	4	0
273	617	584	6	0
494	622	585	4	0
414	622	593	4	0
450	622	588	4	0
445	623	611	0	0
445	623	608	0	0
37	618	612	6	0
348	621	596	2	0
20	617	614	2	0
7	620	619	3	0
8	619	622	5	0
20	620	611	2	0
37	618	613	7	0
8	618	622	3	0
23	621	618	6	0
7	618	619	5	0
274	618	621	0	0
1	617	616	2	1
1	619	596	1	1
450	624	588	0	0
445	630	614	2	0
47	631	597	4	0
70	627	590	6	0
3	626	609	0	0
208	626	606	4	0
70	626	605	6	0
273	628	596	0	0
70	631	607	6	0
47	631	599	4	0
7	628	586	0	0
273	628	594	0	0
445	624	614	2	0
7	628	595	4	0
445	627	614	2	0
63	628	602	6	0
206	626	612	0	0
7	629	593	4	0
15	629	611	0	0
120	631	585	1	1
1	628	609	0	1
45	613	1546	2	0
25	613	1549	2	0
25	613	1548	2	0
45	613	1547	2	0
45	613	1545	2	0
25	608	1548	2	0
45	610	1545	6	0
14	618	1561	4	0
25	620	1562	4	0
3	615	1559	4	0
436	617	1559	4	0
25	608	1549	2	0
3	619	1565	4	0
5	611	1551	4	0
44	611	1545	4	0
45	610	1546	6	0
45	610	1547	6	0
25	617	1562	4	0
44	615	1562	0	0
15	611	1563	6	0
71	612	1565	2	0
45	617	1563	2	0
441	614	1562	4	0
1	618	1564	1	1
1	614	1564	1	1
97	609	1548	0	1
1	612	1562	0	1
338	610	2487	4	0
51	614	2507	6	0
3	615	2506	2	0
51	619	2509	2	0
6	611	2495	4	0
3	612	2508	2	0
51	618	2504	6	0
3	617	2509	2	0
3	616	2509	2	0
51	618	2509	6	0
6	615	2504	0	0
51	615	2507	0	0
51	612	2506	0	0
51	614	2509	6	0
51	619	2507	2	0
438	620	2506	2	0
51	618	2506	6	0
51	614	2506	0	0
51	618	2507	6	0
1	617	2507	0	1
1	613	2506	0	1
1	619	2508	1	1
1	615	2508	1	1
445	623	582	0	0
416	622	578	4	0
445	623	579	0	0
415	622	582	4	0
445	623	576	0	0
447	620	581	4	0
72	618	578	0	0
72	619	576	0	0
72	620	576	1	0
72	617	578	4	0
15	616	580	0	0
16	617	583	4	0
1	616	577	0	1
7	628	583	4	0
70	629	576	6	0
70	627	578	6	0
51	630	572	4	0
20	630	571	0	0
445	623	573	0	0
34	617	573	6	0
51	630	569	0	0
377	625	573	4	0
407	626	575	0	0
8	627	568	0	0
445	623	570	0	0
20	628	571	0	0
51	628	569	0	0
3	627	570	0	0
51	628	572	4	0
1	629	569	0	1
3	608	1523	0	0
290	627	1515	0	0
505	627	1514	0	0
6	608	1526	0	0
290	628	1513	0	0
290	627	1516	0	0
45	631	1511	6	0
45	631	1510	6	0
504	631	1514	0	0
290	629	1513	0	0
290	629	1516	0	0
22	631	1513	2	0
290	627	1513	0	0
208	636	587	2	0
51	632	572	4	0
502	635	564	0	0
23	639	582	0	0
407	637	567	0	0
20	634	569	0	0
3	632	570	0	0
205	636	584	2	0
23	637	578	0	0
23	637	580	0	0
445	633	563	6	0
79	632	590	0	0
51	635	572	4	0
445	639	563	6	0
23	639	580	0	0
23	632	565	0	0
19	638	576	0	0
377	636	573	4	0
23	639	578	0	0
407	638	570	0	0
23	637	582	0	0
3	635	570	0	0
445	636	563	6	0
20	634	571	0	0
41	632	566	2	0
1	632	569	0	1
138	633	573	0	1
205	639	594	2	0
209	637	596	2	0
209	636	597	2	0
209	636	605	2	0
206	639	606	1	0
204	634	607	1	0
217	632	603	5	0
123	637	606	0	1
205	640	588	2	0
208	646	596	4	0
5	642	597	0	0
273	645	584	2	0
20	644	589	0	0
208	643	604	2	0
71	645	587	2	0
20	646	589	0	0
55	640	598	0	0
7	644	584	6	0
55	640	597	0	0
71	646	605	2	0
4	643	595	4	0
1	646	579	1	1
1	642	581	1	1
120	642	600	0	1
1	646	603	1	1
122	645	569	0	1
41	647	565	4	0
278	645	566	0	0
273	644	565	0	0
278	645	564	0	0
445	645	563	6	0
445	642	563	6	0
278	643	565	0	0
55	634	608	1	0
445	633	614	2	0
42	637	608	4	0
55	634	609	7	0
445	636	614	2	0
445	639	614	2	0
123	635	612	0	1
407	640	617	1	0
1	619	628	6	0
0	633	628	6	0
54	643	628	4	0
6	621	626	6	0
37	615	624	4	0
20	611	624	2	0
0	628	625	6	0
0	626	626	6	0
391	628	629	0	0
1	614	628	0	0
0	633	624	6	0
1	621	629	0	0
37	614	624	7	0
20	620	624	4	0
1	618	626	0	0
37	617	624	6	0
1	608	628	0	0
99	622	626	6	0
0	611	629	0	0
0	623	624	6	0
0	625	631	6	0
0	624	628	6	0
37	616	624	2	0
0	630	626	6	0
0	625	624	6	0
0	624	625	6	0
209	651	605	6	0
366	648	589	4	0
4	652	599	4	0
7	649	604	2	0
20	650	589	6	0
20	652	589	6	0
407	648	618	6	0
209	653	605	5	0
510	648	602	6	0
4	653	596	4	0
3	648	604	6	0
4	648	600	4	0
121	648	585	1	1
1	655	591	1	1
1	648	588	0	1
435	618	3453	0	0
146	618	3451	4	0
20	616	3451	0	0
456	639	3440	4	0
457	637	3447	2	0
51	621	3456	0	0
51	618	3450	6	0
48	616	3448	6	0
51	640	3445	6	0
146	638	3450	2	0
51	640	3442	6	0
290	636	3440	4	0
51	617	3456	0	0
41	637	3440	4	0
51	617	3454	4	0
433	622	3436	0	0
146	635	3446	0	0
5	621	3458	0	0
22	640	3440	4	0
290	635	3440	4	0
55	635	3443	6	0
51	635	3445	2	0
55	640	3443	4	0
51	635	3442	2	0
7	648	578	6	0
3	649	578	0	0
7	650	578	2	0
80	623	3410	2	0
80	623	3409	2	0
433	621	3409	0	0
448	620	3414	6	0
51	622	3408	2	0
80	620	3421	6	0
449	636	3422	2	0
80	620	3422	6	0
51	622	3420	2	0
433	621	3422	0	0
51	620	3408	0	0
337	614	3401	0	0
51	610	3399	2	0
337	614	3399	0	0
51	619	3392	6	0
51	610	3396	2	0
51	618	3396	6	0
51	618	3399	6	0
445	648	563	6	0
445	651	563	6	0
50	655	567	0	0
445	654	563	6	0
1	648	554	0	0
0	649	556	0	0
0	653	556	0	0
13	663	584	6	0
13	663	597	6	0
13	661	590	6	0
12	662	587	6	0
13	661	584	6	0
12	661	596	6	0
209	660	599	6	0
407	661	602	0	0
407	656	615	0	0
11	658	591	0	0
3	658	590	0	0
12	661	592	6	0
407	659	612	3	0
209	663	599	6	0
278	658	594	0	0
278	657	594	0	0
13	660	594	6	0
12	663	594	6	0
13	660	578	6	0
213	663	581	6	0
12	667	581	6	0
12	670	583	6	0
13	670	586	6	0
13	668	590	6	0
12	664	577	6	0
407	667	600	6	0
12	668	579	6	0
213	666	585	6	0
12	667	588	6	0
13	665	595	6	0
213	664	591	6	0
143	667	569	2	0
29	664	568	0	0
1	657	570	0	1
50	659	567	0	0
445	660	563	6	0
445	663	563	6	0
143	670	566	2	0
445	657	563	6	0
445	666	563	6	0
143	664	566	0	0
445	669	563	6	0
143	666	566	0	0
0	668	561	0	0
143	668	566	0	0
1	670	567	1	1
1	657	566	0	1
0	656	559	0	0
0	662	559	0	0
377	663	546	4	0
562	661	544	0	0
1	657	545	0	0
377	656	550	4	0
22	656	547	2	0
563	660	551	2	0
15	643	1511	6	0
45	646	1511	6	0
42	647	1509	4	0
15	643	1508	6	0
281	643	1510	6	0
15	634	1513	2	0
45	646	1509	6	0
45	646	1510	6	0
42	632	1510	2	0
15	634	1515	2	0
45	646	1512	6	0
141	633	1513	0	1
0	670	554	0	0
0	664	559	0	0
37	674	559	6	0
0	679	553	1	0
306	677	557	1	0
34	677	554	1	0
4	672	565	2	0
209	673	568	4	0
206	676	565	5	0
205	674	569	2	0
445	672	563	6	0
445	678	563	6	0
22	675	566	2	0
204	677	565	0	0
0	674	552	1	0
407	673	585	0	0
37	673	558	6	0
55	679	565	3	0
407	678	585	0	0
22	675	565	2	0
38	678	582	2	0
34	674	561	1	0
205	673	566	2	0
445	675	563	6	0
8	679	586	3	0
407	676	579	0	0
1	676	568	0	1
37	686	555	4	0
4	681	552	1	0
34	683	555	1	0
38	686	553	6	0
0	682	558	1	0
20	684	588	0	0
24	686	587	6	0
3	683	590	0	0
204	684	568	0	0
0	685	557	1	0
445	687	563	6	0
206	682	567	6	0
206	684	566	7	0
3	681	584	0	0
37	680	557	6	0
22	681	586	0	0
22	685	590	0	0
55	687	565	4	0
445	684	563	6	0
88	680	564	5	0
4	680	569	2	0
38	687	577	2	0
34	680	554	1	0
22	681	587	0	0
29	685	588	0	0
22	687	588	2	0
22	686	590	0	0
445	681	563	6	0
3	682	565	2	0
22	687	589	2	0
34	682	561	1	0
3	686	583	0	0
22	681	566	2	0
23	682	585	0	0
1	685	586	0	1
1	686	567	1	1
55	694	568	5	0
8	688	585	0	0
97	689	570	5	0
204	692	570	2	0
1	692	568	0	1
206	693	566	0	0
3	693	565	2	0
445	693	563	6	0
204	691	565	2	0
88	695	565	5	0
55	688	564	6	0
4	688	561	2	0
88	700	573	5	0
38	697	582	2	0
37	690	560	4	0
88	697	589	0	0
407	698	580	0	0
445	690	563	6	0
206	689	565	7	0
88	699	568	5	0
88	697	571	5	0
407	701	576	0	0
4	696	567	2	0
306	698	552	1	0
37	689	556	6	0
37	696	554	6	0
37	689	552	1	0
37	694	553	6	0
34	693	555	1	0
0	690	554	1	0
4	695	556	1	0
0	690	558	1	0
37	693	546	3	0
34	683	544	1	0
306	684	550	1	0
306	677	548	1	0
38	681	550	6	0
377	664	550	4	0
0	681	547	1	0
38	675	545	6	0
306	687	545	1	0
562	667	544	0	0
0	674	546	1	0
562	664	544	0	0
34	686	547	1	0
306	691	548	2	0
34	694	551	1	0
0	689	549	1	0
34	678	544	1	0
37	702	550	4	0
0	699	549	1	0
34	699	546	1	0
0	696	551	1	0
0	699	544	1	0
34	697	545	1	0
0	697	547	1	0
34	696	546	1	0
34	698	546	1	0
34	711	547	1	0
306	708	548	1	0
34	709	546	1	0
0	709	544	1	0
34	711	545	1	0
0	707	547	1	0
38	705	581	2	0
407	706	559	2	0
407	707	572	2	0
37	707	545	6	0
34	706	544	1	0
0	705	550	1	0
407	705	563	6	0
34	707	552	1	0
407	706	585	0	0
407	710	595	2	0
407	715	556	0	0
407	717	586	3	0
407	714	576	5	0
725	715	580	2	0
209	713	583	4	0
209	713	579	4	0
407	713	572	4	0
407	715	566	6	0
407	712	566	4	0
407	712	560	0	0
407	725	563	2	0
730	723	591	5	0
884	725	581	4	0
407	722	571	6	0
407	727	568	2	0
877	725	588	7	0
884	726	580	0	0
407	721	562	7	0
877	730	588	7	0
877	731	582	7	0
884	729	582	2	0
97	731	571	0	0
407	733	565	7	0
884	728	579	0	0
884	728	583	6	0
272	616	638	0	0
547	617	636	3	0
360	617	639	0	0
272	615	638	0	0
0	613	638	0	0
0	625	635	6	0
0	627	637	6	0
272	616	639	0	0
0	629	633	6	0
272	615	639	0	0
0	613	635	0	0
412	623	634	4	0
0	609	637	4	0
1	611	639	3	0
272	616	637	0	0
0	628	635	6	0
0	627	632	6	0
0	630	632	6	0
51	598	3462	4	0
51	598	3459	0	0
51	596	3461	2	0
51	596	3460	2	0
51	598	3456	0	0
51	603	3456	0	0
51	593	3470	4	0
51	592	3465	6	0
20	607	3458	0	0
51	603	3464	0	0
439	611	3451	0	0
5	613	3450	0	0
51	608	3462	0	0
51	608	3460	4	0
97	608	3465	4	0
17	603	3466	2	0
3	613	3452	4	0
7	612	3452	6	0
20	606	3466	0	0
433	622	3464	0	0
51	593	3460	6	0
51	593	3461	6	0
51	603	3467	4	0
51	593	3462	4	0
51	593	3456	0	0
97	609	3467	4	0
97	608	3466	4	0
51	615	3456	0	0
80	621	3466	0	0
51	608	3456	0	0
51	601	3467	4	0
17	607	3457	4	0
51	615	3454	4	0
7	614	3452	3	0
20	606	3463	0	0
51	606	3457	0	0
51	608	3457	0	0
51	606	3462	0	0
97	609	3463	4	0
3	613	3451	4	0
51	612	3456	0	0
51	606	3467	4	0
51	608	3467	4	0
97	608	3464	4	0
1	602	3468	0	1
1	607	3468	0	1
1	607	3462	0	1
1	597	3463	0	1
1	601	3460	1	1
1	615	642	3	0
0	609	641	3	0
0	605	642	3	0
1	620	647	7	0
0	604	647	3	0
0	601	642	3	0
0	609	647	3	0
0	619	644	3	0
1	606	646	4	0
1	607	644	3	0
0	606	640	3	0
306	616	646	4	0
1	611	643	3	0
0	614	641	3	0
1	612	645	3	0
0	614	646	3	0
0	631	640	6	0
1	617	644	3	0
20	605	3477	2	0
51	603	3474	0	0
51	601	3476	4	0
51	602	3474	0	0
432	621	3477	0	0
80	619	3473	4	0
80	620	3473	4	0
5	616	3479	7	0
51	596	3479	4	0
51	593	3479	4	0
51	597	3477	6	0
372	594	3475	0	0
51	597	3475	6	0
373	594	3478	4	0
51	589	3462	2	0
51	590	3479	4	0
51	587	3479	4	0
23	585	3462	4	0
80	586	3471	4	0
372	590	3475	0	0
20	591	3468	2	0
51	587	3465	4	0
51	589	3465	2	0
51	584	3479	4	0
51	588	3457	0	0
375	591	3474	0	0
51	585	3457	0	0
51	584	3462	0	0
51	584	3465	4	0
20	584	3477	4	0
374	590	3478	4	0
51	587	3462	0	0
51	591	3459	0	0
22	584	3457	1	1
1	591	3462	1	1
1	584	3463	1	1
1	584	3464	1	1
1	588	3463	1	1
1	588	3464	1	1
5	583	3462	0	0
5	582	3456	0	0
51	582	3462	6	0
5	581	3461	0	0
5	583	3461	0	0
5	579	3456	0	0
5	581	3456	0	0
362	580	3470	0	0
279	580	3456	4	0
364	581	3469	6	0
51	581	3468	0	0
5	578	3469	2	0
51	577	3467	6	0
51	579	3473	2	0
51	578	3479	4	0
365	581	3471	2	0
51	576	3469	2	0
363	582	3470	4	0
51	577	3475	6	0
51	576	3477	2	0
5	578	3458	0	0
5	577	3457	0	0
51	579	3468	0	0
7	578	3468	6	0
51	578	3464	0	0
51	582	3465	6	0
51	576	3473	2	0
5	579	3461	0	0
51	576	3465	2	0
5	580	3473	4	0
51	581	3479	4	0
51	577	3471	6	0
51	580	3464	0	0
51	579	3476	2	0
1	578	3459	1	1
111	583	3476	1	1
1	579	3462	1	1
1	583	3468	1	1
0	567	646	0	0
307	563	646	0	0
1	572	641	0	0
0	585	641	0	0
48	579	640	0	0
0	571	647	0	0
1	581	640	0	1
0	588	641	4	0
20	603	3481	2	0
80	590	3481	4	0
51	581	3484	4	0
20	576	3484	4	0
51	578	3484	4	0
51	610	3485	4	0
80	606	3480	4	0
51	604	3480	0	0
51	607	3480	0	0
51	610	3482	0	0
55	608	3485	3	0
51	588	3483	2	0
51	607	3485	4	0
51	613	3482	0	0
51	592	3482	6	0
51	582	3482	2	0
80	614	3482	4	0
51	590	3484	4	0
17	615	3482	4	0
51	614	3485	4	0
1	577	3483	0	1
1	585	3482	1	1
1	588	3482	1	1
1	603	3480	0	1
80	622	3480	0	0
80	617	3482	4	0
17	616	3482	4	0
51	616	3485	4	0
80	613	3494	4	0
51	614	3494	0	0
51	616	3488	0	0
5	617	3489	4	0
20	621	3495	4	0
51	618	3490	4	0
51	620	3494	0	0
51	616	3490	4	0
51	618	3488	0	0
1	617	3491	0	1
1	617	3494	0	1
1	617	3488	0	1
51	618	3500	4	0
51	621	3500	4	0
22	617	3500	4	0
20	613	3497	4	0
20	618	3497	4	0
97	620	3497	4	0
51	616	3498	4	0
51	614	3498	4	0
22	617	3499	4	0
20	621	3499	4	0
1	571	654	0	0
406	571	651	0	0
1	582	648	2	0
1	582	651	2	0
1	582	655	2	0
98	590	655	4	0
98	590	651	3	0
98	590	654	4	0
1	607	649	3	0
1	578	655	2	0
1	585	651	2	0
1	606	651	3	0
1	604	654	2	0
1	598	650	2	0
1	600	652	2	0
1	596	655	2	0
1	599	655	2	0
98	588	654	4	0
98	588	652	3	0
98	588	655	5	0
1	580	653	2	0
1	580	649	2	0
1	586	648	2	0
1	593	652	2	0
1	595	652	2	0
1	593	650	2	0
1	595	648	2	0
306	567	652	0	0
0	564	650	0	0
0	564	654	0	0
9	562	649	1	1
51	561	1593	6	0
51	582	1580	6	0
365	581	1583	2	0
363	582	1582	4	0
45	560	1590	2	0
45	560	1591	2	0
362	580	1582	0	0
6	578	1580	2	0
364	581	1581	6	0
5	578	1581	2	0
51	582	2527	6	0
51	582	2524	6	0
363	582	2526	4	0
365	581	2527	2	0
364	581	2525	6	0
6	578	2525	2	0
362	580	2526	0	0
278	578	2524	4	0
20	559	643	6	0
55	558	644	6	0
20	558	643	0	0
43	558	645	4	0
286	555	646	0	0
55	559	644	7	0
377	668	541	6	0
562	664	540	0	0
562	666	542	0	0
562	667	539	0	0
150	666	536	0	1
0	677	543	1	0
0	674	543	1	0
0	674	540	1	0
0	675	536	1	0
597	674	531	4	0
597	675	531	4	0
597	677	531	4	0
597	676	531	4	0
597	672	529	6	0
597	672	531	6	0
597	672	530	6	0
597	678	531	4	0
597	673	531	4	0
597	672	528	6	0
597	679	531	4	0
313	679	522	3	0
306	679	525	2	0
0	678	520	1	0
34	677	525	2	0
0	676	527	1	0
0	676	523	1	0
597	672	523	6	0
597	672	526	6	0
597	672	522	6	0
313	676	522	3	0
597	672	520	6	0
597	672	524	6	0
0	674	520	1	0
597	672	525	6	0
34	676	520	2	0
597	672	527	6	0
597	672	521	6	0
306	674	524	1	0
0	676	517	1	0
313	678	516	0	0
310	678	518	0	0
306	679	514	4	0
34	676	513	2	0
597	672	517	6	0
597	672	518	6	0
597	672	516	6	0
597	672	519	6	0
0	674	513	1	0
313	674	516	3	0
597	672	513	6	0
597	672	514	6	0
597	672	512	6	0
597	672	515	6	0
306	680	542	5	0
0	684	539	1	0
0	686	542	1	0
37	688	543	2	0
306	694	542	1	0
34	693	540	1	0
0	691	543	1	0
306	688	537	2	0
306	695	532	2	0
597	691	531	4	0
0	683	533	1	0
597	690	531	4	0
597	694	531	4	0
597	689	531	4	0
597	682	531	4	0
597	681	531	4	0
597	686	531	4	0
0	686	535	1	0
597	687	531	4	0
597	685	531	4	0
597	692	531	4	0
597	693	531	4	0
597	688	531	4	0
306	680	534	2	0
597	684	531	4	0
597	680	531	4	0
0	691	533	1	0
597	695	531	4	0
597	683	531	4	0
37	687	521	6	0
313	686	524	0	0
307	685	526	4	0
37	685	522	5	0
313	683	523	0	0
5	684	520	2	0
0	683	525	1	0
313	693	524	3	0
34	682	527	2	0
313	694	520	0	0
313	694	527	0	0
313	690	525	3	0
313	688	521	0	0
0	693	521	1	0
34	680	523	2	0
0	695	524	1	0
34	688	523	2	0
420	681	520	0	0
420	689	525	0	0
306	695	526	1	0
0	691	522	1	0
5	692	525	2	0
0	697	525	1	0
313	698	521	3	0
34	701	538	4	0
597	700	531	4	0
34	699	540	1	0
34	701	534	4	0
34	700	541	1	0
34	697	522	2	0
313	696	523	3	0
34	699	535	4	0
34	699	542	1	0
597	701	531	4	0
565	700	540	2	0
34	701	541	3	0
34	700	539	4	0
34	700	537	2	0
313	699	527	3	0
565	702	531	2	0
597	698	531	4	0
597	696	531	4	0
37	697	538	4	0
34	697	543	1	0
626	703	531	2	0
597	697	531	4	0
34	702	533	4	0
597	699	531	4	0
565	701	535	2	0
34	708	538	4	0
34	710	542	1	0
34	708	540	4	0
34	707	541	7	0
34	708	537	4	0
34	708	543	1	0
34	706	536	4	0
0	711	540	1	0
306	711	543	5	0
565	707	540	2	0
34	707	539	5	0
34	707	536	4	0
34	706	540	6	0
306	711	533	2	0
597	709	531	4	0
54	709	533	1	0
597	710	531	4	0
565	705	531	2	0
597	706	531	4	0
34	706	534	4	0
565	706	535	2	0
597	711	531	4	0
597	708	531	4	0
597	707	531	4	0
34	707	533	4	0
34	711	524	2	0
0	709	525	0	0
0	710	521	0	0
34	707	523	2	0
310	718	520	0	0
597	719	531	4	0
306	713	523	1	0
0	718	522	0	0
306	718	525	1	0
597	715	531	4	0
597	714	531	4	0
0	715	532	1	0
597	713	531	4	0
0	717	545	1	0
34	712	541	1	0
0	713	537	1	0
34	716	548	1	0
597	718	531	4	0
0	713	540	1	0
38	713	546	6	0
0	714	545	1	0
0	715	542	1	0
0	713	548	1	0
597	716	531	4	0
306	718	532	2	0
34	715	546	1	0
306	716	540	1	0
597	712	531	4	0
597	717	531	4	0
34	725	545	3	0
306	724	547	0	0
34	721	547	3	0
37	724	543	3	0
306	724	539	0	0
37	726	536	3	0
34	723	536	0	0
0	721	543	3	0
0	726	542	3	0
34	720	539	3	0
34	721	537	0	0
34	722	537	0	0
597	727	531	4	0
597	726	531	4	0
306	722	528	4	0
34	720	535	0	0
597	722	531	4	0
597	721	531	4	0
597	724	531	4	0
597	723	531	4	0
597	725	531	4	0
597	720	531	4	0
34	725	534	3	0
310	727	526	0	0
306	727	523	0	0
5	724	521	4	0
420	724	524	2	0
34	721	526	3	0
0	721	522	0	0
0	709	512	0	0
313	693	517	0	0
0	696	517	1	0
0	696	514	1	0
306	689	519	4	0
306	696	519	1	0
34	719	513	2	0
37	689	517	4	0
34	696	513	0	0
34	693	513	0	0
0	717	518	0	0
34	718	516	2	0
597	689	512	0	0
420	714	519	2	0
597	692	512	0	0
0	720	515	1	0
34	726	512	3	0
37	692	518	0	0
306	721	517	6	0
597	695	512	0	0
0	694	518	1	0
37	694	515	1	0
597	694	512	0	0
306	698	514	0	0
313	688	513	0	0
5	714	516	0	0
5	691	515	2	0
597	693	512	0	0
597	688	512	0	0
597	691	512	0	0
306	709	515	3	0
34	710	519	2	0
597	690	512	0	0
597	696	512	0	0
34	706	519	2	0
419	714	512	6	0
420	688	515	0	0
38	725	519	2	0
34	712	513	2	0
420	726	515	0	0
306	716	513	0	0
306	722	513	2	0
313	712	517	3	0
34	687	518	2	0
313	687	519	0	0
597	687	512	0	0
0	687	517	1	0
313	686	516	3	0
34	686	513	2	0
307	683	516	4	0
597	683	512	0	0
597	684	512	0	0
597	680	512	0	0
597	685	512	0	0
0	684	514	1	0
37	685	519	4	0
0	680	517	1	0
597	686	512	0	0
34	681	513	2	0
597	682	512	0	0
597	681	512	0	0
34	683	513	0	0
718	734	529	2	0
306	734	522	0	0
0	730	551	3	0
306	735	543	0	0
34	728	549	3	0
34	733	546	3	0
5	729	515	6	0
306	731	535	0	0
0	735	533	3	0
38	734	520	2	0
0	728	520	0	0
306	730	514	0	0
718	733	530	2	0
597	731	531	4	0
0	732	531	4	0
0	735	539	3	0
34	733	538	3	0
37	735	537	3	0
718	735	528	2	0
0	730	517	1	0
306	733	517	4	0
597	728	531	4	0
23	735	516	6	0
597	729	531	4	0
34	732	520	3	0
5	731	522	4	0
23	735	512	6	0
34	732	522	3	0
420	731	525	2	0
0	731	540	3	0
306	729	528	4	0
306	729	543	0	0
597	730	531	4	0
0	728	537	3	0
34	734	549	3	0
306	732	554	0	0
0	740	527	3	0
718	740	523	2	0
718	743	520	2	0
718	742	521	2	0
718	741	522	2	0
34	736	536	3	0
718	737	526	2	0
718	736	527	2	0
34	743	533	3	0
0	743	536	3	0
306	741	540	3	0
38	737	523	2	0
34	738	539	3	0
718	739	524	2	0
36	738	520	1	0
718	738	525	2	0
34	737	530	3	0
192	740	520	0	0
306	738	532	0	0
0	742	544	3	0
34	738	546	3	0
0	739	535	3	0
37	743	538	3	0
306	739	549	3	0
0	748	524	3	0
0	751	523	3	0
0	751	526	3	0
37	750	521	3	0
306	744	524	3	0
34	746	522	3	0
69	748	542	3	0
306	744	530	3	0
34	751	535	3	0
0	750	540	3	0
34	749	538	3	0
34	749	529	3	0
0	750	546	3	0
306	746	535	0	0
0	750	531	3	0
0	747	533	3	0
0	745	547	3	0
38	750	549	3	0
306	746	551	0	0
407	746	562	5	0
407	741	571	5	0
407	737	569	5	0
407	746	570	5	0
407	750	574	5	0
877	742	580	7	0
729	751	586	3	0
877	748	584	7	0
877	742	587	7	0
729	740	585	5	0
730	737	591	3	0
877	742	584	7	0
877	738	588	7	0
407	755	566	5	0
877	754	584	7	0
729	755	590	2	0
407	759	563	5	0
407	752	561	5	0
306	759	544	3	0
34	755	545	3	0
0	755	548	3	0
34	753	550	3	0
37	754	543	3	0
0	758	538	3	0
0	752	543	3	0
0	753	537	3	0
34	755	541	3	0
0	758	534	3	0
34	755	535	3	0
306	754	531	3	0
0	758	531	3	0
306	754	521	3	0
0	756	526	3	0
36	736	519	1	0
192	737	518	0	0
597	746	513	2	0
597	746	517	2	0
718	744	519	2	0
192	738	519	0	0
192	744	514	0	0
306	750	515	3	0
718	745	518	2	0
194	736	514	0	0
597	746	514	2	0
597	746	516	2	0
597	746	515	2	0
597	746	512	2	0
0	748	519	3	0
192	737	512	0	0
37	756	515	3	0
0	752	519	3	0
0	758	519	3	0
0	759	513	3	0
0	756	517	3	0
0	752	513	3	0
0	725	510	0	0
668	727	504	0	0
306	724	506	5	0
38	723	507	2	0
306	721	509	5	0
36	737	511	7	0
306	731	508	1	0
34	735	508	3	0
597	746	508	2	0
597	746	505	2	0
310	734	506	0	0
668	733	504	0	0
597	746	506	2	0
597	746	511	2	0
0	737	506	0	0
34	754	508	3	0
34	749	511	3	0
34	737	504	3	0
0	754	504	3	0
37	752	504	3	0
0	755	511	3	0
34	720	504	3	0
192	741	505	0	0
34	721	505	3	0
36	740	507	7	0
668	729	505	0	0
36	743	507	7	0
164	728	507	0	0
36	740	504	7	0
597	746	507	2	0
597	746	509	2	0
0	721	506	1	0
36	744	508	1	0
306	759	507	3	0
192	739	510	0	0
597	746	510	2	0
0	750	507	3	0
597	746	504	2	0
668	731	505	0	0
192	744	511	0	0
34	764	542	3	0
0	763	538	3	0
34	766	548	3	0
0	764	546	3	0
34	762	541	3	0
34	760	548	3	0
407	764	570	5	0
306	766	540	3	0
37	763	536	3	0
407	764	560	5	0
306	763	532	3	0
34	761	535	3	0
0	765	529	3	0
38	760	533	3	0
34	766	535	3	0
34	761	529	3	0
0	766	526	3	0
0	767	522	3	0
0	764	522	3	0
34	760	523	3	0
306	761	526	3	0
0	760	519	3	0
306	763	517	3	0
306	766	512	3	0
0	766	515	3	0
0	763	509	3	0
0	764	506	3	0
37	766	509	3	0
34	766	506	3	0
667	731	501	0	0
668	733	500	0	0
668	731	499	0	0
668	734	502	0	0
668	729	499	0	0
36	742	502	7	0
36	740	502	7	0
597	746	497	2	0
36	739	500	7	0
164	736	501	0	0
192	737	496	0	0
597	746	500	2	0
597	746	502	2	0
597	746	498	2	0
36	736	497	7	0
597	746	501	2	0
306	738	502	6	0
597	746	499	2	0
0	766	498	3	0
0	754	497	3	0
164	730	498	0	0
597	746	496	2	0
0	748	499	3	0
597	746	503	2	0
0	750	497	3	0
0	751	496	3	0
306	749	502	3	0
0	759	503	3	0
69	756	500	3	0
0	765	503	3	0
0	760	497	3	0
306	761	499	3	0
34	743	488	3	0
34	750	488	3	0
597	746	488	2	0
306	756	492	3	0
597	746	489	2	0
597	746	494	2	0
34	751	493	3	0
38	740	491	2	0
36	734	491	6	0
0	741	489	1	0
34	742	492	3	0
306	748	490	3	0
34	740	493	3	0
0	738	492	1	0
597	746	490	2	0
34	753	492	3	0
0	766	488	3	0
34	766	494	3	0
0	754	489	3	0
597	746	495	2	0
0	741	495	1	0
36	736	490	3	0
34	733	488	3	0
0	761	492	3	0
0	761	488	3	0
0	763	495	3	0
306	764	490	3	0
36	736	495	7	0
597	746	493	2	0
23	733	490	6	0
36	736	492	5	0
597	746	491	2	0
597	746	492	2	0
34	743	494	3	0
36	737	491	1	0
34	733	486	3	0
0	735	481	1	0
306	731	485	1	0
0	735	485	1	0
38	733	484	2	0
306	750	483	3	0
36	742	481	3	0
23	743	481	2	0
597	746	486	2	0
0	748	481	3	0
597	746	487	2	0
34	766	482	3	0
306	763	482	3	0
0	758	483	3	0
597	746	480	2	0
34	765	487	3	0
0	757	481	3	0
0	767	485	3	0
597	746	481	2	0
36	739	482	4	0
597	746	484	2	0
34	736	483	3	0
306	737	483	6	0
597	746	483	2	0
192	741	481	0	0
597	746	482	2	0
597	746	485	2	0
34	752	481	3	0
0	754	483	3	0
0	752	487	3	0
306	756	486	3	0
0	762	486	3	0
0	760	483	3	0
34	760	485	3	0
34	762	481	3	0
192	737	475	0	0
192	739	474	0	0
192	738	477	0	0
0	749	478	0	0
5	743	473	0	0
164	749	474	3	0
597	746	477	2	0
0	751	473	0	0
36	738	478	1	0
597	746	474	2	0
34	754	479	3	0
34	754	474	3	0
597	746	476	2	0
597	746	475	2	0
306	752	476	0	0
0	763	473	0	0
34	762	477	3	0
0	760	475	0	0
306	764	476	0	0
36	739	479	1	0
164	740	476	7	0
192	741	479	0	0
164	736	475	7	0
307	736	478	0	0
597	746	479	2	0
164	748	472	3	0
597	746	478	2	0
192	736	473	0	0
194	736	472	0	0
36	736	474	1	0
597	746	473	2	0
597	746	472	2	0
0	757	473	0	0
0	757	477	0	0
164	742	471	0	0
597	746	471	2	0
164	741	468	0	0
34	755	466	3	0
36	741	467	1	0
306	754	470	0	0
36	737	468	1	0
597	746	469	2	0
597	746	470	2	0
164	754	464	3	0
164	741	469	7	0
164	749	467	0	0
597	746	468	2	0
192	747	464	0	0
192	739	471	0	0
192	738	468	0	0
34	752	471	3	0
192	742	465	0	0
164	740	470	0	0
307	740	471	7	0
306	761	469	0	0
192	736	471	0	0
306	757	465	0	0
164	753	465	0	0
34	760	467	3	0
313	752	467	3	0
420	743	471	6	0
36	736	470	1	0
313	743	470	3	0
192	743	464	0	0
192	739	469	0	0
36	736	469	1	0
597	746	466	2	0
307	737	465	0	0
36	740	465	1	0
192	740	466	0	0
34	765	469	3	0
597	746	465	5	0
428	743	467	6	0
0	763	465	0	0
0	766	471	0	0
0	766	467	0	0
36	740	468	1	0
313	749	470	3	0
597	746	467	2	0
5	743	461	0	0
0	738	461	1	0
718	750	461	2	0
36	749	463	1	0
718	751	460	2	0
420	743	458	6	0
193	745	463	0	0
306	747	457	7	0
36	753	461	1	0
36	756	458	1	0
597	744	456	0	0
164	739	460	7	0
702	745	456	6	0
718	739	456	0	0
597	740	457	0	0
192	754	461	0	0
307	745	461	7	0
37	739	462	1	0
597	743	456	0	0
36	742	463	1	0
597	742	456	2	0
313	742	460	3	0
597	741	457	0	0
0	737	459	1	0
718	748	463	2	0
34	765	456	3	0
597	749	462	5	0
34	762	456	3	0
306	764	459	0	0
0	762	462	0	0
34	765	463	3	0
597	752	456	2	0
0	766	456	0	0
34	763	457	3	0
597	752	457	2	0
597	752	459	2	0
597	752	458	2	0
597	742	455	4	0
718	738	455	0	0
718	736	453	0	0
597	749	451	2	0
597	749	449	2	0
597	749	450	2	0
718	737	454	0	0
192	756	454	0	0
597	749	448	2	0
192	758	451	0	0
36	759	455	3	0
34	766	452	3	0
192	764	450	0	0
597	752	448	2	0
34	765	454	3	0
597	752	449	2	0
597	752	451	2	0
36	755	455	3	0
597	752	450	2	0
597	752	454	2	0
36	755	453	1	0
597	752	453	2	0
37	755	448	7	0
597	752	452	2	0
194	760	448	0	0
192	761	453	0	0
597	752	455	2	0
718	748	454	2	0
565	748	453	0	0
718	747	455	2	0
37	749	454	7	0
313	750	455	3	0
597	749	452	2	0
597	749	453	2	0
718	735	453	2	0
0	728	473	1	0
164	732	472	3	0
0	733	473	1	0
0	728	477	1	0
23	734	471	6	0
0	733	477	1	0
37	729	469	1	0
37	730	462	7	0
164	732	474	0	0
597	731	456	0	0
34	730	458	7	0
164	733	476	0	0
306	735	462	2	0
164	731	465	2	0
0	731	458	1	0
164	730	466	4	0
307	732	468	0	0
0	729	479	1	0
34	734	477	7	0
34	729	465	7	0
0	728	457	1	0
0	729	468	1	0
718	729	455	0	0
23	735	457	0	0
0	735	464	1	0
313	735	459	3	0
313	735	466	3	0
718	734	454	2	0
0	730	464	1	0
642	729	450	6	0
164	731	461	2	0
597	732	456	0	0
597	730	456	0	0
718	728	454	0	0
718	733	455	2	0
597	743	443	4	0
34	742	441	7	0
34	749	441	7	0
565	748	446	0	0
597	742	443	2	0
718	739	443	2	0
718	737	445	2	0
597	740	442	4	0
718	738	444	2	0
718	747	444	0	0
313	755	442	3	0
36	759	441	3	0
597	752	443	2	0
37	744	440	7	0
36	759	440	3	0
597	752	441	2	0
597	752	440	2	0
597	752	442	2	0
597	736	440	2	0
597	746	443	4	0
597	752	446	2	0
306	746	441	7	0
597	752	447	2	0
597	749	446	2	0
597	745	443	4	0
313	751	444	3	0
597	752	445	2	0
597	749	447	2	0
718	736	446	2	0
192	757	445	0	0
597	742	444	0	0
718	748	445	0	0
597	744	443	4	0
597	741	442	4	0
597	752	444	2	0
194	760	446	0	0
306	747	432	7	0
164	750	433	4	0
420	750	435	6	0
5	745	438	4	0
597	752	436	2	0
37	757	432	7	0
313	744	435	3	0
597	752	435	2	0
164	756	438	3	0
597	752	434	2	0
597	752	432	2	0
192	760	435	0	0
597	752	433	2	0
192	759	438	0	0
164	757	437	4	0
5	750	438	4	0
37	755	435	7	0
164	758	435	5	0
597	752	438	2	0
164	755	436	4	0
313	760	432	3	0
597	752	437	2	0
597	752	439	2	0
420	745	435	6	0
36	761	433	3	0
192	762	434	0	0
34	710	505	2	0
597	697	507	2	0
597	697	508	2	0
313	711	508	3	0
54	709	508	2	0
306	706	504	0	0
597	697	509	2	0
306	718	509	1	0
597	697	505	2	0
525	693	508	6	0
597	697	511	2	0
34	699	509	2	0
419	714	504	6	0
597	697	504	2	0
597	697	506	2	0
575	714	508	0	0
597	697	510	2	0
306	711	510	0	0
313	716	511	3	0
306	716	504	0	0
421	724	1468	2	0
576	714	1452	0	0
6	724	1465	0	0
6	729	1459	2	0
421	726	1459	0	0
6	714	1460	0	0
418	714	1456	6	0
6	731	1466	0	0
418	714	1448	6	0
421	731	1469	2	0
421	714	1463	2	0
597	680	504	2	0
597	680	511	2	0
597	680	507	2	0
597	680	505	2	0
597	680	509	2	0
597	680	510	2	0
597	680	506	2	0
597	680	508	2	0
525	683	508	2	0
34	679	507	0	0
313	679	504	3	0
306	675	509	1	0
34	678	511	2	0
313	678	510	3	0
313	675	504	3	0
597	672	511	6	0
313	675	506	3	0
597	672	510	6	0
597	672	504	6	0
597	672	509	6	0
597	672	507	6	0
597	672	508	6	0
597	672	506	6	0
597	672	505	6	0
580	683	1452	2	0
6	684	1464	2	0
421	681	1464	0	0
121	686	1464	0	0
648	693	1452	6	0
6	691	1459	2	0
121	693	1459	0	0
121	694	1469	0	0
6	692	1469	2	0
421	688	1459	0	0
421	689	1469	0	0
597	672	503	6	0
306	675	500	1	0
597	672	496	6	0
597	672	498	6	0
597	680	503	2	0
34	679	501	0	0
597	672	499	6	0
313	673	496	3	0
597	680	502	2	0
597	680	501	2	0
653	681	502	4	0
597	672	497	6	0
597	680	499	2	0
653	683	502	4	0
654	683	497	4	0
647	692	503	0	0
597	680	498	2	0
597	680	496	2	0
655	692	496	0	0
597	680	497	2	0
597	680	500	2	0
647	690	503	0	0
647	694	503	0	0
597	672	500	6	0
653	685	502	4	0
597	672	501	6	0
597	672	502	6	0
530	692	499	0	0
34	678	499	0	0
655	692	497	0	0
655	692	498	0	0
34	699	503	2	0
597	697	499	2	0
597	697	502	2	0
597	697	501	2	0
597	697	497	2	0
306	699	496	2	0
597	697	496	2	0
597	697	498	2	0
597	697	503	2	0
597	697	500	2	0
34	708	501	2	0
306	711	501	0	0
313	710	500	3	0
306	709	497	1	0
313	714	502	0	0
313	717	501	3	0
34	716	501	0	0
306	718	499	1	0
34	719	503	2	0
420	714	497	6	0
313	717	497	3	0
5	714	500	0	0
34	716	502	2	0
0	722	503	1	0
306	720	497	1	0
668	727	500	0	0
38	724	497	2	0
306	725	499	6	0
34	720	501	3	0
313	719	494	3	0
0	723	495	1	0
34	722	495	3	0
0	725	494	1	0
34	715	494	2	0
313	716	490	3	0
34	718	489	2	0
34	726	492	3	0
306	714	492	1	0
313	712	495	3	0
0	713	490	0	0
34	718	495	2	0
313	717	492	0	0
310	718	493	0	0
34	721	495	3	0
306	723	491	1	0
0	721	493	1	0
306	708	489	1	0
306	711	492	0	0
34	709	494	2	0
0	710	495	0	0
597	697	494	2	0
607	697	492	2	0
597	697	493	2	0
597	697	495	2	0
597	694	491	0	0
597	695	491	0	0
597	693	491	0	0
34	695	490	0	0
306	691	489	0	0
597	691	491	0	0
597	690	491	0	0
597	688	491	0	0
34	693	490	0	0
655	692	495	0	0
597	692	491	0	0
597	689	491	0	0
597	687	491	0	0
34	681	490	0	0
597	680	495	2	0
597	686	491	0	0
306	685	489	0	0
597	685	491	0	0
597	681	491	0	0
597	682	491	0	0
597	680	494	2	0
597	684	491	0	0
34	684	489	0	0
597	680	493	2	0
652	683	495	0	0
597	680	491	0	0
597	680	492	2	0
597	683	491	0	0
420	690	480	6	0
34	693	480	2	0
661	703	486	0	0
313	687	482	3	0
608	701	482	0	0
5	685	483	4	0
5	690	483	4	0
313	696	480	3	0
306	707	481	7	0
5	714	481	0	0
608	706	482	0	0
640	712	482	2	0
420	714	484	2	0
420	680	480	6	0
5	680	483	4	0
34	694	482	4	0
420	685	480	6	0
668	702	485	1	0
306	698	481	0	0
668	702	487	1	0
313	682	481	3	0
668	704	485	1	0
313	710	480	3	0
668	704	487	1	0
34	718	480	2	0
306	723	481	1	0
34	721	480	3	0
6	714	1425	0	0
421	714	1428	2	0
6	714	1444	0	0
421	714	1441	6	0
668	739	3333	2	0
668	735	3335	2	0
22	727	3334	2	0
668	737	3332	2	0
97	737	3334	2	0
668	735	3333	2	0
5	729	3334	2	0
55	735	3332	2	0
668	739	3335	2	0
668	737	3336	2	0
34	722	473	7	0
575	714	476	0	0
0	720	474	7	0
34	723	477	7	0
313	714	472	3	0
164	725	478	4	0
0	725	472	7	0
306	726	474	2	0
37	720	479	7	0
525	710	476	2	0
34	706	474	3	0
34	699	472	3	0
0	697	475	3	0
0	691	473	3	0
306	695	477	2	0
34	691	477	3	0
313	693	476	3	0
306	686	474	0	0
313	685	478	3	0
420	680	473	2	0
0	683	477	3	0
306	675	474	2	0
313	677	482	3	0
34	679	478	3	0
597	672	484	6	0
597	672	482	6	0
313	673	481	3	0
597	672	487	6	0
597	672	480	6	0
597	672	494	6	0
597	679	491	6	0
597	672	492	6	0
597	672	493	6	0
597	672	481	6	0
597	672	490	6	0
597	672	491	6	0
34	674	479	3	0
597	672	486	6	0
597	672	474	6	0
597	672	489	6	0
597	672	495	6	0
597	672	485	6	0
597	672	479	6	0
597	672	478	6	0
306	677	494	1	0
597	672	483	6	0
306	674	488	1	0
597	672	473	6	0
597	672	477	6	0
597	672	475	6	0
597	672	472	6	0
597	672	476	6	0
313	676	491	3	0
597	672	488	6	0
598	687	2395	6	0
598	687	2397	6	0
598	689	2394	6	0
649	683	2396	2	0
598	689	2393	6	0
598	687	2394	6	0
581	693	2396	6	0
650	689	2395	6	0
598	687	2398	6	0
598	689	2396	6	0
590	687	2396	6	0
598	689	2397	6	0
0	694	471	3	0
34	687	465	3	0
0	687	469	3	0
428	691	467	6	0
34	698	464	3	0
306	697	468	2	0
34	685	467	3	0
306	681	464	2	0
0	683	471	3	0
5	680	470	0	0
6	680	1414	0	0
6	680	1427	4	0
421	685	1424	6	0
6	685	1427	4	0
641	685	1425	0	0
641	680	1425	0	0
421	680	1424	6	0
6	690	1427	4	0
421	680	1417	2	0
119	692	1412	0	0
427	692	1411	6	0
641	690	1425	0	0
0	710	465	3	0
306	708	470	7	0
0	707	467	3	0
34	707	464	3	0
428	716	466	6	0
313	719	464	3	0
306	717	470	2	0
0	713	469	3	0
34	727	471	7	0
0	727	468	7	0
0	726	465	1	0
0	723	466	7	0
313	723	471	3	0
306	721	469	0	0
37	721	466	1	0
576	714	1420	0	0
632	715	1418	4	0
620	716	1419	4	0
3	714	1418	0	0
7	716	1421	0	0
182	716	1412	0	0
22	712	1419	4	0
635	710	1420	2	0
22	713	1418	4	0
427	716	1410	6	0
0	714	463	3	0
575	717	458	4	0
34	724	463	7	0
37	726	461	7	0
306	721	460	4	0
34	713	459	3	0
313	716	456	3	0
37	723	456	7	0
0	712	456	3	0
306	727	462	2	0
428	725	457	0	0
0	719	462	3	0
419	720	457	0	0
419	716	461	6	0
313	709	462	3	0
0	708	459	3	0
425	711	457	2	0
306	711	461	7	0
525	703	462	4	0
524	703	459	2	0
524	703	456	2	0
34	700	462	3	0
306	695	461	7	0
34	689	459	3	0
575	691	458	0	0
425	695	458	4	0
419	691	462	2	0
428	682	458	4	0
0	686	463	3	0
34	684	462	3	0
419	687	458	0	0
0	685	452	3	0
34	685	455	3	0
419	690	449	2	0
525	690	451	2	0
526	691	453	0	0
419	695	454	2	0
34	682	451	3	0
34	683	454	3	0
0	684	449	3	0
419	695	448	2	0
0	680	450	3	0
419	691	454	2	0
524	693	451	0	0
524	695	451	0	0
427	683	1402	4	0
22	682	1401	0	0
22	681	1402	0	0
580	691	1395	2	0
574	692	1397	0	0
602	690	1401	2	0
5	691	1394	0	0
418	691	1393	2	0
418	688	1402	0	0
418	692	1398	2	0
418	692	1406	2	0
587	693	1402	0	0
578	694	1395	0	0
576	692	1402	0	0
5	690	1394	0	0
526	696	452	2	0
512	703	450	4	0
419	697	451	0	0
524	698	451	0	0
153	703	455	0	1
419	711	454	2	0
419	711	448	2	0
526	711	451	6	0
313	709	454	3	0
524	710	451	0	0
419	709	451	4	0
524	708	451	0	0
526	716	450	4	0
525	716	451	6	0
419	715	448	2	0
524	713	451	0	0
34	718	454	3	0
419	716	453	2	0
718	727	453	0	0
565	727	452	4	0
597	726	452	6	0
597	726	451	6	0
597	726	450	6	0
0	721	453	7	0
597	726	449	6	0
306	723	450	4	0
313	723	448	3	0
597	726	448	6	0
597	727	441	0	0
718	727	446	2	0
565	727	447	4	0
597	725	441	0	0
597	726	441	0	0
575	723	444	0	0
597	724	441	0	0
718	735	446	0	0
597	721	441	0	0
597	726	447	6	0
597	735	441	0	0
718	734	445	0	0
36	731	440	0	0
597	731	443	4	0
597	723	441	0	0
597	728	441	0	0
36	729	440	2	0
0	722	447	7	0
597	731	441	0	0
597	720	441	0	0
597	722	441	0	0
597	734	441	0	0
718	733	444	0	0
597	732	443	4	0
718	728	445	2	0
597	730	443	4	0
718	729	444	2	0
606	732	441	0	0
597	730	441	0	0
597	729	441	0	0
281	723	1389	0	0
281	722	1388	0	0
281	723	1387	0	0
418	720	1401	0	0
576	723	1388	0	0
6	743	1405	0	0
641	743	1419	0	0
597	743	1408	0	0
3	740	1384	4	0
427	725	1401	0	0
6	743	1418	0	0
281	724	1388	0	0
3	750	1384	6	0
3	745	1384	6	0
421	743	1402	6	0
597	743	1409	0	0
427	743	1411	6	0
421	743	1415	6	0
597	719	440	6	0
313	717	441	3	0
575	716	445	4	0
419	715	440	6	0
34	719	447	3	0
419	719	444	0	0
34	713	441	3	0
0	708	443	3	0
425	711	444	0	0
0	709	440	3	0
525	703	440	0	0
524	703	446	2	0
524	703	443	2	0
0	700	443	3	0
313	699	440	3	0
0	696	441	3	0
419	690	441	6	0
425	695	445	6	0
575	690	445	0	0
0	687	440	3	0
428	681	445	4	0
419	686	445	0	0
34	680	432	3	0
0	695	436	3	0
428	690	436	2	0
306	684	435	0	0
313	685	438	3	0
313	681	438	3	0
34	682	439	3	0
0	702	434	3	0
0	697	434	3	0
306	699	437	3	0
306	708	435	3	0
313	705	435	3	0
597	719	439	6	0
428	715	435	2	0
0	712	433	3	0
597	719	434	6	0
597	719	433	6	0
597	719	438	6	0
597	719	437	6	0
0	717	438	3	0
597	719	436	6	0
597	719	435	6	0
597	727	432	4	0
597	726	432	4	0
597	725	432	4	0
36	725	433	2	0
36	726	433	0	0
597	721	432	4	0
36	723	435	2	0
36	722	433	2	0
597	720	432	4	0
597	723	432	4	0
597	724	432	4	0
36	723	436	2	0
597	722	432	4	0
597	734	432	4	0
597	732	432	4	0
597	733	432	4	0
597	731	432	4	0
597	730	432	4	0
36	730	435	5	0
36	729	438	2	0
597	735	432	4	0
597	729	432	4	0
36	731	433	7	0
36	733	435	2	0
597	728	432	4	0
36	728	435	2	0
34	742	432	7	0
313	739	433	3	0
420	740	435	6	0
597	736	439	2	0
597	736	438	2	0
0	737	435	7	0
597	736	435	2	0
5	740	438	4	0
0	737	437	7	0
0	737	439	7	0
597	736	437	2	0
597	736	433	2	0
597	736	434	2	0
597	736	436	2	0
6	750	1382	4	0
6	745	1382	4	0
6	740	1382	4	0
421	750	1379	6	0
641	750	1380	0	0
421	745	1379	6	0
641	745	1380	0	0
421	740	1379	6	0
51	414	162	2	0
51	414	164	2	0
51	415	161	0	0
419	410	163	4	0
526	409	164	2	0
419	422	163	4	0
419	408	160	6	0
51	417	161	0	0
51	418	164	6	0
419	408	166	6	0
643	416	161	0	0
5	416	163	0	0
419	424	166	6	0
154	416	166	0	1
419	424	160	6	0
51	418	162	6	0
567	416	162	4	0
526	424	163	6	0
51	417	1106	0	0
51	415	1108	4	0
51	417	1108	4	0
51	415	1106	0	0
6	416	1107	4	0
568	416	1106	4	0
156	414	1107	1	1
155	419	1107	1	1
22	690	1379	0	0
427	682	1389	4	0
418	687	1389	0	0
418	691	1385	2	0
427	691	1380	2	0
22	691	1378	0	0
576	691	1389	0	0
578	696	1395	0	0
418	696	1392	2	0
578	704	1387	2	0
578	704	1390	2	0
418	696	1398	2	0
580	704	1384	0	0
426	711	1388	0	0
426	696	1389	6	0
578	708	1395	4	0
570	703	1394	4	0
578	710	1395	0	0
418	719	1388	0	0
580	704	1406	4	0
574	697	1396	2	0
426	711	1401	2	0
603	716	1386	0	0
574	711	1395	6	0
418	698	1395	4	0
418	711	1398	2	0
578	699	1395	0	0
600	697	1398	0	0
426	696	1402	4	0
418	709	1395	0	0
418	715	1384	2	0
418	711	1392	2	0
576	716	1401	0	0
578	713	1395	0	0
5	716	1393	0	0
576	715	1388	0	0
418	715	1392	2	0
574	716	1394	4	0
5	715	1394	4	0
418	716	1405	2	0
604	717	1404	0	0
605	712	1398	0	0
589	714	1388	4	0
418	716	1397	2	0
601	712	1393	0	0
580	716	1395	6	0
119	715	1377	4	0
427	715	1379	2	0
578	704	1400	2	0
578	704	1403	2	0
155	707	1395	1	1
156	701	1395	1	1
579	703	2333	2	0
579	703	2331	2	0
579	703	2343	2	0
581	704	2328	0	0
579	709	2338	0	0
571	703	2337	4	0
579	703	2346	2	0
579	711	2338	0	0
579	699	2338	0	0
579	696	2338	0	0
581	715	2339	6	0
579	698	2338	0	0
579	712	2338	0	0
3	716	2337	0	0
6	715	2337	6	0
581	703	2349	4	0
160	702	2338	1	1
158	705	2338	1	1
159	703	2340	0	1
157	703	2337	0	1
6	691	2339	4	0
581	691	2338	2	0
579	694	2338	0	0
3	690	2339	0	0
569	416	2050	4	0
585	416	2051	4	0
157	416	2050	0	1
159	416	2053	0	1
160	415	2051	1	1
158	418	2051	1	1
618	413	2999	6	0
617	420	2992	4	0
617	420	2991	4	0
584	418	2994	0	0
586	417	2993	0	0
51	617	3385	2	0
51	619	3389	6	0
51	619	3385	6	0
5	618	3383	4	0
433	600	3410	0	0
51	601	3408	0	0
7	343	1555	6	0
48	343	1547	0	0
3	343	1556	0	0
158	343	1557	2	0
5	341	1572	0	0
5	340	1575	6	0
3	346	1551	0	0
45	347	1572	0	0
6	349	1556	0	0
45	346	1572	0	0
3	344	1545	6	0
5	338	1570	0	0
3	345	1545	6	0
6	347	1571	0	0
3	344	1550	6	0
3	345	1551	0	0
121	346	1555	0	0
3	345	1550	6	0
3	344	1556	0	0
3	346	1550	6	0
3	346	1545	6	0
45	348	1572	0	0
6	349	1569	6	0
49	345	1571	1	1
1	346	1570	0	1
49	345	1573	1	1
49	345	1575	1	1
45	404	556	6	0
45	404	553	6	0
45	402	557	2	0
45	404	552	6	0
45	404	557	6	0
45	402	553	2	0
3	403	567	0	0
45	415	569	0	0
45	402	554	2	0
45	421	569	0	0
45	418	569	0	0
45	402	555	2	0
45	414	569	0	0
45	404	555	6	0
1	400	560	6	0
45	413	569	0	0
23	409	560	6	0
1	400	564	6	0
37	413	568	0	0
23	407	560	6	0
34	409	570	6	0
118	419	559	6	0
7	430	558	4	0
191	430	556	0	0
45	404	554	6	0
243	414	571	0	0
45	402	552	2	0
3	430	559	0	0
45	422	568	2	0
242	421	571	0	0
23	407	563	6	0
45	420	569	0	0
242	419	571	0	0
45	416	569	0	0
285	433	554	0	0
285	433	556	0	0
285	433	558	0	0
11	426	560	0	0
5	429	561	4	0
45	419	569	0	0
36	400	557	0	0
302	429	566	0	0
1	402	567	0	1
0	435	556	0	0
26	428	563	4	0
1	436	554	0	0
191	428	557	0	0
1	414	560	1	1
21	424	570	2	0
37	401	559	0	0
191	426	556	0	0
0	406	569	6	0
191	426	554	0	0
1	410	566	0	1
45	417	569	0	0
241	423	571	0	0
191	430	554	0	0
34	417	564	6	0
45	422	567	2	0
0	434	570	0	0
45	402	556	2	0
191	428	555	0	0
1	404	570	6	0
0	416	562	2	0
19	412	562	6	0
23	409	563	6	0
37	411	567	0	0
21	424	569	0	0
0	429	568	2	0
251	433	565	4	0
0	438	569	0	0
285	433	560	0	0
34	418	558	6	0
15	436	561	4	0
1	426	559	1	1
1	406	562	1	1
1	434	563	1	1
7	412	549	2	0
192	420	551	0	0
0	400	546	0	0
1	434	544	0	0
37	414	547	0	0
1	417	544	0	0
3	408	547	4	0
244	426	548	0	0
192	422	551	0	0
3	411	549	0	0
7	411	550	0	0
1	421	545	2	0
206	425	547	0	0
36	424	551	0	0
2	408	549	1	1
34	423	538	0	0
34	411	537	0	0
103	426	536	3	0
1	430	543	2	0
1	410	541	0	0
37	412	541	0	0
111	428	536	2	0
111	429	536	2	0
1	412	543	0	0
34	422	543	2	0
0	417	538	2	0
1	420	539	0	0
1	414	538	0	0
34	422	536	2	0
1	436	541	0	0
1	435	538	0	0
6	429	1505	4	0
15	426	1502	4	0
3	430	1502	4	0
36	421	532	2	0
36	423	531	2	0
37	409	533	0	0
0	410	530	0	0
36	416	528	2	0
34	416	531	0	0
34	419	535	0	0
0	413	534	0	0
111	430	534	2	0
113	430	533	3	0
103	428	532	0	0
1	427	534	2	0
103	428	533	0	0
1	434	533	0	0
51	415	3390	6	0
205	411	3378	6	0
51	415	3394	6	0
51	409	3399	0	0
205	415	3382	6	0
219	411	3398	2	0
205	410	3393	6	0
51	411	3395	4	0
51	409	3377	2	0
205	409	3380	6	0
51	413	3373	0	0
218	411	3399	2	0
205	414	3388	6	0
51	411	3388	0	0
205	414	3397	6	0
51	423	3381	2	0
205	423	3378	4	0
51	420	3379	4	0
65	406	3392	1	1
51	415	3402	6	0
219	410	3404	0	0
51	412	3405	4	0
219	410	3403	4	0
205	414	3404	6	0
205	414	3400	6	0
218	411	3400	2	0
205	410	3405	6	0
219	411	3403	0	0
245	412	3402	0	0
218	410	3401	0	0
205	409	3404	6	0
218	411	3402	2	0
51	408	3402	2	0
220	411	3401	2	0
205	408	3400	6	0
51	391	3369	6	0
15	271	1379	0	0
6	271	1377	2	0
147	278	1379	0	0
5	281	1388	0	0
148	275	1377	4	0
146	276	1377	0	0
15	274	1379	0	0
32	279	1383	0	1
309	559	654	0	0
23	583	661	4	0
286	571	661	0	0
98	590	657	6	0
98	587	658	3	0
25	582	661	4	0
21	582	663	4	0
98	588	657	5	0
51	588	659	2	0
1	566	660	0	0
51	590	659	6	0
98	588	656	4	0
25	585	661	4	0
98	591	658	5	0
1	576	657	2	0
98	590	656	3	0
9	589	658	0	1
9	589	663	0	1
1	568	667	0	0
1	573	667	0	0
51	580	664	2	0
286	567	664	0	0
51	591	667	4	0
19	588	667	4	0
51	587	667	4	0
7	585	668	0	0
9	583	666	0	1
0	588	679	0	0
360	578	678	4	0
3	566	683	0	0
51	572	685	0	0
51	570	685	0	0
1	588	682	0	0
0	584	683	4	0
1	586	685	0	0
3	570	680	0	0
1	591	685	0	0
0	582	682	0	0
55	568	684	0	0
0	591	681	4	0
1	571	685	0	1
1	573	681	1	1
1	565	685	0	1
1	571	682	0	1
23	594	661	4	0
8	599	661	4	0
281	592	663	0	0
25	593	661	4	0
51	598	664	6	0
71	595	666	6	0
7	593	668	0	0
5	593	664	0	0
0	593	678	4	0
21	596	663	4	0
25	596	661	4	0
7	592	662	3	0
0	599	683	4	0
0	594	682	4	0
0	593	686	0	0
9	594	667	1	1
44	558	1589	4	0
51	556	1593	2	0
14	557	1594	6	0
45	557	1590	6	0
45	559	1592	0	0
45	558	1592	0	0
45	557	1591	6	0
1	604	662	2	0
1	603	666	2	0
8	600	663	2	0
1	604	658	2	0
1	606	664	2	0
1	605	668	2	0
0	602	672	0	0
71	602	679	4	0
0	601	676	0	0
0	614	652	3	0
0	614	657	3	0
1	613	658	3	0
97	613	661	3	0
1	614	649	3	0
1	613	664	3	0
0	612	667	3	0
97	612	656	3	0
0	613	655	3	0
1	608	652	3	0
0	611	654	3	0
1	611	651	3	0
1	609	658	3	0
1	609	662	3	0
0	611	660	3	0
1	608	656	3	0
7	611	677	4	0
25	610	676	4	0
0	608	674	0	0
5	613	676	0	0
0	614	673	0	0
71	613	678	4	0
3	611	678	2	0
0	621	652	7	0
0	622	649	7	0
102	618	655	0	0
1	622	654	3	0
417	619	649	7	0
102	618	658	0	0
0	617	667	3	0
1	619	663	3	0
102	616	657	0	0
102	619	657	0	0
0	620	668	3	0
25	618	679	4	0
102	616	656	0	0
0	620	661	3	0
1	617	651	7	0
102	617	655	0	0
6	617	657	4	0
102	617	658	0	0
0	616	660	3	0
71	617	673	6	0
1	616	674	1	1
51	612	682	0	0
7	614	687	4	0
7	604	682	4	0
7	613	687	4	0
22	610	680	4	0
3	604	683	4	0
56	604	687	4	0
0	615	681	0	0
56	622	681	4	0
5	601	682	0	0
7	616	687	4	0
7	615	687	4	0
22	602	683	4	0
22	601	683	4	0
382	607	683	4	0
47	607	681	4	0
1	603	681	1	1
1	611	682	0	1
1	602	687	1	1
1	614	684	0	1
1	606	685	0	1
2	619	683	1	1
1	611	680	0	1
0	627	653	0	0
0	630	664	0	0
0	631	670	0	0
284	630	677	2	0
284	628	676	0	0
0	627	664	0	0
284	630	676	0	0
284	627	674	0	0
0	630	661	0	0
0	625	658	0	0
284	626	678	2	0
284	629	676	0	0
284	628	679	2	0
284	630	678	2	0
284	628	678	2	0
284	628	672	0	0
0	628	656	0	0
284	627	676	0	0
284	631	672	0	0
0	628	659	0	0
284	624	676	2	0
284	624	678	2	0
284	628	683	2	0
284	624	679	2	0
284	630	687	2	0
284	629	684	0	0
284	628	685	2	0
284	626	684	2	0
284	630	684	0	0
284	628	684	2	0
0	629	658	0	0
284	628	682	2	0
0	626	666	0	0
0	626	656	0	0
284	630	672	0	0
0	625	655	0	0
284	631	674	0	0
0	625	663	0	0
284	628	674	0	0
0	624	653	0	0
284	626	676	0	0
284	630	679	2	0
284	631	678	0	0
284	629	672	0	0
284	626	685	2	0
0	631	653	0	0
284	630	682	2	0
0	629	654	0	0
284	626	680	2	0
284	629	674	0	0
284	626	677	2	0
284	624	674	0	0
284	630	674	0	0
284	624	677	2	0
284	627	672	0	0
284	626	681	2	0
0	627	661	0	0
284	626	682	2	0
284	625	682	0	0
284	626	686	2	0
284	626	672	0	0
284	624	684	2	0
284	626	674	0	0
284	624	683	2	0
284	624	672	0	0
284	631	684	0	0
284	630	686	2	0
284	625	672	0	0
284	624	673	2	0
284	629	680	0	0
284	626	687	2	0
284	630	681	2	0
284	628	680	0	0
284	628	686	2	0
38	627	685	1	0
284	627	680	0	0
284	627	686	0	0
284	625	674	0	0
284	630	680	2	0
284	624	681	2	0
284	624	686	2	0
284	624	680	2	0
284	624	682	2	0
284	624	685	2	0
284	624	687	2	0
0	638	661	0	0
0	635	662	0	0
0	638	670	0	0
0	635	659	0	0
0	637	659	0	0
0	638	668	0	0
284	638	676	0	0
0	636	669	0	0
284	638	678	0	0
0	634	670	0	0
284	637	678	0	0
284	635	672	0	0
284	634	679	2	0
284	638	672	0	0
284	634	676	0	0
284	633	672	0	0
284	636	676	0	0
284	632	678	0	0
284	636	679	2	0
284	635	676	0	0
284	632	674	0	0
284	633	676	0	0
284	632	675	2	0
0	633	660	0	0
0	633	662	0	0
284	632	672	0	0
0	632	657	0	0
209	634	667	5	0
0	639	671	0	0
284	636	678	2	0
284	636	674	2	0
284	637	676	0	0
284	636	675	2	0
284	637	672	0	0
284	639	672	0	0
284	639	678	0	0
284	639	676	0	0
284	633	678	0	0
284	634	673	2	0
284	632	676	0	0
284	634	672	0	0
284	636	672	0	0
284	634	678	0	0
284	634	674	2	0
621	605	3509	4	0
51	602	3509	0	0
407	641	641	6	0
377	645	637	3	0
209	641	663	4	0
0	645	660	4	0
0	643	661	0	0
0	642	659	0	0
0	641	661	0	0
0	640	662	0	0
55	652	629	4	0
22	651	630	4	0
22	651	629	4	0
407	649	625	0	0
377	653	638	3	0
392	648	629	2	0
377	648	634	3	0
377	652	636	4	0
377	652	643	3	0
377	648	641	3	0
143	648	632	0	0
1	649	632	0	1
377	650	662	0	0
377	653	659	0	0
377	659	635	4	0
55	654	664	6	0
393	658	632	0	0
22	652	665	4	0
143	656	633	0	0
97	654	668	6	0
22	653	665	4	0
377	662	645	4	0
55	662	669	5	0
389	658	664	4	0
360	663	667	6	0
388	657	663	4	0
377	656	658	4	0
55	663	669	2	0
0	645	670	0	0
0	642	669	0	0
0	640	668	0	0
55	651	666	4	0
97	648	666	6	0
55	653	664	4	0
217	660	646	6	0
377	662	642	4	0
377	661	659	0	0
377	658	640	3	0
377	662	638	4	0
1	663	632	0	1
206	656	628	6	0
5	665	635	0	0
55	668	635	6	0
143	665	639	0	0
206	666	633	6	0
55	668	636	4	0
407	667	620	0	0
407	664	640	2	0
407	664	646	6	0
360	664	654	0	0
377	664	662	0	0
306	678	647	0	0
34	674	656	0	0
1	677	659	0	0
938	687	635	0	0
938	686	645	0	0
938	686	643	0	0
217	687	648	1	0
104	687	644	7	0
938	687	651	0	0
938	686	649	0	0
938	685	650	0	0
938	684	648	0	0
938	687	647	0	0
217	685	646	3	0
104	685	643	0	0
938	683	645	0	0
68	685	637	0	0
938	684	641	0	0
104	686	646	4	0
938	685	640	0	0
208	686	648	4	0
1	683	651	0	0
938	684	655	1	0
3	665	1581	6	0
409	662	1578	6	0
3	664	1577	6	0
6	665	1579	0	0
4	680	658	0	0
1	682	657	0	0
1	680	661	0	0
214	686	659	0	0
214	685	656	0	0
97	668	664	6	0
308	685	670	4	0
306	683	666	0	0
34	682	669	0	0
34	679	665	0	0
283	687	669	4	0
730	722	605	7	0
171	727	607	0	1
877	734	605	7	0
730	732	600	4	0
171	728	607	0	1
730	725	614	4	0
730	723	608	1	0
729	734	614	6	0
877	729	612	7	0
884	723	617	5	0
884	724	619	3	0
407	724	620	7	0
729	734	617	0	0
100	695	630	0	0
938	692	633	0	0
667	694	632	2	0
938	690	629	0	0
938	694	629	0	0
938	692	635	0	0
938	693	634	0	0
938	689	633	0	0
214	689	636	3	0
209	691	636	4	0
208	689	631	3	0
100	692	630	7	0
938	691	632	0	0
100	693	630	0	0
938	691	635	0	0
38	690	638	0	0
938	693	631	0	0
407	699	633	0	0
938	697	629	0	0
70	701	634	0	0
217	698	632	1	0
104	696	630	0	0
314	697	634	0	0
70	706	634	0	0
88	713	634	0	0
88	712	638	0	0
729	722	625	0	0
729	724	637	0	0
877	724	627	0	0
314	688	646	0	0
938	688	642	0	0
938	690	645	0	0
938	689	646	0	0
407	725	647	6	0
938	689	642	0	0
88	690	646	7	0
938	689	647	0	0
938	688	647	0	0
938	690	641	0	0
0	710	644	0	0
70	706	646	0	0
938	691	647	0	0
70	708	641	0	0
13	689	650	0	0
205	689	651	4	0
941	694	649	0	0
938	703	650	3	0
938	696	649	0	0
13	693	653	1	0
70	701	651	1	0
217	704	652	4	0
938	704	651	3	0
938	703	651	1	0
144	699	650	0	0
88	704	654	0	0
216	694	655	0	0
88	689	654	0	0
215	703	653	3	0
1	707	648	0	0
4	703	655	0	0
938	688	649	0	0
217	696	653	1	0
38	710	655	0	0
12	692	655	1	0
938	689	648	0	0
407	727	650	6	0
34	711	649	0	0
878	724	654	6	0
283	713	654	4	0
407	692	658	2	0
306	694	663	4	0
407	689	659	2	0
13	690	657	1	0
1	688	660	0	0
88	697	656	0	0
12	688	656	1	0
0	700	659	0	0
0	697	660	0	0
38	707	663	0	0
1	707	661	0	0
38	712	657	0	0
38	706	663	0	0
283	708	657	4	0
1	709	661	0	0
38	711	658	0	0
38	709	659	0	0
407	727	657	6	0
306	694	667	4	0
34	694	671	0	0
306	693	669	4	0
283	692	666	4	0
283	689	670	4	0
283	690	668	4	0
308	703	668	4	0
283	692	671	4	0
306	702	670	4	0
306	700	669	4	0
306	688	671	4	0
283	696	669	4	0
38	707	666	0	0
38	708	664	0	0
38	706	666	0	0
283	701	667	4	0
1	703	664	0	0
34	691	670	0	0
0	696	667	0	0
306	712	669	0	0
306	700	671	4	0
38	705	665	0	0
0	688	666	0	0
309	697	669	4	0
306	689	668	4	0
283	696	670	4	0
38	705	664	0	0
38	708	665	0	0
284	667	677	2	0
284	668	672	0	0
284	667	679	2	0
284	671	678	2	0
284	667	678	2	0
284	669	674	0	0
284	668	674	0	0
284	669	672	0	0
284	671	672	0	0
284	667	674	0	0
284	671	679	2	0
284	670	672	0	0
284	669	675	2	0
284	671	673	2	0
284	666	674	0	0
284	671	674	2	0
284	667	672	0	0
284	664	676	0	0
284	665	672	0	0
284	665	677	2	0
284	665	676	0	0
284	669	676	2	0
284	666	672	0	0
284	665	679	2	0
284	667	676	2	0
284	669	678	2	0
284	669	679	2	0
284	671	675	2	0
284	665	674	0	0
284	669	677	2	0
284	671	677	2	0
0	695	674	0	0
284	665	678	2	0
284	671	676	2	0
55	695	679	2	0
1	700	678	2	0
8	695	676	2	0
0	690	674	0	0
38	699	678	2	0
1	681	673	0	0
1	702	677	2	0
284	664	672	0	0
0	673	678	0	0
8	694	675	2	0
1	694	673	0	0
1	696	674	0	0
4	692	675	0	0
0	682	672	0	0
0	682	676	0	0
34	699	675	2	0
38	691	676	2	0
55	696	679	1	0
284	664	674	0	0
1	695	677	2	1
34	710	678	7	0
1	706	678	0	0
0	705	678	0	0
306	704	675	0	0
283	707	679	0	0
283	708	679	0	0
283	709	678	0	0
283	708	677	0	0
0	708	675	1	0
1	718	678	1	0
34	719	679	1	0
47	716	679	4	0
283	716	677	4	0
1	718	679	1	0
928	712	679	4	0
7	715	679	4	0
0	694	687	0	0
97	692	680	4	0
38	691	687	2	0
34	703	685	2	0
0	688	685	2	0
1	703	680	0	0
283	704	682	0	0
283	717	685	0	0
0	719	682	1	0
0	702	684	2	0
1	718	682	1	0
34	718	680	1	0
34	719	681	1	0
47	716	683	4	0
1	719	680	1	0
34	707	683	2	0
25	711	684	0	0
286	709	683	0	0
11	711	681	0	0
1	690	682	0	0
37	712	687	4	0
0	718	681	1	0
47	716	681	4	0
119	691	682	0	0
0	689	681	0	0
97	700	680	4	0
283	716	686	0	0
38	689	686	2	0
2	713	686	4	0
283	705	681	0	0
278	713	681	0	0
37	714	685	4	0
37	715	686	4	0
34	708	682	2	0
38	699	681	2	0
0	705	683	0	0
37	711	686	4	0
34	707	681	2	0
1	719	683	1	0
4	718	683	1	0
34	718	684	1	0
34	718	686	7	0
204	700	686	4	0
1	717	683	0	0
0	717	680	0	0
34	716	687	2	0
37	714	687	4	0
7	715	684	0	0
3	715	683	0	0
3	715	680	0	0
37	716	685	4	0
1	698	685	1	1
1	711	683	1	1
1	695	683	1	1
283	695	689	0	0
38	701	688	2	0
34	697	688	2	0
306	696	690	0	0
283	710	692	0	0
51	714	693	0	0
1	717	689	0	0
915	699	695	2	0
4	696	688	0	0
283	714	689	0	0
933	713	695	0	0
51	712	693	0	0
0	697	693	0	0
283	715	690	0	0
34	715	688	2	0
46	712	694	0	0
46	714	694	0	0
1	685	682	0	0
38	684	692	2	0
34	686	694	0	0
34	685	695	0	0
0	684	695	0	0
4	680	694	0	0
1	682	693	0	0
4	675	689	0	0
284	671	694	2	0
284	671	695	2	0
284	665	681	2	0
284	665	682	2	0
284	664	686	0	0
284	671	686	2	0
284	667	682	2	0
284	667	685	2	0
284	665	684	2	0
284	671	685	2	0
284	667	683	2	0
284	667	684	2	0
284	669	682	2	0
284	669	684	2	0
284	671	684	2	0
284	669	683	2	0
284	669	681	2	0
284	671	692	2	0
284	671	690	2	0
284	671	683	2	0
284	665	680	2	0
284	664	684	0	0
284	671	691	2	0
284	665	686	0	0
284	671	687	2	0
284	671	689	2	0
284	669	680	2	0
284	670	686	0	0
284	668	680	0	0
284	667	686	0	0
284	668	686	0	0
284	667	680	2	0
284	669	686	0	0
284	671	681	2	0
284	671	693	2	0
284	665	683	2	0
284	666	688	0	0
284	666	682	0	0
284	671	682	2	0
284	671	680	2	0
284	671	688	2	0
284	667	695	2	0
284	669	695	2	0
284	664	690	0	0
284	669	688	0	0
284	668	688	0	0
284	669	693	2	0
284	669	694	2	0
284	667	690	0	0
284	667	688	0	0
284	669	690	2	0
284	669	689	2	0
284	666	690	0	0
284	669	692	2	0
284	669	691	2	0
284	667	694	2	0
284	667	692	2	0
284	667	693	2	0
284	667	691	2	0
284	665	688	0	0
284	665	690	0	0
284	664	688	0	0
284	671	703	2	0
284	669	698	2	0
284	667	700	0	0
284	664	700	0	0
284	669	702	0	0
284	667	702	0	0
284	669	699	2	0
284	671	698	2	0
284	671	701	2	0
284	671	702	2	0
36	688	703	0	0
284	667	698	2	0
38	687	701	2	0
284	667	699	2	0
192	694	702	0	0
283	680	698	0	0
284	667	696	2	0
34	682	700	0	0
36	694	701	0	0
284	669	697	2	0
284	669	696	2	0
284	666	700	0	0
192	696	703	4	0
34	683	701	2	0
283	685	701	0	0
34	693	697	2	0
34	684	703	2	0
307	688	702	0	0
34	694	698	2	0
0	696	696	0	0
284	664	702	0	0
284	666	702	0	0
36	693	700	0	0
284	667	697	2	0
284	671	700	2	0
34	686	700	0	0
1	675	703	2	0
284	669	700	2	0
284	669	701	2	0
0	679	696	0	0
1	684	698	0	0
284	665	702	0	0
284	668	702	0	0
284	671	696	2	0
284	671	697	2	0
284	671	699	2	0
284	665	700	0	0
938	710	703	4	0
51	709	698	2	0
27	709	697	2	0
46	710	698	0	0
46	710	696	0	0
283	705	697	0	0
51	709	696	2	0
46	712	700	0	0
51	717	698	6	0
51	717	696	6	0
46	714	700	0	0
46	716	698	0	0
925	713	699	0	0
46	716	696	0	0
27	717	697	6	0
938	718	701	4	0
51	714	701	4	0
51	712	701	4	0
34	685	706	0	0
283	682	711	0	0
0	681	711	0	0
34	682	709	0	0
34	681	706	2	0
283	711	707	0	0
36	693	710	2	0
36	689	708	0	0
306	701	704	0	0
36	697	706	0	0
307	697	708	0	0
34	704	707	2	0
0	708	706	2	0
34	690	710	2	0
1	707	707	2	0
36	692	709	0	0
283	706	705	0	0
283	680	709	0	0
938	715	704	0	0
23	697	704	2	0
38	701	708	2	0
34	700	707	2	0
1	680	705	2	0
36	696	708	0	0
1	701	706	2	0
283	714	706	0	0
192	689	706	0	0
34	708	710	2	0
283	712	705	0	0
938	713	706	0	0
283	711	705	0	0
283	716	704	0	0
306	714	710	0	0
34	685	713	2	0
0	703	712	2	0
34	703	713	2	0
1	681	713	0	0
38	699	714	2	0
0	687	717	0	0
1	683	712	0	0
1	709	715	0	0
1	701	713	2	0
1	717	712	0	0
38	686	712	2	0
0	679	704	2	0
1	678	708	0	0
283	679	705	0	0
283	675	707	0	0
0	673	707	2	0
1	675	716	0	0
284	663	676	0	0
284	663	672	0	0
284	663	674	0	0
284	657	674	0	0
284	657	675	2	0
284	661	678	2	0
284	663	677	2	0
284	659	678	0	0
284	659	674	0	0
284	663	684	0	0
284	659	672	0	0
284	659	676	0	0
284	660	672	0	0
284	657	676	0	0
284	663	682	2	0
284	663	681	2	0
284	657	677	2	0
284	661	679	2	0
284	658	678	0	0
284	663	680	0	0
284	661	676	2	0
284	661	684	0	0
284	662	674	0	0
284	658	680	0	0
284	659	684	0	0
284	659	686	0	0
284	657	686	2	0
284	659	682	0	0
284	659	680	0	0
284	662	672	0	0
284	658	684	0	0
284	661	686	0	0
284	663	678	2	0
284	660	690	0	0
284	660	676	0	0
284	657	687	2	0
284	656	690	0	0
284	663	690	0	0
284	661	688	0	0
284	661	683	2	0
284	661	682	0	0
284	657	681	2	0
284	657	680	0	0
284	656	682	0	0
284	663	686	0	0
284	657	690	0	0
284	660	688	0	0
284	658	674	0	0
284	661	674	0	0
390	660	695	5	0
284	661	675	2	0
284	657	672	0	0
284	660	674	0	0
284	661	672	0	0
284	662	684	0	0
284	657	688	0	0
284	656	676	0	0
284	656	672	0	0
284	661	685	2	0
284	658	672	0	0
284	661	677	2	0
284	661	690	0	0
284	657	678	0	0
284	658	688	0	0
284	657	684	2	0
284	657	683	2	0
284	659	690	0	0
284	663	688	0	0
284	659	688	0	0
284	657	682	2	0
284	658	690	0	0
284	662	688	0	0
284	662	690	0	0
284	662	680	0	0
284	660	680	0	0
284	660	682	0	0
284	660	686	0	0
284	657	685	2	0
284	662	686	0	0
284	661	680	0	0
284	654	678	0	0
284	649	676	0	0
284	653	672	0	0
284	652	672	0	0
284	650	672	0	0
284	650	676	0	0
284	655	674	2	0
284	649	678	0	0
37	654	673	5	0
284	652	678	0	0
284	650	678	0	0
284	652	676	0	0
284	651	672	0	0
284	651	676	0	0
284	653	678	0	0
284	654	672	0	0
284	653	676	0	0
284	648	672	0	0
284	655	678	2	0
284	649	672	0	0
284	648	678	0	0
284	654	676	0	0
284	651	678	0	0
284	655	672	0	0
284	655	673	2	0
284	655	679	2	0
284	655	676	0	0
284	648	676	0	0
37	646	675	3	0
284	647	672	0	0
23	647	673	4	0
284	647	678	0	0
284	646	678	0	0
284	643	679	2	0
284	646	676	0	0
284	646	672	0	0
37	644	675	4	0
284	644	672	0	0
284	642	676	0	0
284	647	676	0	0
284	642	672	0	0
284	643	678	2	0
23	642	673	4	0
284	641	679	2	0
284	640	672	0	0
284	640	678	0	0
284	643	676	0	0
37	643	675	5	0
284	643	677	2	0
284	645	678	2	0
284	645	676	0	0
284	640	676	0	0
284	644	676	0	0
284	643	672	0	0
284	641	676	0	0
284	645	672	0	0
284	641	678	0	0
284	645	679	2	0
284	641	672	0	0
284	655	684	0	0
284	653	686	0	0
284	652	686	0	0
284	655	686	0	0
284	653	684	0	0
284	652	682	2	0
284	650	686	0	0
284	649	685	2	0
284	650	680	0	0
284	649	686	0	0
284	648	686	0	0
284	648	680	0	0
284	652	681	2	0
284	652	680	0	0
284	655	683	2	0
284	654	684	0	0
284	654	686	0	0
284	655	682	0	0
284	655	687	2	0
284	655	680	2	0
284	649	683	2	0
284	651	680	0	0
284	652	684	0	0
284	652	683	2	0
284	649	684	2	0
284	651	686	0	0
284	649	680	0	0
284	649	682	2	0
284	647	686	0	0
284	647	682	2	0
284	647	684	0	0
284	647	681	2	0
284	647	683	2	0
284	647	680	2	0
284	645	681	0	0
284	646	686	0	0
284	646	684	0	0
284	644	684	0	0
284	641	680	2	0
284	642	686	0	0
37	642	683	5	0
284	644	681	0	0
284	642	681	0	0
284	644	686	0	0
284	641	686	0	0
37	645	683	3	0
284	645	686	0	0
284	645	684	0	0
284	641	684	0	0
284	643	684	0	0
284	645	680	2	0
284	642	684	0	0
284	643	686	0	0
284	640	686	0	0
284	643	681	0	0
284	641	681	2	0
284	640	684	0	0
284	636	681	2	0
284	638	680	2	0
284	635	686	0	0
284	639	686	0	0
284	635	682	0	0
37	633	687	0	0
284	633	686	0	0
284	634	680	2	0
284	638	682	2	0
284	636	680	2	0
284	637	686	0	0
284	638	681	2	0
284	632	680	2	0
284	633	682	0	0
284	632	683	2	0
284	632	682	2	0
284	632	681	2	0
284	634	684	0	0
284	632	685	2	0
284	632	686	2	0
284	639	684	0	0
284	632	687	2	0
284	636	684	0	0
284	638	683	2	0
284	635	684	0	0
284	634	686	0	0
284	638	684	0	0
284	636	682	2	0
284	634	682	0	0
284	632	684	2	0
284	636	686	0	0
284	638	686	0	0
284	637	684	0	0
7	603	690	2	0
56	606	690	2	0
3	602	690	0	0
70	602	693	4	0
7	602	689	4	0
7	614	689	0	0
284	624	693	2	0
284	624	689	2	0
284	624	691	2	0
284	624	692	2	0
284	626	688	2	0
284	624	694	2	0
284	624	688	2	0
284	624	695	2	0
284	628	692	2	0
284	625	694	0	0
284	624	690	2	0
284	626	694	0	0
284	626	692	2	0
284	627	692	0	0
284	630	693	2	0
284	626	691	2	0
284	630	691	2	0
284	630	694	2	0
284	630	692	2	0
70	621	690	4	0
47	610	688	0	0
51	615	691	0	0
7	613	689	0	0
54	608	692	1	0
284	626	689	2	0
284	626	690	2	0
284	630	688	2	0
284	630	690	2	0
284	628	693	2	0
284	630	695	2	0
284	629	689	0	0
284	628	694	2	0
284	628	689	0	0
284	628	695	2	0
55	610	691	0	0
284	630	689	2	0
284	632	695	2	0
284	632	689	2	0
428	636	690	2	0
51	613	691	0	0
58	614	694	0	0
56	618	689	4	0
9	613	688	0	0
1	633	695	0	0
284	632	691	2	0
56	612	690	2	0
56	610	690	2	0
55	611	691	0	0
7	616	689	0	0
284	632	690	2	0
284	632	688	2	0
1	605	689	1	1
284	632	693	2	0
284	632	694	2	0
420	639	691	0	0
1	605	691	0	1
0	638	695	0	0
284	632	692	2	0
1	614	691	0	1
7	559	689	0	0
98	556	689	0	0
51	562	688	2	0
3	559	688	0	0
3	558	693	0	0
0	588	695	4	0
104	556	691	0	0
51	562	690	2	0
101	555	695	0	0
107	557	690	0	0
1	560	691	0	1
1	561	692	1	1
100	553	698	0	0
3	591	702	4	0
51	570	697	4	0
7	591	701	4	0
51	568	697	4	0
99	554	697	6	0
3	562	699	0	0
50	558	703	0	0
50	560	703	0	0
50	568	699	2	0
51	558	700	4	0
25	589	701	4	0
3	572	700	0	0
1	560	701	0	1
1	569	698	0	1
1	564	698	0	1
1	590	703	1	1
1	563	702	1	1
3	598	689	4	0
51	596	691	0	0
56	599	689	4	0
51	598	691	0	0
377	598	699	5	0
7	593	702	2	0
23	596	703	6	0
3	592	702	4	0
3	593	697	4	0
71	594	696	4	0
7	593	696	4	0
2	596	688	0	1
2	599	688	0	1
1	595	698	1	1
1	594	703	1	1
1	597	691	0	1
23	601	696	4	0
503	608	696	0	0
371	621	700	0	0
51	622	702	6	0
284	629	702	0	0
284	626	696	0	0
377	617	699	3	0
284	630	700	0	0
284	630	697	2	0
284	626	699	2	0
284	624	696	2	0
284	628	698	0	0
284	630	696	2	0
284	628	699	2	0
284	632	700	2	0
37	629	703	1	0
284	632	701	2	0
284	630	702	2	0
284	628	702	0	0
284	628	701	2	0
20	620	702	4	0
284	626	701	2	0
284	626	702	2	0
37	635	703	0	0
284	632	697	2	0
284	632	696	2	0
5	636	696	0	0
113	621	699	0	1
284	627	696	0	0
284	624	697	2	0
37	637	701	0	0
284	626	698	2	0
284	628	696	2	0
284	624	698	2	0
284	626	697	2	0
284	630	703	2	0
284	632	702	2	0
0	634	698	0	0
284	632	703	2	0
37	633	700	0	0
284	632	698	2	0
284	632	699	2	0
284	630	698	2	0
284	631	700	0	0
284	624	701	2	0
284	628	700	2	0
284	624	700	2	0
284	629	698	0	0
284	624	703	2	0
284	626	700	2	0
284	626	703	2	0
420	636	699	2	0
284	624	702	2	0
284	624	699	2	0
51	622	708	6	0
283	623	710	0	0
371	621	707	0	0
51	622	706	6	0
51	622	704	6	0
20	620	705	4	0
284	630	705	2	0
284	630	704	2	0
377	617	709	0	0
284	628	708	2	0
284	628	707	2	0
284	628	709	2	0
284	630	707	0	0
284	629	711	0	0
284	628	706	2	0
284	626	708	2	0
284	624	708	2	0
284	624	706	2	0
284	626	705	2	0
284	630	709	0	0
284	631	709	0	0
284	631	707	0	0
284	630	711	0	0
284	624	707	2	0
284	626	707	2	0
284	628	705	2	0
284	626	706	2	0
284	629	705	0	0
284	628	711	0	0
284	627	711	0	0
284	629	709	0	0
284	628	710	2	0
284	634	709	0	0
284	638	710	2	0
284	636	711	0	0
284	624	709	2	0
284	626	711	0	0
284	626	709	2	0
284	634	707	0	0
284	624	710	2	0
284	626	704	2	0
284	632	706	2	0
284	632	711	0	0
284	632	709	0	0
284	634	705	0	0
284	636	709	0	0
284	636	707	0	0
284	636	708	2	0
284	624	705	2	0
284	638	706	2	0
284	624	711	2	0
284	638	709	0	0
284	634	710	2	0
284	639	707	0	0
284	638	707	2	0
284	624	704	2	0
284	633	707	0	0
284	634	711	0	0
284	638	705	2	0
284	633	709	0	0
284	633	711	0	0
284	638	711	2	0
284	639	709	0	0
284	635	705	0	0
284	637	705	0	0
284	635	711	0	0
284	637	709	0	0
1	637	704	0	0
284	632	705	2	0
284	632	707	0	0
284	636	705	0	0
284	635	707	0	0
284	632	704	2	0
115	615	711	2	1
101	634	705	0	1
101	633	705	0	1
114	619	711	1	1
1	622	711	0	1
17	613	1623	0	0
18	595	1611	0	0
51	597	1608	6	0
18	594	1611	0	0
6	601	1626	0	0
71	610	1621	0	0
17	602	1625	2	0
51	597	1609	6	0
5	595	1608	0	0
25	596	1607	0	0
6	613	1620	0	0
6	593	1608	0	0
51	592	1608	2	0
382	600	1623	0	0
1	592	1609	1	1
3	584	1610	0	0
51	586	1608	6	0
51	591	1608	0	0
51	587	1608	0	0
1	587	1609	1	1
23	596	709	6	0
23	596	706	6	0
1	593	709	0	1
7	589	707	6	0
7	591	711	5	0
25	589	711	4	0
7	590	706	4	0
3	590	707	4	0
0	585	706	0	0
56	588	704	2	0
0	584	711	0	0
105	573	708	0	0
103	569	708	0	0
105	574	710	0	0
107	569	710	0	0
103	569	706	0	0
103	571	710	0	0
103	573	706	0	0
105	571	708	0	0
16	590	717	0	0
16	590	718	0	0
3	590	715	0	0
377	598	712	7	0
16	590	713	0	0
16	590	716	0	0
3	590	714	0	0
7	594	718	0	0
7	592	716	0	0
3	592	712	1	0
3	592	715	0	0
20	607	717	0	0
51	612	718	4	0
51	608	718	4	0
20	611	717	0	0
7	593	713	1	0
3	593	718	0	0
3	594	717	0	0
51	606	718	4	0
3	594	716	0	0
51	610	718	4	0
377	613	713	2	0
371	604	716	2	0
58	600	717	0	0
114	615	715	0	1
1	615	718	1	1
1	594	713	3	1
113	603	717	1	1
115	568	713	0	0
101	569	712	0	0
102	569	715	0	0
105	574	716	0	0
105	572	713	0	0
104	571	715	0	0
104	565	717	0	0
107	561	705	0	0
102	567	716	0	0
114	565	715	0	0
102	560	714	0	0
114	562	712	0	0
100	562	716	0	0
104	561	713	0	0
101	565	704	0	0
104	562	714	0	0
100	559	705	0	0
51	552	705	4	0
102	559	717	0	0
51	554	705	4	0
5	553	709	0	0
1	553	706	0	1
45	547	704	4	0
45	546	704	4	0
45	545	704	4	0
45	548	704	4	0
45	544	704	4	0
21	545	705	6	0
1106	537	703	0	0
385	538	699	4	0
386	540	699	4	0
1106	536	703	0	0
1106	536	702	0	0
1181	543	707	4	0
45	543	704	4	0
1102	540	709	6	0
45	539	704	4	0
45	540	704	4	0
1101	540	707	6	0
1180	541	707	4	0
45	542	704	4	0
45	538	704	4	0
45	536	704	4	0
45	541	704	4	0
45	537	704	4	0
21	541	705	6	0
21	537	705	6	0
1179	529	707	4	0
21	534	702	6	0
45	535	704	4	0
21	534	704	4	0
45	535	703	2	0
384	534	699	4	0
45	535	702	2	0
395	572	723	0	0
395	558	723	0	0
397	568	727	0	0
397	563	723	0	0
1	573	727	0	0
1	562	727	0	0
7	551	1652	0	0
3	552	1651	0	0
7	553	1651	2	0
6	553	1653	0	0
402	579	727	0	0
1	591	722	1	0
34	584	726	0	0
397	569	734	0	0
308	570	731	0	0
1	573	735	0	0
394	564	735	0	0
1	566	731	0	0
1	584	733	1	0
1	582	729	1	0
1	590	732	1	0
394	595	731	0	0
1	598	723	1	0
402	593	726	0	0
395	598	727	2	0
397	573	743	0	0
394	570	738	0	0
402	592	741	1	0
34	577	739	0	0
397	565	743	0	0
1	596	736	1	0
402	584	741	0	0
395	586	739	0	0
402	599	741	1	0
116	607	724	6	0
146	602	727	6	0
117	600	723	2	0
45	605	721	6	0
45	605	722	6	0
7	601	742	4	0
146	602	726	2	0
34	600	737	0	0
45	602	722	2	0
44	603	722	0	0
45	605	725	0	0
3	601	743	0	0
45	605	724	6	0
5	605	742	4	0
117	606	726	4	0
45	602	725	0	0
45	602	724	2	0
968	607	742	4	0
45	602	721	2	0
45	602	723	2	0
45	605	723	6	0
34	605	729	0	0
306	600	728	2	0
283	607	743	4	0
24	607	722	3	1
24	601	721	1	1
24	607	721	1	1
24	601	722	3	1
24	606	722	2	1
24	600	722	2	1
34	609	739	0	0
402	609	741	1	0
34	609	732	0	0
1	608	743	3	0
1	609	725	1	0
34	614	730	0	0
395	609	721	2	0
306	609	727	2	0
51	622	714	6	0
395	620	739	0	0
51	622	712	6	0
378	616	712	1	0
41	621	714	0	0
1	621	721	1	0
51	620	718	4	0
51	622	716	6	0
51	617	714	4	0
51	618	718	4	0
51	616	718	4	0
51	618	713	6	0
284	628	713	2	0
284	626	715	2	0
284	630	712	2	0
284	630	713	0	0
284	625	719	0	0
284	630	715	0	0
284	627	717	0	0
284	626	713	0	0
284	626	716	2	0
284	629	719	0	0
284	628	715	2	0
284	631	719	0	0
284	628	714	2	0
284	631	717	0	0
284	627	713	0	0
284	628	717	0	0
284	631	713	0	0
284	626	717	0	0
284	630	717	0	0
284	624	718	2	0
284	624	719	0	0
284	630	719	0	0
284	624	716	2	0
34	627	723	0	0
0	628	723	0	0
284	629	715	0	0
284	624	715	2	0
284	624	714	2	0
284	625	713	0	0
284	626	719	0	0
0	626	725	0	0
960	631	721	0	0
960	630	720	0	0
284	628	719	0	0
284	624	717	2	0
0	628	721	0	0
308	631	725	0	0
1	628	734	6	0
960	630	722	0	0
284	628	716	2	0
960	627	729	0	0
960	627	732	0	0
34	628	730	0	0
284	627	719	0	0
285	628	741	0	0
284	627	741	0	0
960	628	725	0	0
284	624	713	2	0
37	629	726	0	0
26	627	739	0	0
37	625	727	0	0
0	628	727	0	0
1	629	724	6	0
960	631	727	0	0
285	631	729	2	0
993	629	720	0	0
960	626	721	0	0
992	629	730	0	0
960	630	733	0	0
960	625	723	0	0
284	624	712	2	0
960	625	720	0	0
960	629	736	0	0
4	625	731	6	0
285	624	741	0	0
284	629	741	0	0
1	631	731	6	0
1029	624	742	4	0
285	626	741	0	0
23	630	741	0	0
960	625	735	0	0
284	625	741	0	0
284	638	715	0	0
284	638	717	0	0
284	639	715	0	0
284	639	713	0	0
284	637	713	0	0
284	636	712	2	0
284	633	717	0	0
284	636	713	0	0
284	632	713	0	0
284	632	714	2	0
284	632	715	2	0
284	636	719	0	0
284	639	727	2	0
284	633	719	0	0
284	637	715	0	0
284	635	719	0	0
284	634	715	0	0
284	634	713	0	0
284	639	719	0	0
284	636	717	0	0
37	634	727	0	0
284	636	715	0	0
284	638	719	0	0
37	634	723	2	0
284	639	717	0	0
2	637	727	2	0
961	638	722	0	0
960	636	724	0	0
284	639	721	2	0
960	638	725	0	0
1	637	723	6	0
284	639	722	2	0
284	632	717	0	0
284	633	713	0	0
284	635	715	0	0
960	638	733	0	0
960	632	720	0	0
284	639	732	2	0
284	639	734	2	0
960	637	720	0	0
284	639	735	2	0
285	634	721	2	0
0	633	726	0	0
284	634	717	0	0
284	639	728	2	0
991	638	728	0	0
284	639	724	2	0
284	635	717	0	0
308	638	720	0	0
284	634	719	0	0
658	637	735	0	0
284	637	719	0	0
284	639	725	2	0
284	637	717	0	0
284	639	726	2	0
23	636	736	4	0
284	632	719	0	0
659	638	741	0	0
960	638	721	0	0
284	632	716	2	0
284	639	723	2	0
37	637	726	2	0
1	632	735	6	0
284	638	713	0	0
284	639	720	2	0
37	633	729	2	0
960	635	722	0	0
990	633	724	0	0
284	639	733	2	0
37	635	725	2	0
34	636	729	0	0
960	633	721	0	0
960	637	729	0	0
960	637	731	0	0
0	635	728	0	0
0	632	723	0	0
284	639	729	2	0
960	635	720	0	0
284	639	730	2	0
37	636	721	2	0
960	634	733	0	0
284	639	731	2	0
34	632	727	0	0
284	633	741	0	0
309	633	731	0	0
285	632	741	0	0
3	639	737	0	0
0	635	730	0	0
960	636	732	0	0
1	634	739	1	1
37	608	745	5	0
48	629	746	0	0
3	633	747	0	0
3	616	749	0	0
0	609	744	3	0
7	601	744	0	0
3	615	749	0	0
3	614	749	0	0
7	633	748	0	0
11	629	749	0	0
7	633	746	4	0
1	619	750	1	1
1	603	746	0	1
1	631	751	0	1
1	612	750	1	1
1	599	745	0	0
1	595	745	0	0
20	585	751	2	0
20	583	744	0	0
402	577	746	0	0
20	580	747	0	0
1	571	744	0	0
21	552	749	0	0
209	555	746	1	0
21	552	746	4	0
394	559	744	0	0
21	550	746	4	0
21	546	749	0	0
21	548	746	4	0
21	550	749	0	0
21	548	749	0	0
21	546	746	4	0
51	605	3552	0	0
143	604	3572	0	0
51	585	3573	6	0
143	610	3570	4	0
213	584	3553	4	0
614	601	3560	2	0
205	598	3536	2	0
615	601	3562	2	0
615	601	3561	2	0
51	599	3565	2	0
143	593	3574	2	0
51	584	3564	6	0
657	607	3568	6	0
55	606	3555	0	0
614	601	3558	2	0
51	599	3537	4	0
143	610	3565	0	0
143	604	3575	4	0
143	617	3568	7	0
51	602	3565	6	0
51	584	3561	6	0
645	617	3567	2	0
43	603	3554	0	0
614	601	3559	2	0
51	602	3552	0	0
55	601	3554	0	0
51	585	3570	6	0
55	606	3554	0	0
143	599	3567	2	0
34	602	3536	2	0
143	615	3564	0	0
3	614	3564	2	0
51	584	3558	6	0
143	616	3570	4	0
625	609	3573	2	0
656	606	3568	2	0
598	599	3582	4	0
25	584	3582	0	0
143	593	3576	2	0
598	595	3582	4	0
598	597	3582	4	0
51	609	3583	0	0
51	606	3583	0	0
627	598	3582	4	0
598	596	3582	4	0
51	591	3587	4	0
51	606	3589	4	0
598	597	3584	4	0
55	595	3590	6	0
25	584	3589	0	0
55	594	3590	6	0
628	596	3584	4	0
51	596	3591	6	0
162	593	3590	0	1
51	603	3591	2	0
598	595	3584	4	0
51	595	3587	4	0
598	598	3584	4	0
598	599	3584	4	0
51	591	3585	0	0
150	604	3587	0	1
51	609	3589	4	0
5	610	3588	2	0
51	603	3584	2	0
5	605	3591	2	0
51	594	3585	0	0
51	591	3591	2	0
1007	628	3556	0	0
51	603	3588	2	0
205	630	3567	4	0
1007	631	3562	0	0
1007	626	3559	0	0
205	625	3567	4	0
205	631	3569	4	0
117	631	3575	4	0
117	625	3566	2	0
164	631	3568	6	0
205	628	3570	4	0
38	630	3569	4	0
38	627	3574	6	0
164	627	3569	6	0
205	626	3575	4	0
51	631	3580	0	0
58	630	3581	0	0
116	627	3582	2	0
205	627	3589	0	0
117	627	3580	2	0
191	629	3573	0	1
187	626	3591	1	1
190	627	3563	3	1
43	591	3593	6	0
22	596	3593	6	0
22	595	3597	6	0
51	591	3596	2	0
51	596	3596	6	0
51	582	3564	2	0
51	581	3570	2	0
103	581	3572	2	0
25	581	3589	0	0
213	582	3585	0	0
143	576	3583	2	0
103	581	3575	2	0
25	581	3582	0	0
143	576	3581	2	0
630	576	3577	0	0
51	582	3561	2	0
636	581	3573	2	0
205	599	3534	2	0
629	576	3529	0	0
51	597	3532	0	0
143	581	3527	6	0
633	581	3525	2	0
51	595	3534	2	0
277	577	3526	4	0
143	576	3527	2	0
143	579	3524	0	0
143	579	3531	4	0
277	579	3527	4	0
51	607	3516	6	0
286	606	3513	4	0
34	607	3518	4	0
51	607	3512	6	0
34	604	3512	4	0
51	605	3519	4	0
286	603	3512	4	0
277	604	3524	4	0
402	607	3525	4	0
277	607	3527	2	0
286	601	3514	4	0
51	600	3512	2	0
51	600	3516	2	0
277	604	3520	4	0
51	602	3523	2	0
51	611	3525	6	0
402	610	3524	4	0
277	610	3526	4	0
51	608	3522	0	0
277	610	3523	4	0
51	607	3529	2	0
277	607	3533	2	0
51	609	3530	6	0
277	605	3535	2	0
51	602	3521	2	0
284	647	707	0	0
284	645	712	2	0
306	643	738	0	0
284	646	709	0	0
398	643	732	7	0
34	646	732	0	0
1	644	724	1	0
0	642	725	1	0
284	643	713	2	0
284	641	712	2	0
284	641	717	0	0
284	643	714	2	0
309	647	723	0	0
284	643	717	0	0
284	641	713	0	0
284	643	715	0	0
284	645	717	0	0
284	641	715	0	0
284	640	715	0	0
284	640	713	0	0
284	645	713	0	0
284	643	707	0	0
284	645	715	0	0
284	643	709	0	0
284	643	711	0	0
284	647	719	0	0
284	640	707	0	0
284	642	711	0	0
284	640	709	0	0
284	647	715	0	0
284	644	711	0	0
284	642	709	0	0
284	646	707	0	0
284	645	719	0	0
284	644	709	0	0
284	640	719	0	0
284	642	719	0	0
284	644	707	0	0
284	640	717	0	0
284	647	709	0	0
284	644	719	0	0
284	642	715	0	0
284	644	715	0	0
284	643	719	0	0
284	642	717	0	0
284	641	719	0	0
284	646	715	0	0
0	645	704	0	0
284	645	709	0	0
284	644	717	0	0
284	642	707	0	0
284	645	711	0	0
284	645	707	0	0
284	646	719	0	0
284	641	711	2	0
284	641	707	0	0
284	646	713	0	0
284	647	713	0	0
284	641	709	0	0
43	643	748	4	0
284	655	718	2	0
284	655	717	2	0
284	654	719	0	0
284	655	715	2	0
284	655	716	2	0
26	653	717	0	0
284	653	719	0	0
284	653	713	0	0
284	653	715	0	0
284	655	714	2	0
284	655	719	0	0
284	655	713	0	0
284	654	715	0	0
402	654	745	7	0
34	650	746	0	0
395	649	750	2	0
284	650	719	0	0
34	652	736	0	0
402	652	742	7	0
34	648	740	0	0
395	648	743	4	0
284	648	719	0	0
284	652	719	0	0
284	649	715	0	0
284	649	719	0	0
284	651	719	0	0
284	649	713	0	0
284	651	713	0	0
284	652	713	0	0
284	651	715	0	0
284	648	715	0	0
284	650	715	0	0
284	650	713	0	0
1	649	735	7	0
284	648	713	0	0
1	648	727	1	0
284	654	713	0	0
284	652	715	0	0
284	655	707	0	0
284	655	706	2	0
284	655	711	0	0
284	654	711	0	0
284	655	704	2	0
284	652	707	0	0
284	652	711	0	0
6	652	709	2	0
0	652	704	0	0
284	655	705	2	0
284	653	711	0	0
284	654	707	0	0
284	653	707	0	0
284	651	711	0	0
284	648	709	0	0
284	648	707	0	0
284	648	711	2	0
284	649	707	0	0
284	648	710	2	0
284	650	710	2	0
284	650	711	0	0
284	651	707	0	0
284	649	709	0	0
284	650	707	0	0
284	650	709	0	0
145	639	1685	2	0
145	635	1685	0	0
44	643	1692	4	0
145	636	1680	4	0
44	643	1701	0	0
1021	636	1685	0	0
145	634	1680	4	0
145	638	1680	4	0
145	639	1681	2	0
6	638	1685	0	0
25	639	1680	4	0
145	634	1681	6	0
145	634	1685	6	0
23	635	1680	4	0
145	634	1683	6	0
145	637	1685	0	0
145	639	1683	2	0
959	634	2629	0	0
959	634	2624	2	0
959	639	2624	4	0
1014	636	2624	4	0
1017	636	2629	0	0
959	639	2629	6	0
3	616	753	0	0
43	643	757	0	0
395	652	757	0	0
1	620	758	0	0
1	623	752	0	0
1	653	752	7	0
63	647	753	0	0
285	631	755	0	0
63	642	753	4	0
285	630	756	0	0
37	632	755	0	0
285	634	756	0	0
37	630	755	0	0
285	633	755	0	0
37	633	756	0	0
37	631	756	0	0
285	632	756	0	0
37	634	755	0	0
1	643	752	0	1
1	643	756	0	1
1	616	759	0	1
3	614	753	0	0
1	613	758	0	0
205	608	757	3	0
1	615	754	0	1
205	603	753	5	0
41	602	756	4	0
20	602	755	2	0
38	602	752	2	0
38	607	754	1	0
205	607	755	5	0
51	601	754	0	0
47	606	757	4	0
20	603	755	4	0
205	605	753	5	0
6	605	759	2	0
205	605	752	3	0
51	604	754	0	0
38	597	759	4	0
205	599	753	0	0
283	595	759	0	0
205	599	754	4	0
146	599	757	1	1
284	594	759	0	0
284	592	759	0	0
284	593	759	0	0
205	597	754	2	0
283	591	759	0	0
20	585	754	0	0
139	584	755	4	0
283	585	759	0	0
1	572	752	0	0
50	583	764	1	0
38	598	763	7	0
50	584	763	1	0
42	591	761	6	0
283	584	760	0	0
205	597	761	6	0
94	581	761	2	1
1	591	765	0	1
405	544	755	0	0
396	538	758	1	0
396	537	746	1	0
396	536	749	1	0
405	538	746	0	0
399	538	753	1	0
404	536	751	7	0
405	537	753	7	0
398	539	752	1	0
32	537	764	0	0
32	539	762	0	0
32	537	761	0	0
396	538	760	1	0
396	539	745	1	0
405	541	755	3	0
405	542	744	0	0
396	540	748	1	0
63	536	754	0	0
396	541	754	1	0
396	537	755	1	0
398	542	756	1	0
32	536	757	0	0
32	539	753	0	0
399	536	752	1	0
398	536	759	1	0
32	535	751	0	0
7	534	753	2	0
51	533	756	4	0
51	532	754	2	0
397	533	751	1	0
51	534	756	4	0
51	532	755	2	0
32	534	759	0	0
396	535	758	1	0
404	534	749	1	0
398	533	758	1	0
404	535	760	1	0
399	529	757	1	0
399	530	752	1	0
405	529	754	1	0
397	532	749	1	0
400	530	753	1	0
405	530	756	7	0
397	531	754	1	0
32	533	750	0	0
32	531	757	0	0
404	530	755	1	0
402	533	748	1	0
405	531	750	0	0
404	532	751	7	0
400	531	751	1	0
404	532	759	3	0
3	533	753	0	0
1	589	771	3	0
34	591	771	0	0
283	590	769	0	0
283	590	773	0	0
1	596	770	3	0
0	598	771	6	0
34	596	772	3	0
308	596	774	6	0
283	592	772	0	0
306	599	773	6	0
0	592	774	3	0
283	599	774	0	0
306	592	770	6	0
283	594	773	0	0
0	598	776	3	0
283	595	776	0	0
51	601	761	4	0
47	602	761	2	0
205	602	762	5	0
51	604	761	4	0
283	607	777	0	0
306	607	779	6	0
1	605	776	3	0
0	606	775	3	0
283	602	774	0	0
34	602	771	3	0
205	600	763	5	0
205	606	761	4	0
205	604	762	5	0
205	606	763	4	0
34	600	777	3	0
7	615	760	4	0
7	615	762	0	0
7	614	762	0	0
7	614	760	4	0
273	614	761	0	0
38	608	760	7	0
283	613	775	0	0
308	611	774	6	0
283	609	775	0	0
34	610	777	3	0
283	612	779	0	0
1	609	776	3	0
306	612	776	6	0
283	611	778	0	0
1	609	781	3	0
283	610	780	0	0
1	621	761	0	0
274	618	761	2	0
278	631	764	0	0
16	630	765	0	0
278	629	763	0	0
3	630	764	0	0
16	630	766	0	0
1	631	761	0	1
5	627	761	0	0
16	628	764	2	0
16	629	764	2	0
278	631	766	0	0
68	631	773	3	0
281	634	762	0	0
302	638	762	4	0
279	634	766	0	0
3	634	765	0	0
278	634	761	0	0
279	633	766	0	0
278	634	763	0	0
566	632	760	2	0
962	626	779	6	0
972	628	778	7	0
394	630	781	3	0
962	631	783	2	0
971	638	779	0	0
398	632	783	3	0
34	634	778	3	0
283	644	767	6	0
283	643	767	6	0
399	645	783	0	0
205	649	780	0	0
962	654	775	3	0
1	649	761	7	0
97	654	777	0	0
377	653	778	3	0
970	649	770	1	0
962	654	780	1	0
377	652	774	1	0
395	652	775	0	0
377	663	732	0	0
38	660	737	2	0
938	661	741	2	0
394	662	732	1	0
97	662	736	7	0
395	660	742	0	0
377	661	739	3	0
377	659	736	2	0
38	662	730	2	0
283	662	743	7	0
97	661	736	7	0
97	661	735	7	0
395	661	729	0	0
205	657	749	0	0
1	662	742	7	0
938	658	736	2	0
1	661	744	7	0
1	658	744	7	0
938	658	740	2	0
938	657	741	2	0
938	659	735	2	0
377	660	733	1	0
286	663	761	6	0
283	663	741	2	0
97	662	735	7	0
283	663	743	5	0
34	661	746	0	0
377	662	763	3	0
938	659	742	2	0
34	659	738	0	0
938	659	741	2	0
377	658	756	1	0
38	659	740	2	0
286	660	757	6	0
377	662	755	0	0
938	659	739	2	0
34	663	751	0	0
395	659	755	1	0
938	657	736	2	0
377	659	762	3	0
377	657	758	2	0
286	658	759	6	0
377	657	761	3	0
938	656	738	2	0
394	657	775	0	0
1	629	1709	1	1
6	627	1705	0	0
15	631	1705	2	0
15	631	1709	2	0
1	629	1706	1	1
34	656	783	2	0
962	662	781	7	0
664	664	757	7	0
662	665	754	0	0
283	664	745	5	0
205	665	752	0	0
0	665	744	7	0
398	669	765	7	0
377	665	762	5	0
404	665	771	0	0
395	665	780	0	0
399	665	782	2	0
989	666	772	0	0
398	664	781	6	0
405	666	775	3	0
395	669	771	0	0
205	667	767	1	0
377	666	760	6	0
962	668	771	7	0
398	669	742	7	0
283	665	743	5	0
377	665	739	6	0
283	665	741	2	0
0	664	742	2	0
34	664	740	0	0
34	665	734	0	0
398	668	735	3	0
394	662	723	7	0
34	658	723	0	0
1	669	727	7	0
34	665	725	0	0
398	658	726	1	0
284	664	717	0	0
284	664	715	0	0
284	665	712	2	0
284	667	714	2	0
284	667	717	0	0
284	667	712	2	0
284	671	714	2	0
284	669	715	0	0
284	669	713	0	0
284	657	717	0	0
284	667	715	2	0
284	657	713	0	0
284	671	717	2	0
284	671	719	0	0
284	667	713	2	0
284	667	716	2	0
284	665	719	0	0
284	666	717	0	0
284	669	719	0	0
284	668	717	0	0
284	670	719	0	0
284	666	719	0	0
284	671	713	2	0
284	671	712	2	0
284	668	713	0	0
284	667	719	0	0
284	668	719	0	0
284	669	717	0	0
284	658	719	0	0
284	670	715	0	0
284	661	717	0	0
284	657	716	2	0
284	659	717	0	0
284	659	719	0	0
284	659	714	2	0
284	659	715	0	0
284	661	713	2	0
284	659	712	2	0
284	657	719	0	0
284	671	716	2	0
284	660	719	0	0
284	656	719	0	0
284	661	714	2	0
284	660	717	0	0
284	661	715	0	0
284	656	713	0	0
284	657	715	0	0
284	658	717	0	0
284	661	719	0	0
284	671	715	0	0
284	663	717	0	0
284	659	713	2	0
284	658	715	0	0
284	663	719	0	0
284	662	717	0	0
284	663	713	2	0
284	662	715	0	0
284	663	712	2	0
284	663	715	0	0
284	662	719	0	0
284	665	714	2	0
284	664	719	0	0
284	665	713	2	0
284	665	715	0	0
284	665	717	0	0
284	659	711	2	0
284	659	710	2	0
284	658	711	0	0
284	657	704	2	0
284	657	707	0	0
284	663	705	0	0
284	661	711	0	0
284	663	707	0	0
284	659	709	2	0
284	659	707	0	0
284	659	705	0	0
284	657	708	2	0
284	661	707	0	0
284	660	711	0	0
284	661	709	0	0
284	661	705	0	0
284	657	705	0	0
284	663	711	0	0
284	663	709	0	0
284	662	707	0	0
23	661	704	0	0
284	660	705	0	0
284	656	707	0	0
284	657	709	2	0
284	656	711	0	0
284	662	709	0	0
284	657	711	0	0
284	662	711	0	0
284	658	705	0	0
284	658	707	0	0
284	660	707	0	0
284	662	705	0	0
284	671	707	2	0
284	671	710	2	0
284	671	708	2	0
284	669	708	2	0
284	669	711	0	0
284	667	711	2	0
284	669	707	2	0
284	667	709	0	0
284	671	709	2	0
284	667	710	2	0
284	669	709	2	0
284	667	705	0	0
284	665	707	0	0
284	665	709	2	0
284	665	708	2	0
284	665	710	2	0
284	671	706	2	0
284	671	711	2	0
284	671	704	2	0
284	671	705	2	0
284	670	711	0	0
284	669	706	2	0
284	668	705	0	0
284	665	711	2	0
284	669	705	0	0
284	665	705	0	0
284	664	705	0	0
284	666	705	0	0
284	668	709	0	0
284	666	707	0	0
284	664	709	0	0
23	664	704	0	0
284	667	707	0	0
284	664	707	0	0
973	650	788	0	0
38	660	787	2	0
97	630	785	0	0
205	648	790	2	0
205	662	785	2	0
394	640	789	3	0
970	640	786	5	0
217	636	791	3	0
404	634	785	6	0
88	644	786	7	0
377	635	790	0	0
949	630	788	5	0
962	634	784	2	0
34	653	789	2	0
377	628	785	3	0
399	632	784	3	0
962	625	789	6	0
34	633	791	2	0
973	645	790	2	0
965	624	786	6	0
34	637	789	0	0
430	644	788	1	0
395	627	785	3	0
973	652	785	2	0
965	619	788	6	0
965	622	787	6	0
965	620	787	6	0
0	621	791	0	0
965	622	790	6	0
34	611	787	5	0
164	607	795	2	0
70	606	797	1	0
283	621	795	0	0
377	601	798	1	0
597	616	799	2	0
34	606	796	0	0
597	616	798	2	0
283	620	793	0	0
988	630	793	2	0
597	615	799	4	0
377	609	798	6	0
164	609	795	3	0
965	623	792	6	0
597	616	797	2	0
951	638	797	4	0
951	639	797	4	0
951	639	796	0	0
398	629	793	5	0
951	638	796	0	0
394	638	794	4	0
34	636	799	0	0
962	637	792	2	0
97	635	792	0	0
377	635	799	4	0
430	634	792	0	0
962	634	798	4	0
962	632	799	7	0
34	615	802	0	0
201	614	800	1	1
978	606	804	0	0
377	609	802	6	0
34	615	805	5	0
164	602	806	5	0
377	601	801	2	0
394	605	804	1	0
34	603	802	0	0
377	607	805	4	0
397	605	805	2	0
997	604	803	2	0
597	616	802	2	0
34	617	802	4	0
34	616	803	1	0
34	618	804	6	0
597	616	801	2	0
950	625	807	2	0
395	624	803	0	0
402	632	800	0	0
394	632	806	0	0
377	639	806	6	0
397	632	807	0	0
395	635	804	2	0
398	634	806	0	0
377	635	805	0	0
34	638	808	0	0
399	612	812	2	0
951	638	811	0	0
951	637	812	4	0
97	636	808	0	0
34	634	809	0	0
377	635	809	4	0
951	635	811	0	0
962	634	808	1	0
962	624	810	6	0
398	614	811	2	0
951	639	811	0	0
951	638	812	4	0
951	639	812	4	0
951	636	812	4	0
951	637	811	0	0
951	635	812	4	0
951	636	811	0	0
38	619	822	3	0
97	617	823	6	0
402	619	820	6	0
402	621	819	7	0
402	618	822	4	0
402	618	819	0	0
402	620	823	1	0
38	615	821	0	0
402	617	821	0	0
97	616	823	6	0
0	611	817	0	0
407	618	820	0	0
402	614	823	0	0
402	616	819	0	0
962	629	817	0	0
34	634	819	0	0
34	632	818	0	0
283	606	825	6	0
1	603	830	0	0
283	600	824	0	0
0	631	828	1	0
0	617	828	0	0
38	634	828	1	0
36	636	830	1	0
36	632	831	1	0
0	628	829	1	0
0	595	822	0	0
34	596	830	7	0
1	590	824	0	0
34	590	831	2	0
0	584	818	0	0
283	588	828	0	0
34	590	833	1	0
34	590	834	4	0
68	588	834	2	0
34	591	833	1	0
34	588	832	2	0
34	587	835	1	0
34	588	836	5	0
34	589	833	6	0
34	589	835	2	0
34	590	832	3	0
33	593	835	6	0
34	585	835	2	0
34	612	834	7	0
34	585	836	0	0
0	614	839	7	0
34	585	837	7	0
395	612	836	7	0
34	592	832	1	0
34	593	832	0	0
34	586	832	2	0
34	586	834	2	0
1	607	835	0	0
34	587	833	2	0
0	611	832	0	0
283	621	833	0	0
1	623	838	0	0
283	620	838	0	0
395	619	834	7	0
34	618	839	7	0
34	616	834	7	0
283	583	820	2	0
0	580	833	1	0
1	581	826	0	0
283	581	818	6	0
1	579	820	1	0
0	579	827	4	0
0	581	825	6	0
0	579	822	0	0
285	580	831	6	0
283	580	830	6	0
1	580	823	0	0
285	582	835	6	0
0	582	832	3	0
283	581	832	2	0
283	579	825	2	0
285	582	824	6	0
283	577	821	6	0
283	582	828	6	0
1	578	826	1	0
285	581	830	6	0
0	580	837	7	0
285	578	824	6	0
283	577	829	2	0
283	579	834	7	0
285	577	834	6	0
0	577	838	6	0
1	579	832	0	0
283	582	822	6	0
34	583	843	6	0
1	583	845	6	0
37	581	845	6	0
37	584	847	6	0
37	584	841	6	0
285	584	844	6	0
37	587	843	6	0
34	585	845	6	0
0	577	840	6	0
0	579	841	5	0
402	599	840	7	0
402	597	846	0	0
402	598	842	6	0
402	597	847	6	0
402	597	845	0	0
402	598	841	0	0
402	598	843	4	0
401	598	840	5	0
401	597	844	6	0
402	597	843	5	0
402	596	845	1	0
401	595	847	4	0
402	596	847	4	0
283	591	852	7	0
37	586	848	6	0
402	595	848	0	0
283	594	855	7	0
402	594	849	3	0
402	594	850	1	0
402	593	848	3	0
401	592	851	3	0
34	607	856	6	0
283	607	857	7	0
34	605	858	6	0
0	608	856	7	0
99	615	845	1	0
283	608	858	7	0
113	615	842	6	0
34	609	859	6	0
283	614	840	6	0
99	612	843	1	0
99	609	840	0	0
99	621	842	1	0
99	621	843	6	0
283	621	854	7	0
1	617	844	7	0
99	616	841	1	0
283	620	843	6	0
283	620	840	6	0
283	619	846	7	0
113	617	842	5	0
99	621	845	6	0
283	617	840	3	0
99	618	844	7	0
113	618	842	3	0
99	619	841	2	0
1	631	837	1	0
34	630	833	1	0
99	631	843	6	0
99	628	848	6	0
99	626	848	6	0
103	626	853	2	0
99	625	855	6	0
103	629	855	6	0
99	628	852	6	0
103	624	857	6	0
36	637	834	1	0
205	633	836	1	0
32	638	849	5	0
32	637	850	6	0
103	637	853	6	0
962	644	822	0	0
397	645	817	0	0
394	644	818	0	0
286	640	830	3	0
0	642	839	1	0
283	647	840	2	0
205	640	832	1	0
286	647	842	3	0
975	647	812	2	0
951	646	812	4	0
951	646	811	0	0
951	642	812	4	0
951	643	812	4	0
951	643	811	0	0
951	642	811	0	0
951	641	812	4	0
951	641	811	0	0
951	640	812	4	0
951	640	811	0	0
951	645	812	4	0
951	645	811	0	0
951	644	811	0	0
951	644	812	4	0
995	647	802	0	0
951	646	804	2	0
951	646	803	2	0
34	643	806	0	0
377	644	806	2	0
951	646	802	2	0
951	647	803	6	0
951	647	804	6	0
1	651	839	1	0
205	652	828	2	0
377	651	813	0	0
962	651	821	0	0
34	649	821	0	0
398	650	829	6	0
398	653	828	0	0
285	648	838	1	0
283	651	836	2	0
398	649	827	6	0
962	654	803	1	0
309	649	839	1	0
951	655	812	4	0
283	653	837	2	0
395	651	804	6	0
377	649	804	0	0
398	652	831	0	0
164	655	801	6	0
307	655	838	1	0
951	655	811	0	0
962	652	829	1	0
395	651	828	6	0
36	654	841	1	0
38	652	846	1	0
205	653	842	1	0
1	649	843	1	0
283	655	857	1	0
306	652	857	1	0
283	649	859	1	0
307	658	844	1	0
34	659	838	1	0
1	662	847	1	0
36	657	840	1	0
1	663	850	1	0
283	659	854	1	0
0	662	849	1	0
0	661	851	1	0
283	663	853	1	0
216	660	855	2	0
0	657	853	1	0
283	659	850	1	0
1	656	855	1	0
34	657	825	0	0
38	662	831	0	0
962	660	825	0	0
951	663	812	4	0
951	662	811	0	0
951	663	811	0	0
951	662	812	4	0
951	657	812	4	0
951	659	812	4	0
951	658	812	4	0
951	659	811	0	0
951	660	812	4	0
951	660	811	0	0
951	661	812	4	0
951	661	811	0	0
951	658	811	0	0
951	656	811	0	0
951	657	811	0	0
951	656	812	4	0
283	666	837	1	0
283	668	845	1	0
1	664	843	1	0
283	667	844	1	0
286	664	845	3	0
0	666	848	1	0
283	668	848	1	0
306	666	846	1	0
0	667	841	1	0
286	669	841	3	0
0	664	847	1	0
1	665	844	1	0
998	668	825	5	0
0	669	824	3	0
0	664	829	3	0
286	666	827	3	0
286	669	826	3	0
306	664	827	5	0
283	665	826	5	0
1	668	828	3	0
962	666	822	1	0
1	670	822	3	0
962	647	796	7	0
951	643	796	0	0
951	643	797	4	0
398	646	793	0	0
394	643	799	0	0
951	642	797	4	0
951	642	796	0	0
951	640	797	4	0
951	640	796	0	0
951	641	797	4	0
951	641	796	0	0
377	645	794	0	0
397	647	792	0	0
1001	646	795	4	0
394	655	798	0	0
996	648	799	0	0
395	649	793	4	0
34	649	799	0	0
377	649	796	6	0
973	652	792	0	0
34	655	799	0	0
962	655	793	4	0
377	648	792	0	0
398	650	794	3	0
34	664	786	2	0
951	666	789	2	0
955	664	788	2	0
88	665	784	1	0
951	666	787	2	0
951	666	788	2	0
951	666	790	2	0
205	661	798	0	0
164	657	801	4	0
1003	659	804	2	0
962	663	806	3	0
377	656	799	0	0
395	659	796	4	0
97	657	803	0	0
164	657	800	7	0
164	657	802	4	0
962	658	801	3	0
377	658	806	4	0
164	656	800	0	0
164	659	802	2	0
398	659	800	0	0
164	656	801	0	0
419	646	691	4	0
0	643	689	0	0
0	641	694	0	0
306	644	694	0	0
5	642	691	2	0
37	644	700	0	0
419	647	703	0	0
284	655	690	2	0
69	641	696	0	0
306	642	699	0	0
5	643	703	2	0
420	640	703	0	0
284	655	689	2	0
284	655	701	2	0
419	650	699	2	0
284	655	702	2	0
284	655	703	2	0
0	649	689	0	0
284	655	688	2	0
0	653	688	0	0
425	650	691	6	0
425	650	703	4	0
0	649	694	0	0
424	653	696	2	0
0	654	702	0	0
284	655	700	2	0
284	660	702	0	0
284	657	703	2	0
284	656	700	0	0
284	660	700	0	0
284	661	700	0	0
284	657	700	0	0
284	657	702	2	0
284	661	702	0	0
284	662	700	0	0
284	663	700	0	0
284	659	702	0	0
284	659	700	0	0
284	663	702	0	0
284	658	702	0	0
284	658	700	0	0
284	662	702	0	0
427	636	1634	2	0
426	650	1635	6	0
6	636	1640	0	0
6	643	1647	2	0
421	640	1647	0	0
10	652	1639	4	0
421	636	1643	2	0
421	639	1635	0	0
423	653	1640	2	0
426	650	1647	4	0
418	647	1647	0	0
418	646	1635	0	0
418	650	1643	2	0
6	642	1635	2	0
51	647	3535	6	0
51	658	3538	2	0
8	661	3534	3	0
51	661	3532	2	0
51	654	3536	4	0
5	652	3541	2	0
51	661	3535	2	0
51	652	3536	4	0
480	663	3530	6	0
51	648	3539	0	0
51	654	3539	0	0
164	653	3554	5	0
164	652	3554	4	0
290	654	3531	3	0
51	646	3536	2	0
290	653	3533	0	0
51	654	3542	4	0
290	651	3532	0	0
290	653	3536	0	0
51	649	3534	2	0
290	653	3534	0	0
51	646	3533	0	0
51	647	3537	6	0
1007	642	3555	0	0
51	646	3538	2	0
51	653	3531	0	0
164	653	3555	2	0
290	655	3533	0	0
481	655	3535	0	0
205	649	3559	4	0
7	650	3535	3	0
1007	639	3559	0	0
7	665	3532	2	0
51	646	3540	4	0
1007	636	3558	0	0
55	651	3534	3	0
1007	637	3555	0	0
290	655	3534	0	0
51	662	3536	4	0
17	667	3524	4	0
51	665	3528	4	0
22	667	3527	6	0
22	668	3525	6	0
55	668	3524	4	0
1007	658	3559	0	0
55	668	3527	0	0
25	669	3526	6	0
1007	661	3556	0	0
51	669	3527	6	0
51	659	3540	4	0
51	657	3539	0	0
1007	667	3559	0	0
51	658	3536	2	0
51	659	3537	6	0
290	663	3536	3	0
51	669	3525	6	0
51	659	3534	6	0
51	658	3533	0	0
51	650	3542	4	0
116	648	3558	1	0
63	663	3531	6	0
51	663	3527	2	0
290	651	3536	0	0
51	651	3539	0	0
51	650	3536	4	0
290	649	3535	0	0
55	668	3526	6	0
51	666	3524	0	0
55	665	3525	3	0
51	666	3532	6	0
51	668	3528	4	0
3	665	3534	6	0
51	666	3535	6	0
1	661	3533	1	1
3	664	3534	6	0
51	665	3536	4	0
1	649	3533	1	1
188	649	3554	2	1
937	677	3526	4	0
117	678	3526	5	0
51	673	3520	2	0
51	673	3540	2	0
51	673	3547	2	0
116	674	3539	1	0
51	675	3550	4	0
51	673	3531	2	0
51	673	3523	2	0
51	675	3538	0	0
51	678	3543	4	0
51	673	3528	2	0
51	677	3512	4	0
51	674	3512	4	0
116	678	3515	1	0
37	647	3302	0	0
51	641	3308	4	0
63	642	3304	4	0
3	653	3292	0	0
22	664	3291	0	0
71	664	3294	2	0
22	664	3292	0	0
51	667	3290	0	0
0	644	3306	0	0
51	669	3294	4	0
55	669	3291	6	0
51	669	3290	0	0
471	659	3303	6	0
51	659	3292	2	0
51	664	3297	0	0
51	662	3298	6	0
34	646	3302	2	0
51	660	3293	6	0
1	657	3297	1	1
63	659	3295	6	0
51	661	3299	4	0
51	655	3294	4	0
51	653	3290	0	0
20	668	3292	6	0
51	659	3290	2	0
51	658	3288	4	0
51	670	3292	6	0
51	660	3291	6	0
469	659	3304	4	0
51	662	3296	6	0
51	657	3298	2	0
64	659	3289	6	0
20	658	3295	6	0
20	661	3295	0	0
51	655	3297	0	0
51	661	3288	4	0
22	663	3291	0	0
51	653	3297	2	0
51	649	3290	0	0
51	663	3292	2	0
37	648	3303	0	0
47	649	3291	0	0
51	657	3296	2	0
51	658	3299	4	0
37	646	3306	0	0
22	654	3291	0	0
51	663	3290	0	0
7	654	3292	2	0
51	665	3290	0	0
23	642	3302	2	0
63	653	3295	6	0
51	655	3290	0	0
51	649	3294	4	0
467	646	3305	0	0
22	655	3291	0	0
51	651	3294	4	0
51	655	3292	6	0
51	663	3294	4	0
0	645	3303	0	0
55	652	3293	1	0
34	644	3303	2	0
37	643	3303	0	0
7	652	3292	6	0
492	650	3290	0	0
23	642	3306	2	0
2	663	3297	1	1
37	646	3304	0	0
51	651	3290	0	0
51	641	3301	0	0
63	667	3295	6	0
22	670	3290	0	1
51	663	3285	6	0
51	663	3283	6	0
51	663	3287	6	0
7	662	3284	2	0
51	662	3282	0	0
20	661	3287	6	0
51	667	3286	4	0
51	656	3287	2	0
51	668	3283	6	0
20	658	3287	6	0
51	657	3282	0	0
51	669	3286	4	0
51	667	3284	2	0
51	656	3285	2	0
274	659	3282	4	0
9	658	3284	4	0
51	656	3283	2	0
7	657	3284	6	0
51	667	3281	2	0
135	668	3281	0	1
158	175	1431	0	0
11	175	1427	2	0
5	178	1430	0	0
45	176	1430	2	0
3	178	1424	0	0
6	180	1430	0	0
28	178	1426	0	0
45	176	1431	2	0
3	176	1424	0	0
3	181	1424	0	0
3	177	1424	0	0
3	180	1424	0	0
11	175	1426	2	0
3	181	2374	0	0
3	177	2370	0	0
6	177	2374	0	0
3	175	2370	0	0
3	176	2370	0	0
173	179	2371	0	0
145	127	3525	2	0
1	135	3534	7	0
51	127	3513	6	0
72	140	3528	7	0
72	141	3527	7	0
0	135	3524	7	0
51	130	3518	6	0
72	140	3527	7	0
145	130	3516	2	0
72	140	3526	7	0
72	143	3528	7	0
72	142	3528	7	0
72	142	3529	7	0
72	140	3530	7	0
72	142	3530	7	0
72	139	3528	7	0
72	139	3530	7	0
51	122	3518	2	0
145	125	3525	6	0
51	125	3523	2	0
302	136	3533	4	0
51	127	3523	6	0
72	139	3529	7	0
1	137	3535	7	0
72	139	3527	7	0
145	122	3516	6	0
5	128	3518	0	0
51	117	3534	2	0
145	130	3520	2	0
51	125	3513	2	0
25	126	3533	0	0
51	120	3531	2	0
72	143	3527	7	0
72	142	3527	7	0
72	142	3526	7	0
51	123	3528	2	0
72	141	3526	7	0
72	143	3526	7	0
25	128	3535	0	0
0	138	3523	7	0
72	143	3530	7	0
72	141	3530	7	0
0	142	3521	7	0
72	141	3528	7	0
72	139	3526	7	0
72	141	3529	7	0
72	143	3529	7	0
145	122	3520	6	0
72	140	3529	7	0
25	129	3538	0	0
51	118	3541	6	0
51	118	3539	6	0
25	127	3541	0	0
51	113	3536	0	0
51	120	3537	4	0
25	134	3541	0	0
1	142	3537	7	0
67	117	3539	0	1
25	127	3538	0	0
67	116	3537	1	1
51	115	3536	0	0
30	104	3547	3	0
30	115	3548	5	0
25	127	3547	0	0
25	134	3544	0	0
25	134	3547	0	0
25	127	3544	0	0
30	117	3544	5	0
10	130	3547	2	0
30	103	3537	1	0
249	98	3537	0	0
72	148	3526	7	0
72	148	3527	7	0
72	147	3526	7	0
72	147	3528	7	0
72	148	3530	7	0
72	145	3526	7	0
72	144	3526	7	0
72	146	3529	7	0
72	146	3530	7	0
72	144	3527	7	0
72	147	3527	7	0
72	146	3527	7	0
72	146	3528	7	0
72	148	3529	7	0
72	147	3529	7	0
72	146	3526	7	0
72	145	3527	7	0
72	144	3528	7	0
72	147	3530	7	0
72	145	3530	7	0
72	148	3528	7	0
51	145	3539	4	0
51	148	3541	4	0
1	147	3537	0	0
72	144	3530	7	0
72	144	3529	7	0
1	150	3538	0	0
72	145	3529	7	0
51	151	3541	4	0
72	145	3528	7	0
51	145	3519	0	0
51	150	3517	0	0
1	153	3536	0	0
51	156	3514	0	0
26	153	3531	4	0
3	158	3523	0	0
246	159	3533	0	0
26	153	3524	4	0
3	159	3523	0	0
51	167	3526	0	0
51	162	3524	6	0
118	161	3516	4	0
53	162	3533	0	0
51	164	3533	6	0
51	160	3535	4	0
51	160	3519	4	0
51	160	3514	0	0
51	164	3525	4	0
51	165	3538	2	0
51	163	3526	0	0
51	166	3521	0	0
1	160	3526	0	1
1	166	3526	0	1
1	161	3531	0	1
51	174	3521	0	0
1	169	3533	1	1
51	170	3525	4	0
3	170	3531	0	0
51	171	3521	0	0
11	172	3534	0	0
51	174	3541	0	0
51	172	3545	2	0
51	174	3535	4	0
11	172	3536	0	0
58	171	3527	4	0
51	169	3523	6	0
3	170	3532	0	0
3	171	3531	0	0
51	172	3548	2	0
3	172	3532	0	0
51	169	3543	4	0
3	171	3532	0	0
51	175	3550	4	0
3	172	3531	0	0
51	169	3538	0	0
51	179	3548	6	0
51	178	3531	6	0
51	176	3541	0	0
51	180	3545	6	0
51	176	3523	6	0
51	178	3534	6	0
51	176	3535	4	0
51	177	3550	4	0
51	176	3527	6	0
22	602	1686	0	0
51	599	1703	2	0
6	605	1686	4	0
22	601	1686	0	0
51	599	1700	2	0
5	606	1703	6	0
45	604	1702	2	0
45	604	1701	2	0
45	603	1703	0	0
51	604	1705	4	0
51	601	1698	0	0
22	605	1704	5	0
3	600	1704	6	0
45	604	1700	2	0
45	601	1702	6	0
45	601	1701	6	0
51	604	1698	0	0
1	602	1704	1	1
45	601	1700	6	0
42	602	1700	4	0
45	602	1703	0	0
15	601	1688	0	0
1	600	1703	0	1
51	606	1700	6	0
15	605	1688	0	0
51	601	1705	4	0
148	599	2645	1	1
6	606	2647	0	0
63	602	2642	6	0
25	600	2643	5	0
25	605	2648	5	0
147	602	2642	0	1
148	599	2646	1	1
147	603	2642	0	1
20	601	2642	6	0
63	599	2645	0	0
63	602	2649	2	0
25	605	2643	7	0
20	604	2642	0	0
20	599	2644	2	0
20	604	2649	2	0
20	599	2647	0	0
149	602	2650	0	1
20	601	2649	4	0
25	600	2648	7	0
149	603	2650	0	1
6	444	1707	2	0
15	446	1705	4	0
6	446	1713	0	0
15	442	1714	0	0
15	441	1706	4	0
21	454	1699	0	0
21	458	1695	0	0
21	456	1697	0	0
21	457	1703	0	0
21	459	1695	0	0
6	441	1714	0	0
21	454	1700	0	0
21	462	1699	0	0
6	451	1714	0	0
6	449	1706	2	0
6	454	1708	2	0
21	458	1703	0	0
15	452	1714	0	0
21	462	1698	0	0
21	460	1701	0	0
15	451	1707	4	0
6	458	1699	5	0
15	447	1713	0	0
6	429	1715	0	0
734	678	3418	5	0
740	678	3415	0	0
739	677	3414	0	0
734	676	3421	0	0
832	675	3422	0	0
730	681	3417	0	0
750	682	3409	2	0
940	680	3418	0	0
754	687	3423	0	0
730	680	3420	4	0
734	681	3421	3	0
940	684	3420	0	0
738	674	3416	0	0
778	675	3429	0	0
729	686	3418	0	0
734	675	3417	6	0
730	674	3418	7	0
833	684	3438	0	0
797	679	3434	2	0
799	678	3434	2	0
799	679	3435	2	0
729	684	3422	0	0
772	680	3434	6	0
754	689	3419	0	0
730	674	3422	0	0
754	689	3421	2	0
734	682	3415	0	0
754	689	3423	7	0
728	672	3420	6	0
741	685	3410	0	0
799	677	3434	0	0
754	691	3421	4	0
743	692	3412	0	0
798	681	3435	5	0
754	691	3420	3	0
754	691	3423	6	0
734	683	3422	2	0
730	685	3416	0	0
753	684	3430	2	0
730	682	3418	0	0
734	684	3417	0	0
798	682	3435	0	0
734	680	3424	0	0
753	685	3430	6	0
754	690	3424	1	0
730	684	3425	0	0
747	692	3429	2	0
754	690	3425	5	0
734	693	3435	0	0
729	690	3439	6	0
734	693	3434	0	0
754	689	3425	4	0
940	688	3438	0	0
734	693	3438	0	0
729	688	3436	0	0
798	681	3434	0	0
754	692	3417	5	0
742	688	3410	0	0
746	689	3429	4	0
754	692	3419	0	0
745	677	3424	0	0
754	690	3422	5	0
754	691	3418	0	0
754	691	3416	4	0
754	690	3418	6	0
734	695	3436	0	0
754	688	3425	2	0
754	692	3424	7	0
754	688	3416	2	0
754	688	3420	1	0
754	688	3418	7	0
754	688	3424	1	0
754	689	3416	0	0
744	696	3414	0	0
734	703	3432	3	0
749	699	3423	2	0
55	700	3419	1	0
55	702	3422	2	0
730	703	3426	0	0
752	700	3426	4	0
771	697	3436	1	0
748	697	3426	3	0
771	697	3435	3	0
730	697	3438	1	0
729	697	3437	4	0
754	697	3439	2	0
737	701	3416	0	0
734	700	3435	0	0
97	701	3420	0	0
771	698	3436	0	0
732	700	3415	2	0
752	701	3426	4	0
729	703	3425	0	0
731	699	3417	0	0
940	697	3420	0	0
55	700	3420	2	0
729	696	3438	0	0
751	696	3417	0	0
730	696	3437	6	0
772	696	3435	6	0
730	705	3414	0	0
729	704	3414	0	0
733	711	3418	0	0
770	709	3422	3	0
770	706	3429	0	0
730	706	3423	0	0
730	704	3424	0	0
729	705	3423	0	0
729	709	3429	0	0
729	705	3425	0	0
736	707	3420	4	0
726	704	3417	2	0
736	705	3420	0	0
729	711	3422	0	0
734	711	3430	0	0
770	716	3410	2	0
940	715	3412	0	0
734	713	3421	0	0
795	714	3418	0	0
97	714	3426	0	0
734	714	3427	0	0
729	714	3430	0	0
734	712	3424	0	0
770	712	3411	4	0
770	718	3412	0	0
770	714	3410	0	0
729	714	3423	0	0
795	715	3418	0	0
770	711	3436	0	0
734	704	3438	6	0
940	710	3434	0	0
770	714	3433	0	0
734	705	3432	2	0
734	707	3437	5	0
770	707	3439	0	0
729	707	3434	0	0
734	708	3432	4	0
940	706	3433	0	0
770	709	3438	2	0
729	710	3432	0	0
770	713	3435	0	0
770	708	3439	3	0
734	727	3410	0	0
729	726	3414	2	0
734	725	3412	1	0
97	723	3415	2	0
729	726	3412	4	0
729	727	3426	0	0
940	727	3416	0	0
730	725	3410	1	0
940	724	3412	0	0
182	722	3413	2	0
734	727	3418	6	0
729	723	3411	3	0
729	725	3419	5	0
734	724	3426	2	0
734	725	3416	5	0
22	723	3416	3	0
182	722	3417	3	0
734	726	3422	4	0
940	732	3415	0	0
734	732	3426	2	0
734	733	3431	5	0
734	735	3430	2	0
729	733	3429	1	0
734	734	3422	7	0
734	730	3415	2	0
729	735	3428	3	0
734	728	3428	2	0
734	731	3411	5	0
734	734	3413	5	0
729	731	3417	7	0
829	728	3436	0	0
734	732	3419	0	0
734	728	3413	6	0
828	728	3438	0	0
734	735	3417	5	0
734	701	3440	0	0
729	726	3445	2	0
730	730	3442	2	0
734	728	3443	0	0
734	726	3443	0	0
734	703	3446	0	0
810	713	3447	2	0
754	699	3441	2	0
814	727	3447	0	0
825	728	3440	0	0
734	730	3444	0	0
754	696	3443	2	0
754	696	3441	2	0
734	731	3446	0	0
940	729	3446	0	0
730	724	3446	0	0
754	697	3441	2	0
754	696	3440	2	0
754	698	3442	2	0
754	698	3440	2	0
754	698	3443	2	0
754	697	3442	2	0
809	711	3447	2	0
97	706	3447	0	0
811	715	3447	2	0
812	706	3440	0	0
770	691	3441	0	0
754	695	3442	2	0
770	688	3440	0	0
782	688	3445	0	0
782	688	3447	0	0
734	692	3447	0	0
729	692	3440	1	0
801	689	3444	4	0
806	695	3449	2	0
801	689	3453	0	0
793	690	3452	0	0
806	695	3448	2	0
802	690	3449	6	0
782	688	3449	0	0
807	697	3449	2	0
808	699	3448	2	0
809	711	3448	2	0
807	697	3448	2	0
791	688	3453	0	0
808	699	3449	2	0
734	707	3451	0	0
782	688	3451	0	0
813	706	3453	0	0
834	707	3450	1	0
940	704	3449	0	0
811	715	3448	2	0
835	727	3450	0	0
734	723	3448	0	0
144	725	3451	4	0
810	713	3448	2	0
729	725	3449	4	0
782	684	3445	0	0
782	684	3451	0	0
782	686	3447	0	0
785	682	3447	0	0
782	684	3447	0	0
789	686	3451	0	0
782	684	3453	0	0
782	680	3453	0	0
782	680	3451	0	0
790	686	3453	0	0
788	686	3449	0	0
782	680	3449	0	0
786	682	3449	0	0
782	686	3445	0	0
782	682	3453	0	0
787	684	3449	0	0
782	682	3445	0	0
782	680	3445	0	0
782	682	3451	0	0
777	680	3447	0	0
792	679	3446	0	0
729	679	3451	0	0
734	678	3448	0	0
730	679	3463	4	0
730	675	3463	4	0
730	678	3459	4	0
729	679	3458	4	0
730	678	3461	4	0
730	684	3461	4	0
730	684	3459	7	0
730	681	3462	4	0
759	682	3463	6	0
760	681	3460	0	0
940	681	3461	0	0
755	673	3459	6	0
734	676	3458	1	0
730	685	3460	4	0
755	674	3460	4	0
729	685	3462	4	0
729	678	3458	1	0
767	672	3462	6	0
734	686	3458	7	0
729	683	3463	0	0
729	683	3458	4	0
730	683	3462	5	0
760	681	3459	6	0
756	689	3460	1	0
730	688	3462	4	0
757	688	3461	2	0
730	688	3458	4	0
729	680	3457	4	0
729	681	3458	1	0
730	684	3457	4	0
730	676	3462	4	0
765	677	3463	2	0
730	677	3460	0	0
730	687	3457	1	0
940	676	3461	0	0
758	686	3463	4	0
729	686	3461	2	0
730	688	3459	4	0
730	679	3459	4	0
755	694	3459	6	0
729	679	3462	5	0
755	693	3461	7	0
762	680	3462	0	0
755	695	3462	0	0
730	689	3463	4	0
729	675	3468	2	0
764	679	3465	7	0
755	673	3467	4	0
755	673	3466	4	0
729	686	3469	4	0
730	677	3469	7	0
730	674	3469	4	0
755	673	3464	4	0
729	677	3464	0	0
940	675	3466	0	0
730	686	3467	4	0
729	684	3469	4	0
729	685	3466	4	0
729	681	3466	0	0
729	680	3469	0	0
730	689	3464	4	0
729	681	3465	4	0
730	688	3469	4	0
729	682	3469	3	0
730	689	3466	2	0
940	688	3465	0	0
729	679	3467	7	0
729	677	3465	4	0
940	678	3468	0	0
763	679	3464	1	0
766	675	3467	4	0
734	678	3466	7	0
729	675	3465	0	0
755	672	3468	3	0
755	672	3465	2	0
730	683	3466	4	0
734	687	3468	7	0
729	687	3464	6	0
940	684	3465	0	0
729	680	3467	0	0
729	675	3464	4	0
730	683	3464	4	0
755	694	3465	6	0
734	683	3467	7	0
729	681	3464	4	0
729	688	3467	4	0
730	689	3469	6	0
734	726	3467	2	0
816	708	3472	4	0
734	688	3482	5	0
868	725	3463	0	0
734	726	3463	3	0
817	715	3481	4	0
841	722	3461	6	0
729	689	3484	5	0
729	726	3462	0	0
734	723	3469	2	0
730	726	3458	7	0
734	723	3462	0	0
800	710	3465	0	0
734	724	3459	5	0
734	727	3460	2	0
730	724	3478	7	0
734	726	3478	6	0
730	722	3481	5	0
730	725	3475	6	0
729	723	3479	6	0
729	723	3476	6	0
734	727	3474	6	0
729	721	3470	2	0
734	727	3479	6	0
734	724	3482	1	0
730	723	3472	2	0
734	721	3483	6	0
730	725	3481	6	0
940	724	3471	0	0
734	728	3450	0	0
729	729	3462	3	0
940	729	3460	0	0
730	730	3449	2	0
730	734	3452	0	0
730	734	3454	0	0
734	728	3468	2	0
730	742	3430	1	0
729	741	3429	1	0
734	741	3427	2	0
734	741	3461	5	0
209	737	3454	4	0
730	743	3426	1	0
734	741	3431	4	0
815	740	3452	2	0
730	741	3424	1	0
819	736	3446	2	0
820	739	3446	2	0
818	738	3452	0	0
839	743	3457	2	0
821	742	3446	2	0
170	738	3429	1	1
170	738	3430	1	1
940	741	3464	0	0
204	743	3465	0	0
730	740	3465	5	0
734	746	3459	0	0
734	744	3463	2	0
290	745	3457	0	0
730	747	3434	1	0
734	744	3432	2	0
940	750	3462	0	0
824	751	3446	2	0
204	748	3457	0	0
730	751	3465	5	0
22	750	3458	0	0
823	748	3446	2	0
734	751	3460	3	0
822	745	3446	2	0
940	749	3440	0	0
290	750	3467	0	0
734	749	3466	0	0
22	744	3467	0	0
170	750	3460	0	1
840	746	3470	0	0
206	748	3467	0	0
170	744	3465	0	1
170	750	3465	0	1
167	744	3460	0	1
734	749	3429	0	0
734	746	3431	0	0
729	751	3428	1	0
730	748	3425	1	0
729	748	3430	1	0
734	744	3428	6	0
730	754	3430	6	0
729	758	3439	7	0
730	753	3433	5	0
729	758	3436	5	0
940	755	3434	0	0
730	756	3432	3	0
729	759	3441	7	0
730	756	3437	2	0
170	753	3430	1	1
848	729	3479	4	0
848	736	3474	4	0
940	742	3475	0	0
848	731	3479	4	0
848	742	3477	4	0
848	742	3473	4	0
853	740	3477	0	0
848	731	3477	4	0
857	734	3476	0	0
734	731	3474	0	0
848	733	3473	4	0
848	729	3473	4	0
848	742	3479	4	0
846	738	3476	2	0
734	742	3478	4	0
848	729	3477	4	0
856	734	3479	0	0
855	738	3479	0	0
848	737	3476	4	0
858	732	3475	0	0
854	738	3472	0	0
847	732	3478	2	0
844	751	3474	2	0
734	750	3477	4	0
849	750	3476	0	0
848	749	3478	4	0
734	744	3478	4	0
848	744	3476	4	0
843	751	3481	2	0
848	751	3482	4	0
845	740	3484	2	0
859	731	3481	0	0
848	749	3482	4	0
860	744	3481	0	0
734	732	3483	4	0
734	740	3482	4	0
852	741	3481	0	0
848	738	3481	4	0
848	749	3484	4	0
848	739	3483	4	0
848	749	3480	4	0
848	735	3483	4	0
850	745	3480	0	0
848	737	3483	4	0
848	736	3481	4	0
848	741	3483	4	0
734	741	3484	4	0
848	747	3480	4	0
848	751	3480	4	0
848	747	3482	4	0
734	745	3482	2	0
851	743	3483	0	0
848	731	3483	4	0
734	744	3480	4	0
848	733	3483	4	0
734	759	3487	6	0
940	758	3487	0	0
734	757	3486	4	0
730	754	3476	5	0
730	761	3474	5	0
735	760	3484	3	0
734	761	3485	0	0
734	765	3486	3	0
940	761	3476	0	0
837	764	3464	4	0
940	761	3478	0	0
734	764	3484	2	0
836	762	3463	4	0
735	763	3486	5	0
838	766	3472	0	0
735	765	3482	3	0
734	762	3483	1	0
730	761	3458	5	0
730	765	3460	5	0
730	764	3440	0	0
729	761	3441	5	0
730	762	3443	7	0
714	775	3443	2	0
713	770	3443	2	0
831	773	3444	2	0
715	774	3443	2	0
144	768	3442	2	0
715	770	3454	2	0
734	775	3460	2	0
714	773	3454	2	0
716	771	3454	2	0
714	769	3454	2	0
734	775	3478	2	0
713	772	3454	2	0
734	774	3475	4	0
734	769	3470	2	0
734	771	3470	2	0
831	783	3446	2	0
734	780	3442	2	0
895	777	3447	0	0
714	782	3443	0	0
716	782	3442	0	0
716	777	3446	0	0
715	777	3444	0	0
713	777	3464	0	0
713	777	3441	0	0
714	777	3463	0	0
713	776	3443	2	0
715	777	3450	0	0
734	779	3461	4	0
716	777	3440	0	0
713	777	3442	0	0
734	779	3479	4	0
715	779	3472	2	0
715	780	3472	2	0
714	781	3472	2	0
714	781	3466	2	0
899	777	3460	0	0
714	780	3466	2	0
900	777	3478	0	0
715	779	3466	2	0
713	777	3443	0	0
717	777	3454	7	0
713	783	3472	2	0
714	782	3454	2	0
831	776	3467	0	0
716	783	3454	2	0
713	777	3475	0	0
714	777	3459	0	0
713	782	3466	2	0
714	777	3474	0	0
716	777	3477	0	0
734	776	3469	2	0
713	783	3466	2	0
714	777	3458	0	0
713	782	3472	2	0
831	776	3473	0	0
734	779	3468	2	0
713	781	3454	2	0
713	777	3476	0	0
715	777	3445	0	0
734	771	3484	2	0
717	777	3484	7	0
940	747	3495	0	0
714	771	3492	2	0
734	761	3488	7	0
713	775	3492	2	0
715	772	3492	2	0
734	779	3495	2	0
940	764	3495	0	0
735	756	3488	7	0
735	763	3489	1	0
940	761	3489	0	0
713	774	3492	2	0
713	777	3488	0	0
831	776	3493	0	0
734	773	3490	2	0
713	777	3494	0	0
716	773	3492	2	0
713	777	3489	0	0
714	777	3495	0	0
713	777	3490	0	0
169	755	3489	1	1
873	737	3489	2	0
169	739	3489	1	1
170	730	3489	1	1
874	733	3489	2	0
873	735	3489	2	0
864	732	3494	2	0
730	729	3490	0	0
730	726	3491	6	0
730	722	3491	3	0
730	727	3488	7	0
730	727	3493	5	0
169	725	3491	1	1
862	750	3496	6	0
862	733	3496	6	0
867	761	3500	0	0
734	763	3499	5	0
863	749	3497	2	0
734	747	3498	5	0
734	737	3498	5	0
871	761	3499	0	0
734	735	3496	5	0
866	728	3498	0	0
865	756	3502	0	0
734	754	3498	5	0
940	730	3501	0	0
865	739	3502	0	0
734	752	3496	5	0
863	766	3497	2	0
168	744	3499	0	1
940	724	3499	0	0
867	720	3501	0	0
51	718	3507	6	0
51	704	3505	0	0
51	708	3507	4	0
51	716	3505	0	0
51	711	3505	0	0
117	713	3511	5	0
51	714	3507	4	0
5	712	3511	0	0
714	775	3510	2	0
714	774	3510	2	0
715	773	3510	2	0
831	770	3511	0	0
714	771	3512	0	0
716	768	3505	0	0
831	770	3517	0	0
714	769	3510	2	0
716	768	3508	0	0
714	771	3514	0	0
716	771	3513	0	0
715	777	3499	0	0
716	781	3497	2	0
901	778	3497	6	0
714	777	3503	0	0
716	781	3510	2	0
713	783	3497	2	0
714	782	3497	2	0
714	777	3498	0	0
714	777	3496	0	0
714	782	3510	2	0
713	780	3510	2	0
713	779	3510	2	0
714	777	3497	0	0
831	776	3502	0	0
713	777	3504	0	0
904	777	3505	0	0
831	778	3511	2	0
714	777	3508	0	0
906	780	3519	6	0
713	779	3519	2	0
734	773	3429	1	0
872	765	3439	2	0
730	764	3433	1	0
729	762	3438	0	0
729	764	3437	3	0
714	777	3438	0	0
734	774	3438	2	0
714	778	3435	2	0
714	777	3439	0	0
716	777	3436	0	0
940	760	3439	0	0
730	762	3433	2	0
730	760	3435	1	0
729	762	3435	1	0
714	777	3437	0	0
713	782	3433	0	0
713	782	3437	0	0
714	782	3432	0	0
715	779	3435	2	0
715	780	3435	2	0
896	782	3439	0	0
713	782	3428	0	0
831	781	3434	6	0
714	782	3438	0	0
714	782	3431	0	0
892	782	3425	0	0
716	782	3429	0	0
734	776	3431	2	0
716	782	3424	0	0
734	777	3427	7	0
715	782	3430	0	0
886	777	3435	0	0
729	757	3418	5	0
729	755	3417	1	0
734	757	3416	3	0
730	759	3419	1	0
734	758	3420	3	0
97	760	3417	4	0
734	761	3419	3	0
830	760	3416	6	0
97	761	3416	4	0
940	762	3421	0	0
876	766	3417	0	0
722	764	3417	0	0
889	773	3418	0	0
734	772	3421	5	0
716	773	3417	6	0
875	768	3416	4	0
714	772	3417	6	0
716	777	3417	6	0
714	782	3423	0	0
891	774	3417	6	0
713	782	3422	0	0
713	782	3421	0	0
717	782	3417	6	0
734	776	3419	5	0
734	779	3421	5	0
715	778	3417	2	0
883	737	3420	2	0
729	737	3419	4	0
734	741	3413	3	0
882	737	3411	2	0
734	762	3414	3	0
734	736	3412	5	0
734	736	3409	5	0
729	760	3414	7	0
734	775	3415	5	0
734	762	3412	3	0
940	756	3415	0	0
734	758	3413	3	0
729	758	3415	6	0
734	772	3411	5	0
729	737	3410	4	0
734	778	3413	5	0
713	782	3413	0	0
716	782	3412	0	0
734	781	3412	5	0
714	782	3410	0	0
715	782	3411	0	0
886	791	3414	0	0
714	791	3412	0	0
714	791	3413	0	0
713	791	3418	0	0
888	790	3420	0	0
831	790	3419	6	0
714	786	3417	2	0
886	791	3417	0	0
922	785	3418	0	0
714	791	3423	0	0
714	791	3430	0	0
713	791	3431	2	0
888	791	3409	6	0
715	790	3431	2	0
715	791	3422	0	0
831	790	3409	6	0
831	790	3425	6	0
888	788	3430	6	0
714	790	3417	2	0
715	789	3417	2	0
714	787	3417	2	0
713	788	3417	2	0
713	791	3424	0	0
713	791	3428	0	0
888	790	3426	0	0
713	791	3429	0	0
734	785	3427	5	0
888	791	3435	0	0
831	789	3432	2	0
734	786	3424	5	0
831	795	3413	6	0
734	799	3411	2	0
888	798	3408	6	0
714	798	3411	0	0
713	799	3414	2	0
713	799	3420	2	0
888	796	3413	6	0
734	799	3418	2	0
715	797	3421	0	0
714	798	3420	2	0
831	799	3410	2	0
715	796	3416	0	0
715	798	3413	0	0
714	798	3414	2	0
714	798	3412	0	0
714	796	3417	0	0
715	796	3419	0	0
714	794	3414	2	0
734	793	3413	2	0
715	793	3414	2	0
715	795	3420	2	0
714	792	3414	2	0
714	799	3431	2	0
713	796	3435	2	0
715	797	3427	0	0
716	796	3420	2	0
715	797	3420	2	0
714	796	3418	0	0
715	797	3422	0	0
713	793	3431	2	0
716	794	3431	2	0
734	794	3427	2	0
715	797	3429	0	0
831	798	3425	2	0
888	798	3424	4	0
715	797	3426	0	0
715	797	3428	0	0
714	797	3435	2	0
734	796	3438	2	0
716	795	3435	2	0
888	797	3432	2	0
734	793	3416	2	0
734	794	3422	2	0
714	793	3420	2	0
714	794	3420	2	0
713	799	3435	2	0
716	798	3435	2	0
713	792	3431	2	0
714	795	3431	2	0
713	792	3432	0	0
831	793	3436	2	0
714	792	3433	0	0
713	794	3435	2	0
831	798	3432	2	0
713	805	3414	2	0
714	806	3414	2	0
888	803	3420	4	0
713	802	3422	0	0
831	807	3413	6	0
734	805	3413	4	0
716	802	3423	0	0
888	807	3422	0	0
713	800	3420	2	0
715	802	3418	0	0
831	803	3415	2	0
831	807	3428	6	0
888	807	3429	0	0
714	802	3429	0	0
715	807	3426	2	0
713	800	3414	2	0
888	802	3413	6	0
714	804	3414	2	0
715	802	3424	0	0
714	802	3417	0	0
715	801	3431	2	0
714	802	3428	0	0
831	801	3419	6	0
713	802	3416	0	0
898	804	3426	6	0
713	802	3431	0	0
831	803	3427	2	0
716	800	3431	2	0
714	802	3430	0	0
734	801	3430	4	0
888	801	3426	0	0
714	801	3437	0	0
887	807	3434	0	0
714	802	3432	0	0
893	801	3438	0	0
886	802	3433	4	0
831	800	3434	6	0
888	801	3434	6	0
715	811	3414	2	0
714	808	3420	0	0
716	812	3414	2	0
714	813	3414	2	0
897	808	3417	0	0
734	811	3415	4	0
713	814	3426	2	0
713	808	3416	0	0
831	809	3435	2	0
714	808	3424	0	0
715	808	3425	0	0
713	808	3432	0	0
713	808	3427	0	0
715	809	3426	2	0
831	809	3423	2	0
713	808	3431	0	0
714	810	3414	2	0
831	811	3425	6	0
888	808	3413	6	0
888	812	3425	6	0
734	810	3428	4	0
715	810	3426	2	0
886	808	3426	0	0
715	794	3445	2	0
714	797	3445	2	0
714	801	3441	0	0
734	805	3441	4	0
715	796	3445	2	0
716	795	3445	2	0
714	805	3445	2	0
717	801	3445	0	0
714	806	3445	2	0
713	793	3445	2	0
714	807	3445	2	0
734	809	3447	4	0
734	806	3444	4	0
713	812	3445	2	0
894	808	3445	6	0
713	811	3445	2	0
734	798	3448	4	0
713	785	3445	2	0
717	789	3445	0	0
713	784	3445	2	0
714	789	3452	0	0
734	790	3451	4	0
714	784	3454	2	0
716	785	3454	2	0
713	787	3454	2	0
714	786	3454	2	0
713	789	3450	0	0
713	789	3449	0	0
831	788	3453	6	0
714	789	3451	0	0
729	794	3458	1	0
729	794	3462	0	0
729	796	3457	6	0
730	796	3460	5	0
730	791	3463	0	0
715	784	3466	2	0
729	798	3465	0	0
784	798	3469	2	0
729	797	3471	4	0
715	786	3466	2	0
715	787	3466	2	0
714	785	3466	2	0
869	793	3469	2	0
730	797	3465	5	0
729	793	3464	2	0
734	806	3459	4	0
730	802	3465	5	0
729	801	3464	2	0
729	804	3462	6	0
729	806	3465	7	0
913	802	3469	2	0
730	803	3467	5	0
877	742	605	7	0
729	740	602	6	0
407	743	598	1	0
407	751	592	6	0
407	736	607	0	0
730	744	602	7	0
877	738	600	7	0
730	750	595	5	0
730	737	598	1	0
877	742	592	7	0
877	758	600	7	0
407	750	599	3	0
729	759	605	1	0
877	758	596	7	0
877	738	596	7	0
730	757	603	7	0
407	740	592	5	0
729	754	595	0	0
877	755	605	7	0
730	749	601	6	0
877	746	605	7	0
729	751	596	4	0
877	750	605	7	0
729	746	595	7	0
730	752	602	6	0
877	755	600	7	0
729	761	597	3	0
890	766	585	2	0
729	760	603	2	0
407	742	608	7	0
877	742	612	7	0
877	757	614	7	0
407	759	609	7	0
877	755	610	7	0
884	758	615	2	0
729	740	612	1	0
877	746	612	7	0
729	763	610	7	0
884	759	612	1	0
729	739	614	2	0
884	761	613	7	0
877	737	612	7	0
877	737	617	7	0
877	746	621	7	0
877	741	621	7	0
877	746	616	7	0
877	751	621	7	0
729	742	617	1	0
877	734	630	0	0
877	736	626	0	0
729	732	628	0	0
877	728	627	0	0
877	741	626	0	0
877	749	630	0	0
730	744	631	6	0
877	747	627	0	0
877	734	633	0	0
877	741	633	0	0
877	732	638	0	0
877	728	633	0	0
877	741	638	0	0
730	743	633	7	0
877	728	638	0	0
877	754	630	0	0
877	757	621	7	0
729	758	628	0	0
877	758	625	0	0
884	759	618	6	0
729	752	617	4	0
729	755	634	0	0
877	749	635	0	0
730	746	633	6	0
729	762	616	6	0
884	760	616	3	0
877	763	638	0	0
877	749	640	0	0
877	749	644	0	0
729	746	643	0	0
877	752	644	0	0
877	757	644	0	0
877	743	644	0	0
877	737	646	0	0
407	738	641	6	0
407	741	655	7	0
877	743	650	0	0
877	739	650	0	0
377	765	654	7	0
407	752	655	3	0
21	760	652	4	0
377	758	651	1	0
407	754	651	5	0
377	757	654	1	0
407	747	648	5	0
729	748	650	2	0
22	763	655	4	0
21	760	654	4	0
21	763	652	0	0
21	763	654	0	0
407	739	661	1	0
407	749	658	6	0
729	747	656	0	0
22	759	656	4	0
50	758	661	0	0
21	758	658	4	0
377	754	657	1	0
21	754	661	4	0
22	764	659	3	0
21	765	660	0	0
377	752	660	1	0
278	760	660	0	0
377	757	658	1	0
278	757	662	0	0
21	758	656	4	0
97	762	660	5	0
8	756	663	1	0
50	760	661	0	0
51	763	663	0	0
8	755	662	3	0
22	764	660	3	0
21	765	657	0	0
1	764	663	0	1
729	740	666	0	0
729	748	666	2	0
920	739	668	0	0
51	757	667	5	0
3	756	665	6	0
3	763	665	6	0
51	756	664	0	0
880	764	665	3	0
3	756	666	6	0
21	754	664	4	0
21	754	667	4	0
51	760	664	0	0
55	760	667	4	0
51	759	667	4	0
377	752	664	1	0
22	763	666	3	0
51	762	664	0	0
1	758	664	0	1
407	729	647	6	0
407	732	645	6	0
407	733	651	3	0
884	733	642	0	0
407	731	654	6	0
407	735	658	6	0
884	735	641	2	0
881	730	650	3	0
729	733	661	2	0
215	362	3351	3	0
145	364	3347	4	0
203	379	3349	4	0
51	374	3344	2	0
51	365	3358	4	0
51	370	3371	0	0
51	372	3362	2	0
275	371	3354	6	0
276	367	3353	3	0
51	373	3371	0	0
51	378	3344	6	0
203	379	3346	4	0
51	373	3351	6	0
51	378	3352	6	0
51	374	3352	2	0
20	390	3340	4	0
51	378	3348	2	0
271	373	3374	0	0
51	392	3347	6	0
51	391	3366	6	0
51	376	3354	4	0
5	376	3352	6	0
51	384	3353	2	0
51	385	3367	2	0
51	392	3354	6	0
51	392	3350	6	0
51	385	3365	2	0
145	367	3347	4	0
214	363	3355	3	0
51	368	3358	4	0
51	383	3354	6	0
51	379	3351	0	0
51	378	3355	2	0
51	382	3344	0	0
51	373	3353	6	0
276	369	3352	0	0
276	369	3355	2	0
51	382	3352	4	0
51	391	3356	6	0
51	385	3356	2	0
51	380	3356	4	0
51	377	3368	6	0
51	386	3342	2	0
20	393	3341	0	0
51	396	3347	2	0
51	396	3345	2	0
87	383	3353	0	1
51	387	3331	2	0
277	393	3331	0	0
277	390	3331	0	0
277	399	3331	0	0
277	396	3331	0	0
277	389	3325	0	0
51	387	3327	2	0
51	396	3325	6	0
51	387	3323	2	0
277	393	3321	0	0
277	395	3326	0	0
51	395	3320	6	0
51	395	3323	6	0
51	387	3319	2	0
277	390	3314	6	0
51	396	3314	6	0
277	394	3316	6	0
277	388	3316	6	0
51	386	3315	2	0
51	389	3312	0	0
51	392	3312	0	0
51	387	3304	4	0
143	341	3317	2	0
145	346	3346	2	0
51	344	3347	2	0
51	352	3354	4	0
145	353	3352	4	0
51	338	3358	2	0
51	354	3354	4	0
84	355	3353	1	1
51	351	3356	6	0
215	350	3353	6	0
51	342	3355	6	0
51	345	3357	2	0
51	340	3352	2	0
299	332	3367	6	0
51	340	3363	0	0
51	352	3366	0	0
45	349	3364	6	0
45	349	3367	6	0
45	349	3365	6	0
51	337	3366	2	0
45	349	3366	6	0
51	330	3361	0	0
46	332	3366	6	0
51	340	3365	4	0
46	325	3369	6	0
51	327	3374	4	0
45	349	3368	6	0
277	357	3374	0	0
51	337	3369	2	0
51	341	3371	2	0
51	325	3371	2	0
277	341	3372	0	0
51	330	3374	4	0
46	332	3369	6	0
51	356	3368	0	0
277	354	3371	0	0
277	349	3376	0	0
51	342	3381	2	0
51	346	3379	4	0
51	342	3377	2	0
51	351	3379	4	0
51	357	3378	4	0
277	353	3377	0	0
272	350	3390	0	0
55	351	3391	2	0
55	349	3389	1	0
55	347	3389	6	0
22	340	3391	0	0
51	342	3385	6	0
51	340	3386	2	0
272	345	3389	0	0
272	344	3389	0	0
51	340	3389	2	0
55	342	3395	3	0
22	340	3392	0	0
51	353	3394	6	0
51	353	3399	6	0
51	350	3405	6	0
51	344	3400	2	0
51	348	3405	2	0
15	341	3414	6	0
51	355	3412	6	0
51	341	3412	2	0
15	354	3414	2	0
1	345	3414	1	1
1	352	3414	1	1
15	341	3419	6	0
51	341	3417	2	0
281	341	3421	0	0
273	350	3416	6	0
63	348	3421	2	0
273	346	3416	6	0
15	354	3419	2	0
279	361	3422	4	0
281	355	3421	0	0
51	355	3417	6	0
2	345	3420	1	1
2	352	3420	1	1
25	342	3431	0	0
25	342	3428	0	0
279	361	3430	0	0
9	342	3424	0	0
9	351	3424	0	0
25	354	3431	0	0
25	354	3428	0	0
83	360	3425	0	1
83	360	3428	0	1
143	351	3436	4	0
143	343	3436	4	0
25	342	3434	0	0
25	354	3434	0	0
143	353	3436	4	0
10	349	3435	2	0
143	345	3436	4	0
10	348	3435	2	0
10	347	3435	2	0
143	407	3344	0	0
205	404	3338	4	0
205	406	3338	4	0
51	407	3340	4	0
51	407	3333	0	0
277	401	3334	0	0
273	391	3301	4	0
51	390	3300	2	0
273	398	3298	4	0
51	385	3296	2	0
97	396	3301	4	0
97	396	3300	4	0
51	385	3302	2	0
273	398	3301	4	0
43	385	3298	4	0
51	389	3299	6	0
273	400	3298	4	0
273	400	3301	4	0
51	402	3300	6	0
273	394	3298	4	0
51	400	3302	4	0
51	393	3302	4	0
273	392	3298	4	0
273	393	3301	4	0
51	409	3333	0	0
51	412	3333	0	0
51	413	3298	4	0
205	410	3335	4	0
277	410	3339	4	0
51	412	3340	4	0
5	412	3339	0	0
143	413	3347	6	0
143	413	3349	6	0
143	410	3344	0	0
143	411	3351	4	0
143	408	3351	4	0
301	421	3336	0	0
25	272	2321	6	0
25	278	2325	6	0
55	279	2321	0	0
143	274	2335	4	0
25	275	2328	6	0
25	273	2324	6	0
55	278	2322	0	0
154	278	2323	0	0
143	276	2331	0	0
55	279	2324	0	0
25	275	2324	6	0
55	277	2324	0	0
143	276	2335	4	0
6	281	2332	0	0
144	272	2332	6	0
143	274	2331	0	0
14	274	2321	2	0
14	274	2325	6	0
1	272	2325	0	1
2	279	2333	1	1
2	277	2329	0	1
22	281	2325	0	1
6	342	2523	0	0
6	339	2516	0	0
6	343	2523	0	0
6	338	2514	0	0
6	346	2523	0	0
6	338	2526	0	0
6	341	2516	0	0
6	339	2522	0	0
6	342	2524	0	0
6	345	2513	6	0
50	349	2520	0	1
2	346	2514	0	1
2	341	2518	0	1
50	341	2515	1	1
50	347	2520	0	1
50	345	2520	0	1
6	346	1582	0	0
6	340	1582	6	0
6	350	1582	0	0
6	350	1577	0	0
5	343	1579	0	0
5	342	1579	0	0
2	340	1579	0	1
1	346	1578	0	1
2	348	1580	0	1
49	345	1577	1	1
1	349	1581	0	1
228	344	3458	0	0
5	346	3457	6	0
5	338	3458	6	0
54	349	3460	1	1
54	346	3469	1	1
51	346	3467	1	1
54	348	3463	1	1
52	342	3460	0	1
53	346	3462	0	1
1	349	3458	0	1
260	149	3332	2	0
32	367	713	2	0
78	212	729	1	1
20	216	747	0	0
20	216	743	0	0
20	216	745	0	0
20	219	743	6	0
20	219	741	6	0
20	219	747	6	0
20	219	745	6	0
20	216	741	0	0
145	221	742	2	0
145	221	745	2	0
34	226	742	2	0
1	211	724	0	0
34	220	735	2	0
20	228	733	0	0
20	228	731	0	0
20	228	737	0	0
34	223	732	2	0
34	225	728	2	0
1	208	725	0	0
491	216	731	2	0
9	214	728	0	0
25	211	727	0	0
23	208	727	4	0
23	208	734	0	0
25	211	734	0	0
25	206	734	0	0
25	206	727	0	0
34	204	732	4	0
4	203	725	0	0
20	196	734	0	0
0	195	733	7	0
1	196	731	0	0
4	198	731	7	0
97	198	729	0	0
4	198	727	7	0
493	196	726	0	0
1	199	723	0	0
4	200	728	7	0
34	206	724	4	0
36	201	723	1	0
496	203	740	0	0
56	217	748	2	0
29	197	748	2	0
29	197	752	2	0
3	206	750	6	0
3	206	749	6	0
3	206	748	6	0
7	207	749	2	0
7	205	747	6	0
17	209	749	0	0
70	213	752	0	0
70	212	756	0	0
70	208	758	0	0
70	205	761	0	0
205	214	752	0	0
205	212	761	0	0
70	213	764	0	0
70	207	761	0	0
25	217	762	0	0
273	220	759	6	0
47	218	762	2	0
47	222	757	6	0
273	222	759	0	0
25	225	757	0	0
15	222	761	0	0
25	225	762	0	0
7	232	762	0	0
7	232	760	4	0
281	232	761	0	0
281	229	762	0	0
7	230	762	2	0
7	228	762	6	0
34	218	751	3	0
34	230	748	0	0
2	215	735	1	0
34	223	754	0	0
37	224	753	0	0
37	226	750	0	0
34	230	751	0	0
454	231	753	0	0
396	209	737	0	0
398	208	739	1	0
398	206	743	1	0
399	206	736	1	0
3	206	747	6	0
75	222	743	1	1
78	212	729	1	1
80	206	730	1	1
81	201	734	0	1
82	198	746	0	1
83	204	752	1	1
84	209	754	0	1
85	217	760	1	1
88	222	760	1	1
89	226	760	1	1
90	230	759	0	1
77	220	727	1	1
76	224	737	0	1
20	233	729	2	0
34	233	741	2	0
1	235	736	2	0
0	235	739	2	0
0	235	764	0	0
0	236	761	0	0
7	232	760	4	0
7	232	762	0	0
0	235	759	0	0
281	232	761	0	0
5	147	3330	6	0
227	259	3473	0	0
22	261	3474	0	0
22	257	3474	0	0
55	261	3472	0	0
55	262	3473	0	0
226	258	3471	0	0
227	259	3494	0	0
22	261	3495	0	0
22	257	3495	0	0
55	261	3493	0	0
55	262	3494	0	0
227	281	3473	0	0
22	279	3474	0	0
22	283	3474	0	0
55	284	3473	0	0
226	280	3472	4	0
102	79	1639	0	0
103	77	1639	0	0
103	76	1639	0	0
102	78	1639	0	0
194	74	1639	0	0
110	80	1642	0	0
111	80	1641	0	0
3	66	1639	0	0
3	66	1640	0	0
3	66	1641	0	0
3	66	1642	0	0
3	66	1643	0	0
118	71	1639	6	0
50	69	1639	0	0
110	80	1643	0	0
1	74	1643	0	0
1	75	1643	0	0
97	77	1643	0	0
47	510	1453	2	0
6	510	1451	0	0
6	566	1446	0	0
281	566	1442	0	0
15	563	1442	4	0
6	566	1474	0	0
28	565	1476	0	0
5	566	1480	0	0
6	566	2424	0	0
52	565	2420	2	0
22	567	2418	0	0
273	568	2420	2	0
730	584	3329	6	0
730	587	3331	0	0
986	593	3334	0	0
729	589	3353	0	0
730	590	3332	0	0
986	593	3333	0	0
940	599	3333	0	0
730	590	3329	0	0
940	594	3332	0	0
986	593	3332	0	0
730	591	3321	0	0
730	588	3325	0	0
730	596	3320	0	0
730	599	3320	0	0
730	595	3322	6	0
986	594	3321	0	0
730	586	3325	0	0
730	603	3329	0	0
986	599	3323	0	0
730	601	3331	0	0
730	602	3321	0	0
730	601	3326	0	0
729	606	3321	6	0
940	593	3324	0	0
730	607	3333	0	0
730	606	3327	6	0
730	605	3326	0	0
730	603	3332	0	0
730	601	3347	0	0
730	600	3355	0	0
22	607	3341	0	0
986	600	3332	0	0
730	601	3349	0	0
730	607	3337	0	0
730	603	3345	6	0
22	604	3345	0	0
730	600	3336	0	0
986	600	3333	0	0
22	607	3342	0	0
730	607	3345	0	0
730	594	3318	0	0
986	615	3313	0	0
730	610	3320	0	0
730	615	3342	6	0
729	615	3314	6	0
986	614	3316	0	0
986	614	3313	0	0
730	610	3317	0	0
986	613	3315	0	0
730	612	3318	0	0
22	613	3314	0	0
22	614	3314	0	0
986	608	3321	0	0
986	609	3327	0	0
729	615	3318	0	0
986	609	3328	0	0
986	613	3341	0	0
940	615	3338	0	0
730	609	3342	0	0
986	614	3337	0	0
22	622	3315	0	0
986	620	3319	0	0
940	618	3314	0	0
986	621	3316	0	0
729	621	3315	0	0
986	621	3313	0	0
986	620	3317	0	0
22	622	3314	0	0
986	617	3313	0	0
986	617	3339	0	0
986	619	3313	0	0
986	620	3320	0	0
986	619	3321	0	0
22	616	3332	0	0
22	616	3331	0	0
986	616	3313	0	0
987	620	3313	0	0
5	631	3306	4	0
51	630	3305	0	0
51	633	3299	0	0
51	632	3302	4	0
63	631	3299	6	0
24	634	3300	6	0
51	639	3301	0	0
51	630	3302	4	0
3	632	3301	4	0
51	632	3305	0	0
7	632	3300	4	0
21	639	3303	4	0
24	629	3300	2	0
51	630	3299	0	0
3	633	3301	4	0
51	639	3308	4	0
21	639	3306	4	0
63	635	3304	4	0
22	634	3303	0	1
729	592	3336	6	0
730	595	3338	0	0
730	593	3342	6	0
730	597	3337	0	0
730	590	3342	0	0
730	589	3344	0	0
730	590	3336	0	0
730	586	3345	0	0
730	587	3347	0	0
730	586	3349	0	0
730	599	3347	0	0
729	589	3348	6	0
729	590	3353	6	0
940	586	3353	0	0
730	587	3356	0	0
730	598	3355	0	0
940	599	3352	0	0
729	579	3355	6	0
729	581	3356	5	0
730	578	3352	5	0
983	577	3355	6	0
6	595	2552	0	0
51	593	2553	2	0
51	596	2553	6	0
25	619	1655	0	0
22	621	1662	6	0
51	619	1656	2	0
3	616	1662	6	0
51	618	1660	6	0
51	619	1658	2	0
145	619	1660	6	0
2	619	1661	1	1
145	619	1662	6	0
51	622	1661	6	0
7	615	1662	6	0
3	616	1661	6	0
22	622	1662	6	0
42	621	1658	0	0
5	622	1655	0	0
25	615	1659	6	0
51	622	1657	6	0
51	618	1662	6	0
18	615	2604	6	0
18	615	2605	6	0
51	621	2602	4	0
6	622	2599	0	0
51	621	2606	4	0
51	619	2603	0	0
51	619	2599	0	0
51	621	2603	0	0
51	617	2603	0	0
51	621	2599	0	0
51	619	2606	4	0
1	620	2603	0	1
51	619	2602	4	0
23	543	3276	0	0
22	545	3283	0	1
51	559	3292	6	0
51	539	3290	6	0
51	542	3283	0	0
51	557	3287	6	0
25	549	3280	2	0
25	540	3280	2	0
20	545	3273	2	0
51	548	3283	0	0
23	540	3276	0	0
281	543	3273	0	0
51	548	3285	4	0
20	545	3275	0	0
51	541	3300	2	0
5	554	3283	0	0
278	549	3275	0	0
51	541	3296	2	0
278	549	3278	0	0
51	541	3304	2	0
51	535	3292	4	0
51	535	3286	0	0
51	533	3289	2	0
277	534	3307	4	0
361	533	3305	4	0
277	531	3306	4	0
277	532	3305	4	0
51	535	3308	6	0
51	534	3305	0	0
51	531	3308	2	0
51	535	3315	4	0
51	538	3313	0	0
51	538	3315	4	0
51	535	3313	0	0
51	533	3317	6	0
51	531	3317	2	0
394	666	792	0	0
951	669	806	6	0
205	665	798	0	0
951	668	806	2	0
951	668	803	2	0
999	664	803	0	0
394	666	809	0	0
951	669	804	6	0
407	666	796	4	0
951	669	805	6	0
951	669	808	6	0
401	665	814	6	0
396	667	793	0	0
951	668	807	2	0
401	669	814	6	0
951	664	812	4	0
962	668	812	5	0
394	666	803	6	0
951	668	808	2	0
951	669	807	6	0
205	664	792	2	0
951	668	804	2	0
951	669	803	6	0
951	664	811	0	0
398	666	801	0	0
951	668	805	2	0
195	665	812	1	1
195	665	811	1	1
164	635	3564	2	0
205	632	3581	4	0
1007	644	3571	0	0
205	643	3564	4	0
164	634	3565	1	0
205	636	3568	0	0
117	635	3568	3	0
1007	642	3560	0	0
1007	642	3576	0	0
1007	645	3560	0	0
1007	640	3577	0	0
1007	643	3581	0	0
1007	645	3578	0	0
189	638	3563	0	1
1007	652	3566	0	0
1007	653	3564	0	0
1007	653	3569	0	0
1007	653	3573	0	0
205	651	3560	4	0
1007	649	3568	0	0
1007	651	3581	0	0
1007	650	3576	0	0
1008	654	3580	0	0
1007	665	3577	0	0
1007	666	3565	0	0
1007	668	3569	0	0
1007	666	3567	0	0
1007	665	3563	0	0
1007	664	3572	0	0
1007	667	3580	0	0
51	686	3550	4	0
51	680	3550	4	0
51	687	3546	4	0
116	687	3538	2	0
51	684	3538	4	0
51	680	3538	4	0
51	683	3541	0	0
1007	649	3585	0	0
1007	653	3586	0	0
1007	666	3587	0	0
1007	661	3588	0	0
1007	668	3592	0	0
1007	663	3593	0	0
205	649	3593	0	0
117	649	3592	1	0
1007	666	3595	0	0
969	651	3592	2	0
1007	643	3587	0	0
38	647	3594	1	0
192	646	3596	1	1
730	645	3605	0	0
965	642	3605	6	0
736	644	3604	0	0
1024	648	3604	4	0
940	651	3605	0	0
729	652	3606	0	0
962	647	3607	0	0
965	642	3604	0	0
205	655	3607	0	0
734	653	3611	4	0
962	655	3611	7	0
730	655	3609	0	0
755	652	3615	2	0
755	652	3613	5	0
755	649	3613	0	0
755	649	3615	7	0
755	648	3612	7	0
755	648	3613	0	0
755	650	3614	1	0
755	649	3614	7	0
205	649	3609	5	0
755	650	3613	2	0
755	650	3611	7	0
755	651	3615	1	0
755	650	3612	7	0
755	651	3613	4	0
755	651	3612	7	0
755	649	3612	2	0
755	649	3611	7	0
755	651	3614	3	0
755	652	3614	4	0
205	659	3611	7	0
965	661	3612	3	0
734	657	3608	0	0
215	664	3615	6	0
755	650	3615	1	0
734	663	3615	6	0
734	657	3615	0	0
965	663	3613	7	0
755	648	3611	6	0
734	660	3610	2	0
755	647	3611	4	0
755	647	3613	6	0
755	647	3612	0	0
1026	646	3615	0	0
755	646	3612	6	0
755	645	3612	1	0
755	646	3613	3	0
755	644	3613	0	0
965	644	3608	1	0
755	646	3611	3	0
755	643	3613	0	0
755	644	3612	7	0
734	642	3608	0	0
755	645	3613	0	0
755	645	3611	2	0
962	640	3613	2	0
205	641	3613	4	0
205	632	3593	0	0
1007	639	3585	0	0
1007	636	3594	0	0
965	639	3607	7	0
1007	639	3592	0	0
205	631	3613	3	0
736	636	3608	4	0
965	639	3613	5	0
734	636	3612	6	0
940	638	3614	0	0
965	634	3613	4	0
205	632	3610	5	0
205	632	3615	5	0
205	633	3611	0	0
734	639	3609	6	0
729	637	3615	0	0
1007	660	3561	0	0
1007	660	3583	0	0
1007	657	3566	0	0
1007	661	3579	0	0
1007	657	3578	0	0
1007	659	3574	0	0
1007	661	3570	0	0
205	661	3635	0	0
965	645	3624	6	0
729	654	3625	6	0
965	640	3626	4	0
205	654	3627	0	0
205	660	3632	0	0
965	657	3624	3	0
962	661	3624	4	0
1068	647	3645	2	0
729	648	3624	0	0
940	643	3628	0	0
734	643	3626	0	0
734	650	3629	2	0
205	660	3624	1	0
965	651	3624	0	0
734	635	3619	0	0
755	644	3618	4	0
730	634	3623	0	0
205	647	3623	5	0
755	643	3618	4	0
940	635	3621	0	0
755	647	3620	7	0
729	633	3617	0	0
755	646	3620	6	0
755	649	3618	2	0
730	655	3623	0	0
755	647	3619	2	0
965	633	3619	7	0
734	660	3617	4	0
729	637	3620	0	0
755	645	3620	5	0
965	655	3618	0	0
755	647	3618	0	0
755	652	3617	1	0
755	646	3618	7	0
734	649	3622	4	0
755	652	3616	7	0
755	649	3619	3	0
755	651	3617	7	0
755	645	3619	1	0
755	648	3620	7	0
965	638	3622	1	0
965	660	3620	3	0
755	648	3617	5	0
734	656	3619	6	0
755	652	3618	5	0
216	664	3618	0	0
755	651	3616	1	0
755	649	3617	6	0
962	654	3621	5	0
755	649	3620	0	0
755	650	3617	6	0
755	650	3620	1	0
755	650	3618	2	0
209	663	3622	4	0
217	664	3617	5	0
755	649	3616	0	0
755	650	3616	0	0
755	650	3619	3	0
755	645	3618	6	0
755	644	3619	0	0
730	662	3616	0	0
215	663	3621	1	0
730	640	3622	0	0
755	646	3619	2	0
755	648	3619	2	0
962	641	3621	3	0
729	659	3622	0	0
755	651	3618	3	0
940	658	3618	0	0
755	651	3619	4	0
755	648	3618	0	0
205	631	3617	1	0
959	495	3520	4	0
1015	492	3520	4	0
959	490	3520	2	0
959	490	3525	0	0
1017	492	3525	0	0
959	495	3525	6	0
179	492	3524	0	1
179	493	3522	0	1
179	492	3522	1	1
179	493	3524	0	1
179	494	3523	1	1
179	494	3522	1	1
179	492	3522	0	1
179	492	3523	1	1
22	10	3377	0	0
22	10	3376	0	0
22	9	3375	0	0
22	8	3376	0	0
22	6	3375	0	0
34	8	3373	0	0
164	379	3693	0	0
34	10	3374	0	0
1148	399	3710	2	0
1146	388	3701	6	0
34	6	3378	0	0
55	11	3376	5	0
55	7	3373	0	0
55	6	3374	1	0
34	8	3380	0	0
1042	12	3382	1	0
1044	15	3384	0	0
34	23	3380	0	0
34	21	3378	0	0
34	19	3380	0	0
34	22	3384	0	0
34	20	3386	0	0
1098	19	3386	0	0
1045	9	3399	0	0
1045	12	3399	0	0
1045	9	3400	0	0
1045	12	3400	0	0
1045	9	3401	0	0
1045	12	3401	0	0
1045	9	3402	0	0
1045	12	3402	0	0
1056	9	3403	0	0
1056	14	3403	0	0
402	8	3399	0	0
34	7	3397	0	0
402	8	3398	0	0
195	7	3396	0	0
205	6	3396	0	0
34	5	3396	0	0
402	6	3395	0	0
195	5	3394	0	0
1007	16	3396	0	0
1007	16	3398	0	0
1007	18	3397	0	0
1007	19	3397	0	0
195	18	3396	0	0
195	17	3395	0	0
1007	17	3400	0	0
1007	18	3401	0	0
1007	19	3401	0	0
205	18	3400	0	0
1043	7	3402	1	0
205	6	3401	0	0
34	5	3400	0	0
34	4	3399	0	0
402	5	3399	0	0
402	4	3398	0	0
1043	16	3402	1	0
1042	8	3389	0	0
1042	11	3389	1	0
1050	9	3387	0	0
1042	13	3387	0	0
277	25	3342	7	0
1075	28	3345	4	0
1075	25	3346	6	0
34	27	3339	4	0
1007	17	3352	7	0
1007	19	3349	1	0
205	28	3341	6	0
205	18	3352	4	0
277	29	3342	2	0
1075	24	3346	0	0
1007	16	3348	4	0
1097	26	3347	0	0
34	30	3342	5	0
277	26	3345	0	0
1075	30	3343	2	0
1075	30	3344	7	0
1043	16	3354	5	0
1007	19	3353	1	0
1096	13	3337	0	0
1150	27	3340	6	0
1050	12	3338	0	0
1075	29	3343	1	0
1045	9	3351	4	0
195	17	3347	2	0
34	23	3343	0	0
55	11	3328	5	0
402	8	3350	4	0
1045	9	3352	4	0
1007	16	3350	4	0
22	10	3328	0	0
1050	12	3340	3	0
1007	18	3353	5	0
1045	9	3354	4	0
1096	14	3337	0	0
1050	13	3338	6	0
34	8	3332	5	0
195	18	3348	7	0
1045	9	3353	0	0
1050	12	3339	5	0
34	19	3332	1	0
22	10	3329	0	0
1098	19	3338	0	0
1050	10	3339	5	0
1056	14	3355	2	0
402	8	3351	4	0
34	23	3332	6	0
34	22	3336	7	0
34	21	3330	0	0
34	22	3344	2	0
1056	9	3355	2	0
1045	12	3351	4	0
277	23	3345	2	0
34	20	3338	5	0
1045	12	3352	4	0
22	8	3328	0	0
1045	12	3354	4	0
1045	12	3353	4	0
34	10	3326	0	0
34	8	3325	4	0
22	9	3327	0	0
55	6	3326	2	0
195	7	3348	1	0
1007	1	3349	0	0
402	4	3350	2	0
402	6	3347	1	0
34	5	3352	7	0
402	5	3351	7	0
34	7	3349	4	0
34	5	3328	0	0
22	6	3327	0	0
34	5	3327	3	0
205	6	3353	4	0
1043	7	3354	5	0
205	6	3348	2	0
34	6	3330	0	0
55	7	3325	0	0
34	5	3348	2	0
34	4	3351	0	0
195	5	3346	6	0
205	6	3401	4	0
1075	24	3394	0	0
34	5	3396	2	0
34	5	3400	7	0
277	23	3393	2	0
1150	27	3388	6	0
1075	25	3394	6	0
205	6	3396	2	0
1043	7	3402	5	0
402	4	3398	2	0
277	26	3393	0	0
205	28	3389	6	0
1097	26	3395	0	0
402	5	3399	7	0
1045	9	3402	4	0
195	5	3394	6	0
277	25	3390	7	0
1042	13	3387	2	0
1007	1	3397	0	0
1075	29	3391	1	0
277	29	3390	2	0
1044	15	3384	4	0
1075	28	3393	4	0
402	8	3398	4	0
1075	30	3391	2	0
1045	9	3399	4	0
402	8	3399	4	0
195	18	3396	7	0
195	7	3396	1	0
55	6	3374	2	0
34	27	3387	4	0
34	8	3380	5	0
1045	12	3402	4	0
1007	18	3397	7	0
1056	9	3403	2	0
1056	14	3403	2	0
34	30	3390	5	0
34	5	3375	3	0
1045	12	3399	4	0
1045	12	3400	4	0
34	20	3386	5	0
1045	9	3400	4	0
34	23	3380	6	0
34	23	3391	0	0
1007	18	3401	5	0
34	19	3380	1	0
34	5	3376	0	0
34	22	3384	7	0
402	6	3395	1	0
1007	19	3397	1	0
1043	16	3402	5	0
1075	30	3392	7	0
1007	19	3401	1	0
1007	16	3398	4	0
1007	16	3396	4	0
34	8	3373	4	0
205	18	3400	4	0
195	17	3395	2	0
1045	12	3401	4	0
1007	17	3400	7	0
34	22	3392	2	0
34	7	3397	4	0
669	380	3692	6	0
164	379	3691	0	0
669	397	3706	5	0
23	361	3711	2	0
164	380	3693	0	0
99	394	3703	0	0
205	396	3700	3	0
164	360	3710	5	0
164	381	3691	0	0
164	381	3692	0	0
669	396	3706	1	0
965	388	3707	6	0
164	353	3654	6	0
669	358	3667	1	0
673	351	3650	0	0
674	359	3665	0	0
669	351	3668	7	0
164	354	3652	3	0
673	353	3651	0	0
673	351	3667	0	0
673	353	3667	2	0
670	357	3668	1	0
164	351	3654	5	0
669	355	3665	0	0
673	358	3661	5	0
673	356	3662	0	0
669	358	3658	2	0
673	356	3660	1	0
669	358	3660	2	0
669	352	3669	7	0
697	355	3668	5	0
673	357	3664	1	0
669	359	3666	1	0
669	356	3666	4	0
206	354	3668	5	0
673	352	3665	3	0
673	354	3669	4	0
696	343	3684	4	0
164	342	3683	4	0
164	342	3684	1	0
696	343	3683	4	0
696	344	3684	4	0
719	348	3617	0	0
164	343	3690	2	0
696	343	3689	6	0
719	349	3619	0	0
164	342	3689	6	0
164	342	3690	3	0
696	344	3689	6	0
696	344	3688	0	0
719	348	3618	4	0
719	349	3618	5	0
690	466	3676	4	0
690	466	3675	0	0
689	464	3672	0	0
673	347	3721	4	0
673	348	3723	5	0
204	357	3709	0	0
164	345	3716	3	0
164	357	3712	7	0
164	342	3716	5	0
7	357	3710	2	0
164	356	3714	0	0
164	360	3717	2	0
670	348	3708	3	0
673	352	3720	2	0
673	344	3715	4	0
673	343	3713	7	0
164	349	3720	3	0
673	360	3713	4	0
673	346	3714	6	0
673	345	3721	7	0
683	342	3710	4	0
206	356	3709	4	0
673	344	3724	0	0
673	344	3719	2	0
783	358	3709	4	0
164	350	3719	2	0
673	346	3719	6	0
673	345	3711	0	0
673	359	3715	3	0
673	359	3709	6	0
673	343	3720	1	0
682	359	3722	2	0
583	348	3608	2	0
577	348	3615	2	0
719	349	3617	0	0
719	348	3619	7	0
724	378	3622	0	0
794	377	3632	6	0
3	395	1772	0	0
776	395	1770	0	0
6	396	1770	0	0
6	395	1782	0	0
3	395	1773	0	0
23	394	1772	6	0
23	392	1769	4	0
22	394	1766	6	0
3	398	1770	0	0
29	394	1767	0	0
3	393	1770	0	0
23	392	1771	0	0
23	398	1771	0	0
3	392	1770	0	0
23	396	1772	2	0
22	393	1767	6	0
23	398	1769	4	0
3	399	1770	0	0
6	395	2714	0	0
15	391	2714	0	0
15	394	2710	6	0
15	398	2713	4	0
15	395	2717	2	0
104	81	3626	0	0
104	82	3625	0	0
104	81	3625	0	0
103	82	3618	0	0
103	81	3618	0	0
103	82	3617	0	0
105	87	3616	0	0
105	88	3616	0	0
105	86	3616	0	0
102	85	3621	0	0
102	85	3622	0	0
102	86	3621	0	0
102	86	3622	0	0
102	87	3621	0	0
102	87	3622	0	0
100	92	3620	0	0
100	91	3620	0	0
101	92	3630	0	0
101	92	3629	0	0
101	92	3628	0	0
100	91	3619	0	0
101	91	3629	0	0
103	87	3635	0	0
103	86	3636	0	0
103	87	3634	0	0
103	87	3636	0	0
102	88	3640	0	0
102	87	3640	0	0
102	88	3641	0	0
102	87	3641	0	0
963	82	3639	6	0
114	70	3629	1	0
964	77	3639	2	0
153	60	3639	0	0
153	58	3639	0	0
114	70	3628	0	0
967	71	3635	0	0
967	64	3641	1	0
180	72	3626	1	1
976	62	3639	0	0
104	81	3626	7	0
153	56	3635	2	0
1030	77	3630	4	0
1030	73	3628	2	0
976	56	3631	2	0
967	63	3640	1	0
974	56	3638	4	0
153	56	3633	2	0
164	62	3640	1	0
153	56	3636	2	0
967	73	3635	0	0
967	70	3636	0	0
1030	73	3629	0	0
153	56	3629	6	0
153	56	3637	2	0
1030	73	3627	4	0
150	64	3639	4	0
153	59	3639	0	0
114	69	3629	2	0
1030	77	3629	5	0
153	56	3632	2	0
967	64	3640	1	0
1030	77	3627	7	0
104	82	3625	6	0
153	56	3634	2	0
150	56	3628	6	0
1030	77	3625	4	0
1030	74	3631	1	0
153	63	3639	0	0
1030	73	3630	0	0
1030	75	3631	2	0
1030	77	3626	2	0
1030	76	3631	3	0
153	61	3639	0	0
103	87	3636	2	0
977	72	3636	0	0
1030	77	3624	3	0
967	70	3637	0	0
103	87	3635	1	0
290	61	3641	1	0
101	92	3629	7	0
290	61	3640	0	0
101	91	3629	1	0
102	87	3640	7	0
103	86	3636	3	0
1035	75	3628	0	0
102	88	3640	4	0
153	56	3630	2	0
1030	77	3628	6	0
101	92	3628	6	0
958	81	3633	6	0
102	87	3641	6	0
102	88	3641	5	0
115	70	3623	3	0
1030	77	3619	0	0
1030	77	3620	7	0
1030	77	3623	4	0
115	69	3623	5	0
1030	77	3622	5	0
1030	77	3621	6	0
1030	77	3618	1	0
1030	76	3616	3	0
1035	75	3619	4	0
1030	74	3617	5	0
1030	73	3618	6	0
1030	77	3617	2	0
1030	75	3616	4	0
1030	73	3620	0	0
1030	73	3619	7	0
115	70	3622	4	0
200	51	3617	0	1
107	54	3608	1	0
111	60	3615	7	0
107	53	3609	2	0
111	61	3614	1	0
107	53	3608	3	0
110	61	3610	5	0
111	60	3614	0	0
110	60	3610	6	0
110	59	3610	7	0
103	69	3610	2	0
103	68	3610	3	0
103	68	3611	4	0
102	72	3614	1	0
102	72	3613	0	0
102	73	3613	7	0
106	53	3604	3	0
106	55	3604	3	0
106	54	3604	3	0
106	54	3605	3	0
109	51	3604	0	0
108	58	3605	0	0
109	50	3605	0	0
105	65	3604	2	0
108	59	3604	7	0
109	50	3604	0	0
105	65	3605	3	0
105	66	3604	1	0
104	73	3604	5	0
104	72	3604	6	0
104	73	3605	4	0
108	58	3604	6	0
105	87	3616	2	0
102	86	3622	7	0
102	87	3621	5	0
102	87	3622	6	0
103	82	3618	4	0
102	85	3621	1	0
103	82	3617	5	0
103	81	3618	3	0
102	85	3622	2	0
105	86	3616	3	0
1023	86	1746	0	0
1005	85	1746	2	0
6	86	1743	0	0
1004	89	1744	4	0
1004	88	1747	2	0
7	86	1747	0	0
715	793	3498	0	0
714	795	3500	2	0
715	797	3500	2	0
714	793	3497	0	0
714	793	3503	0	0
714	798	3510	2	0
831	792	3501	0	0
716	793	3496	0	0
716	798	3500	2	0
714	797	3510	2	0
713	799	3500	2	0
714	796	3500	2	0
734	795	3496	4	0
715	784	3531	0	0
713	784	3523	0	0
714	784	3533	0	0
717	793	3510	2	0
734	799	3499	4	0
714	784	3521	0	0
734	798	3502	4	0
713	793	3502	0	0
717	784	3527	2	0
713	784	3535	0	0
716	784	3534	0	0
716	799	3510	2	0
714	791	3527	2	0
713	784	3512	0	0
715	786	3519	2	0
715	784	3513	0	0
714	788	3527	2	0
715	789	3527	2	0
714	784	3522	0	0
831	790	3520	2	0
714	786	3510	2	0
716	784	3514	0	0
907	777	3527	6	0
714	791	3519	2	0
714	784	3532	0	0
714	787	3519	2	0
713	784	3516	0	0
716	790	3527	2	0
831	785	3511	2	0
714	784	3517	0	0
905	787	3510	6	0
714	776	3527	2	0
831	784	3498	0	0
715	784	3515	0	0
831	778	3520	2	0
713	793	3505	0	0
715	780	3527	2	0
714	793	3519	2	0
831	783	3520	0	0
714	798	3527	2	0
713	795	3525	0	0
714	797	3527	2	0
714	792	3519	2	0
715	795	3521	0	0
716	795	3523	0	0
714	795	3524	0	0
715	795	3522	0	0
831	794	3520	0	0
713	797	3519	2	0
716	799	3519	2	0
713	793	3506	0	0
715	798	3519	2	0
713	793	3504	0	0
714	793	3527	2	0
713	792	3527	2	0
912	799	3533	4	0
716	795	3535	0	0
716	795	3531	0	0
909	795	3532	0	0
902	804	3496	6	0
714	795	3529	0	0
713	803	3503	0	0
831	794	3528	0	0
714	795	3530	0	0
713	807	3496	2	0
831	799	3528	0	0
713	803	3502	0	0
714	803	3498	0	0
714	803	3497	0	0
886	803	3496	0	0
831	802	3501	0	0
713	801	3500	2	0
713	800	3500	2	0
713	806	3510	2	0
713	805	3510	2	0
716	803	3506	0	0
714	803	3505	0	0
713	801	3510	2	0
713	803	3508	0	0
714	807	3510	2	0
715	803	3507	0	0
715	800	3510	2	0
831	802	3511	0	0
714	803	3516	0	0
714	803	3517	0	0
713	803	3504	0	0
716	803	3514	0	0
831	802	3520	0	0
715	803	3535	0	0
923	801	3527	2	0
713	803	3513	0	0
715	803	3515	0	0
714	801	3519	2	0
713	803	3512	0	0
714	802	3534	2	0
713	800	3519	2	0
831	801	3535	2	0
713	811	3496	2	0
886	803	3534	2	0
886	808	3510	2	0
714	812	3496	2	0
831	808	3497	0	0
923	809	3519	2	0
714	808	3511	0	0
714	813	3496	2	0
715	808	3517	0	0
715	808	3516	0	0
716	808	3515	0	0
908	808	3512	0	0
831	809	3520	2	0
715	808	3524	0	0
713	808	3525	0	0
716	808	3523	0	0
713	808	3521	0	0
714	808	3522	0	0
716	808	3535	0	0
713	808	3534	0	0
713	808	3533	0	0
870	808	3529	0	0
713	784	3536	0	0
713	795	3538	0	0
886	784	3537	2	0
715	795	3537	0	0
714	799	3542	2	0
870	795	3542	0	0
715	795	3536	0	0
886	805	3537	2	0
886	805	3539	4	0
714	805	3538	0	0
714	804	3542	2	0
714	804	3539	2	0
714	804	3537	2	0
714	803	3542	2	0
886	803	3539	2	0
714	803	3538	0	0
714	803	3537	2	0
714	803	3536	0	0
715	802	3542	2	0
714	802	3537	2	0
716	801	3542	2	0
714	801	3537	2	0
715	800	3542	2	0
886	800	3537	2	0
713	808	3538	0	0
870	808	3542	0	0
715	808	3536	0	0
715	800	3536	0	0
715	808	3537	0	0
734	789	3494	2	0
715	797	3490	2	0
714	793	3495	0	0
734	807	3490	6	0
717	793	3490	0	0
713	801	3490	2	0
713	804	3488	0	0
713	800	3490	2	0
714	802	3490	2	0
715	798	3490	2	0
716	799	3490	2	0
713	793	3494	0	0
831	803	3491	0	0
713	809	3488	0	0
734	810	3494	4	0
734	811	3491	4	0
716	809	3489	0	0
734	804	3494	4	0
903	809	3490	0	0
713	809	3494	0	0
713	809	3493	0	0
716	782	3543	0	0
910	782	3540	0	0
714	782	3537	2	0
715	783	3537	2	0
713	782	3539	0	0
716	780	3537	2	0
715	781	3537	2	0
713	782	3538	0	0
921	783	3549	2	0
713	782	3548	0	0
713	782	3549	0	0
713	779	3537	2	0
713	778	3537	2	0
831	781	3547	0	0
713	782	3544	0	0
713	779	3546	2	0
713	780	3546	2	0
713	776	3546	2	0
715	777	3546	2	0
172	776	3537	1	1
716	778	3546	2	0
715	774	3533	0	0
714	771	3527	2	0
713	770	3527	2	0
885	773	3538	0	0
713	774	3532	0	0
713	774	3530	0	0
714	772	3527	2	0
716	774	3543	0	0
831	773	3528	0	0
715	774	3541	0	0
713	774	3529	0	0
715	774	3542	0	0
870	774	3537	2	0
716	774	3531	0	0
173	773	3536	0	1
831	773	3547	0	0
713	774	3544	0	0
45	451	1617	2	0
42	448	1617	4	0
45	448	1620	4	0
45	451	1618	2	0
45	450	1620	4	0
45	451	1619	2	0
15	452	1620	2	0
45	450	1617	0	0
45	451	1620	2	0
45	449	1620	4	0
51	702	3507	4	0
116	706	3514	1	0
51	711	3525	6	0
51	709	3518	4	0
929	711	3522	4	0
51	712	3527	2	0
5	713	3527	0	0
51	718	3513	6	0
51	718	3527	6	0
51	713	3517	6	0
51	711	3513	2	0
117	714	3527	5	0
935	705	3520	4	0
51	718	3522	6	0
116	698	3533	6	0
51	698	3528	6	0
929	696	3527	4	0
51	698	3532	6	0
51	699	3535	2	0
116	713	3534	1	0
51	718	3533	6	0
51	696	3521	6	0
929	714	3535	2	0
51	699	3505	0	0
51	698	3522	2	0
51	697	3513	4	0
51	703	3527	4	0
51	692	3516	4	0
51	695	3508	4	0
929	695	3505	4	0
55	694	3505	2	0
734	688	3488	4	0
51	689	3505	0	0
51	691	3505	0	0
51	693	3525	4	0
55	691	3514	4	0
927	692	3515	0	0
55	691	3516	4	0
55	689	3516	6	0
51	688	3516	4	0
926	689	3513	2	0
22	688	3514	2	0
3	688	3515	2	0
51	693	3518	0	0
116	686	3507	1	0
51	684	3509	2	0
730	680	3490	5	0
51	682	3505	0	0
729	681	3488	5	0
734	680	3488	3	0
734	683	3491	0	0
51	685	3525	4	0
51	684	3513	2	0
729	686	3489	5	0
51	684	3517	2	0
734	679	3495	0	0
734	678	3491	0	0
937	678	3508	4	0
729	676	3490	5	0
734	677	3488	0	0
51	675	3505	0	0
729	679	3493	5	0
116	673	3507	2	0
116	681	3534	4	0
116	688	3532	1	0
51	683	3534	4	0
929	689	3536	4	0
51	697	3538	6	0
51	707	3542	2	0
51	699	3543	2	0
51	690	3542	4	0
55	708	3537	0	0
51	706	3537	2	0
55	711	3539	2	0
929	709	3537	4	0
51	692	3550	4	0
51	694	3546	4	0
51	701	3550	4	0
51	697	3544	6	0
51	707	3550	4	0
51	711	3550	4	0
117	718	3542	1	0
51	718	3538	6	0
51	716	3550	4	0
930	717	3542	4	0
51	712	3546	4	0
55	712	3540	4	0
730	795	3475	7	0
729	793	3476	3	0
730	793	3474	3	0
729	793	3479	5	0
729	796	3478	7	0
713	784	3472	2	0
715	787	3472	2	0
714	785	3472	2	0
729	797	3472	5	0
715	786	3472	2	0
729	790	3474	4	0
730	791	3477	5	0
734	786	3482	2	0
734	785	3480	2	0
734	796	3487	4	0
734	794	3485	0	0
730	795	3481	1	0
729	804	3477	0	0
729	806	3474	1	0
730	802	3477	4	0
713	807	3485	2	0
730	802	3472	5	0
730	806	3472	5	0
729	800	3473	3	0
715	805	3485	2	0
713	806	3485	2	0
714	804	3486	0	0
734	806	3482	4	0
713	804	3487	0	0
886	804	3485	0	0
734	802	3487	4	0
831	808	3486	0	0
713	809	3487	0	0
110	686	3268	0	0
699	712	3301	5	0
109	715	3265	0	0
699	713	3301	5	0
610	714	3279	6	0
610	716	3286	7	0
637	707	3297	0	0
610	710	3292	7	0
610	705	3297	6	0
109	716	3265	6	0
610	706	3292	7	0
609	708	3296	7	0
610	713	3281	2	0
107	717	3266	6	0
609	713	3273	2	0
610	716	3279	6	0
610	716	3273	2	0
609	714	3288	6	0
610	714	3276	0	0
609	703	3293	4	0
609	688	3293	2	0
110	707	3268	0	0
114	707	3273	0	0
609	711	3287	6	0
609	696	3285	4	0
616	703	3283	0	0
609	710	3275	3	0
610	710	3271	7	0
610	702	3284	5	0
638	701	3280	2	0
610	703	3290	4	0
610	696	3291	1	0
639	701	3279	6	0
610	704	3274	1	0
110	702	3267	0	0
110	705	3270	0	0
110	694	3276	0	0
609	704	3280	2	0
609	692	3285	0	0
112	690	3269	0	0
610	692	3282	7	0
112	692	3270	0	0
610	693	3288	1	0
110	691	3272	0	0
610	692	3292	3	0
110	688	3271	0	0
110	688	3276	0	0
112	689	3267	0	0
609	690	3280	3	0
610	694	3283	6	0
112	691	3265	0	0
610	691	3277	6	0
610	689	3281	3	0
609	695	3288	3	0
110	689	3274	0	0
102	681	3278	0	0
610	683	3280	6	0
195	680	3267	0	0
195	682	3267	0	0
110	684	3270	0	0
195	682	3265	0	0
111	683	3275	2	0
110	685	3278	0	0
610	687	3303	0	0
102	681	3276	0	0
195	682	3270	0	0
610	683	3298	3	0
110	686	3272	0	0
110	685	3275	0	0
114	683	3273	0	0
102	680	3273	0	0
114	680	3269	0	0
699	711	3310	1	0
699	716	3305	4	0
699	709	3304	7	0
699	714	3307	2	0
699	712	3310	1	0
699	709	3308	0	0
699	715	3306	3	0
699	709	3309	0	0
102	679	3275	0	0
102	678	3272	0	0
102	677	3275	0	0
107	679	3277	7	0
107	676	3272	6	0
114	678	3268	0	0
102	677	3270	0	0
51	668	3269	0	0
474	666	3271	6	0
478	666	3276	2	0
473	664	3271	6	0
51	666	3278	4	0
483	665	3270	4	0
484	665	3273	0	0
477	664	3276	2	0
476	668	3273	4	0
63	667	3279	6	0
51	664	3278	4	0
51	662	3269	0	0
51	662	3278	4	0
475	662	3273	0	0
734	680	3483	0	0
729	687	3485	5	0
730	685	3483	5	0
730	685	3486	5	0
734	683	3486	2	0
730	679	3487	5	0
0	213	3562	0	0
0	230	3563	0	0
0	225	3566	0	0
0	228	3564	0	0
0	214	3561	0	0
0	227	3565	0	0
618	221	3567	0	0
0	215	3563	0	0
0	229	3562	0	0
0	224	3559	0	0
0	227	3560	0	0
0	226	3567	0	0
407	487	612	0	0
193	484	610	0	0
22	485	615	0	0
22	484	613	0	0
3	485	618	0	0
194	480	614	0	0
193	485	610	0	0
22	485	617	0	0
99	486	612	3	0
3	484	615	0	0
193	493	621	0	0
458	494	618	2	0
22	482	616	0	0
21	492	618	4	0
194	480	617	0	0
21	495	621	0	0
21	492	621	0	0
21	492	616	0	0
193	488	616	0	0
21	489	616	0	0
21	489	618	4	0
193	491	616	0	0
193	490	621	0	0
193	488	618	0	0
21	489	621	0	0
1	487	619	1	1
1	483	617	0	1
1	485	619	1	1
1	487	615	1	1
45	487	1557	4	0
7	486	1559	2	0
45	490	1557	4	0
3	485	1559	0	0
22	487	1560	2	0
45	489	1557	4	0
45	491	1565	0	0
45	491	1562	2	0
45	494	1565	0	0
6	494	1562	2	0
45	495	1557	4	0
45	490	1565	0	0
45	495	1565	0	0
45	489	1565	0	0
45	492	1561	6	0
45	488	1557	4	0
45	493	1565	0	0
45	491	1563	0	0
45	492	1565	0	0
45	487	1565	0	0
45	496	1562	2	0
45	492	1562	6	0
453	493	1558	4	0
45	492	1563	0	0
45	491	1561	2	0
45	496	1559	2	0
45	488	1565	0	0
45	496	1561	2	0
45	496	1560	2	0
45	491	1560	2	0
45	496	1563	2	0
45	496	1558	2	0
45	496	1564	2	0
1	487	1563	1	1
124	487	1558	1	1
1	484	1561	0	1
45	463	1462	4	0
45	459	1462	6	0
45	462	1462	4	0
45	463	1459	0	0
143	462	1457	0	0
22	460	1458	6	0
45	460	1462	4	0
45	461	1462	4	0
45	460	1459	0	0
44	461	1459	4	0
143	462	1468	4	0
22	459	1458	6	0
143	466	1459	6	0
45	464	1460	2	0
45	464	1461	2	0
45	459	1459	6	0
143	458	1465	2	0
143	466	1466	6	0
45	466	1463	6	0
45	464	1459	2	0
45	459	1460	6	0
45	459	1461	6	0
45	464	1462	2	0
43	461	1464	4	0
45	466	1464	6	0
143	458	1462	2	0
143	467	1462	6	0
6	467	1464	4	0
25	459	2407	2	0
45	459	2408	0	0
44	461	2408	4	0
45	463	2408	0	0
45	460	2408	0	0
25	464	2407	2	0
25	459	2406	2	0
25	464	2406	2	0
45	464	2408	0	0
70	392	7	0	0
25	415	35	2	0
70	410	33	0	0
205	413	20	0	0
191	407	24	0	0
191	411	22	0	0
191	409	27	0	0
70	400	2	0	0
191	408	22	0	0
70	389	31	0	0
205	405	28	0	0
191	405	25	0	0
205	414	27	0	0
72	400	38	1	0
70	395	20	0	0
70	389	12	0	0
70	395	30	0	0
70	400	9	0	0
70	397	15	0	0
70	389	2	0	0
70	404	32	0	0
191	411	25	0	0
70	399	14	0	0
70	387	3	0	0
70	389	21	0	0
7	422	31	4	0
70	420	19	0	0
70	422	7	0	0
7	421	31	4	0
70	423	22	0	0
59	416	24	0	0
7	423	31	4	0
70	419	15	0	0
7	420	31	4	0
70	428	17	0	0
3	426	15	0	0
70	427	27	0	0
1	424	18	1	1
9	423	35	2	0
7	417	39	0	0
7	419	39	0	0
7	422	37	6	0
9	417	38	0	0
9	420	32	0	0
7	420	33	0	0
7	422	33	0	0
7	420	37	4	0
43	417	34	4	0
7	419	32	6	0
7	422	35	6	0
10	424	38	4	0
10	424	36	4	0
274	425	33	2	0
70	413	44	0	0
25	421	40	4	0
25	415	40	2	0
70	410	43	0	0
70	395	43	0	0
70	389	43	0	0
14	422	976	4	0
14	421	981	2	0
45	419	980	2	0
45	419	978	2	0
45	416	980	6	0
5	419	976	0	0
45	419	979	2	0
45	418	980	4	0
45	417	980	4	0
44	417	978	4	0
281	421	976	0	0
45	416	978	6	0
2	420	982	0	1
45	416	979	6	0
2	415	982	0	1
2	421	979	1	1
6	419	1920	0	0
25	417	1925	3	0
25	419	1923	3	0
25	417	1923	1	0
281	418	1924	3	0
25	419	1925	1	0
398	403	901	2	0
394	407	902	6	0
395	376	899	0	0
395	415	901	6	0
403	402	898	2	0
401	398	897	2	0
1176	409	900	0	0
404	398	901	2	0
32	406	899	2	0
394	409	898	6	0
396	415	899	6	0
399	377	903	6	0
399	385	897	7	0
32	412	899	6	0
404	390	901	7	0
397	400	903	2	0
395	387	899	6	0
1176	394	899	0	0
394	379	896	3	0
401	386	902	7	0
395	400	897	2	0
394	400	899	2	0
1176	411	896	0	0
32	389	896	6	0
397	395	902	2	0
403	391	898	7	0
396	396	898	2	0
33	383	899	0	0
395	379	899	7	0
1176	413	902	0	0
1169	395	896	3	0
397	411	903	6	0
405	414	897	6	0
399	422	900	6	0
401	423	896	6	0
430	423	903	0	0
395	418	901	6	0
394	420	897	6	0
394	421	903	6	0
1176	417	898	0	0
32	420	899	6	0
398	429	895	6	0
399	431	902	6	0
1113	427	899	2	0
396	424	899	6	0
32	426	894	7	0
397	429	891	4	0
397	431	890	2	0
398	426	891	6	0
32	428	890	7	0
33	428	893	7	0
405	430	893	5	0
404	427	888	2	0
398	431	897	6	0
32	429	903	6	0
396	425	888	7	0
394	430	888	6	0
1176	430	899	0	0
404	437	890	4	0
397	432	899	5	0
405	432	895	4	0
395	438	893	5	0
1176	434	891	0	0
403	435	888	4	0
404	437	902	5	0
396	435	895	5	0
400	439	896	5	0
397	432	889	4	0
403	439	899	5	0
394	435	899	5	0
398	434	902	5	0
395	444	891	2	0
398	442	893	2	0
1113	442	899	2	0
397	446	890	2	0
398	445	893	2	0
403	446	897	3	0
395	455	893	7	0
430	453	902	0	0
394	448	899	2	0
399	448	895	2	0
396	448	889	2	0
405	452	898	1	0
1176	450	893	0	0
399	453	890	7	0
396	449	902	7	0
395	455	901	2	0
396	455	898	2	0
401	451	900	0	0
1176	461	898	0	0
1176	461	892	0	0
430	457	895	0	0
32	460	901	2	0
32	457	889	7	0
404	460	889	7	0
398	459	893	7	0
1113	462	903	2	0
396	458	902	2	0
32	463	900	2	0
395	463	895	7	0
1169	463	889	4	0
400	459	891	7	0
397	458	897	2	0
405	466	891	0	0
405	465	902	4	0
404	471	900	4	0
397	468	898	4	0
398	465	899	4	0
396	466	895	4	0
397	465	897	4	0
396	469	896	4	0
395	468	894	4	0
397	470	889	0	0
399	469	901	4	0
398	472	898	4	0
398	475	890	0	0
403	472	894	4	0
399	473	903	4	0
394	474	901	4	0
1176	465	909	0	0
398	456	906	2	0
395	460	906	2	0
404	468	907	4	0
394	462	907	2	0
32	470	909	5	0
403	466	906	4	0
395	464	904	2	0
400	469	904	4	0
1176	460	909	0	0
404	457	904	2	0
401	471	905	4	0
397	460	904	2	0
397	474	906	4	0
396	475	904	4	0
1176	454	907	0	0
397	448	907	4	0
399	455	904	2	0
403	451	908	5	0
32	452	905	2	0
32	454	909	2	0
32	447	909	3	0
404	447	904	6	0
398	444	908	0	0
399	440	907	2	0
1176	441	909	0	0
395	442	906	7	0
399	438	904	5	0
401	435	906	6	0
33	425	904	6	0
405	425	906	6	0
33	430	906	6	0
394	416	904	6	0
398	420	905	6	0
32	418	905	6	0
399	416	907	6	0
33	422	906	6	0
1176	412	891	0	0
401	410	885	7	0
400	412	887	7	0
405	409	882	7	0
394	412	889	7	0
395	411	894	6	0
33	415	882	2	0
33	413	893	6	0
397	412	882	2	0
399	414	880	2	0
396	413	884	7	0
396	412	906	6	0
397	408	886	2	0
32	409	905	6	0
1176	414	908	0	0
397	410	891	7	0
33	414	905	6	0
398	410	907	6	0
397	414	886	7	0
404	409	888	7	0
32	405	890	2	0
1176	401	892	0	0
32	403	883	7	0
395	401	886	2	0
395	406	885	7	0
396	403	890	2	0
403	407	880	2	0
400	401	881	2	0
33	405	882	7	0
395	404	887	2	0
1113	406	894	2	0
405	405	905	2	0
1176	401	883	0	0
396	407	889	7	0
403	401	889	2	0
32	400	905	2	0
32	397	882	2	0
404	398	886	2	0
399	399	883	2	0
430	395	883	0	0
33	394	882	2	0
32	393	904	6	0
394	395	892	4	0
394	397	889	2	0
397	396	886	2	0
397	392	889	6	0
399	396	905	2	0
405	393	884	4	0
394	394	880	3	0
398	392	906	7	0
398	398	891	2	0
394	399	894	2	0
395	394	887	4	0
403	390	880	5	0
394	386	904	7	0
394	391	904	6	0
33	385	885	4	0
404	388	889	5	0
32	384	905	6	0
33	387	908	7	0
32	389	883	4	0
398	386	883	6	0
1176	385	894	0	0
397	386	880	6	0
399	391	886	6	0
1113	388	886	2	0
403	391	892	6	0
399	385	889	6	0
396	388	905	7	0
396	388	893	6	0
395	378	887	3	0
1176	380	886	0	0
32	383	890	0	0
33	381	895	1	0
397	380	889	2	0
394	381	891	3	0
1176	380	882	0	0
403	376	880	3	0
403	376	891	5	0
32	381	904	7	0
400	383	884	5	0
33	378	893	0	0
1176	379	908	0	0
395	369	892	7	0
394	371	888	5	0
396	371	893	6	0
397	370	890	2	0
404	373	890	3	0
405	375	895	7	0
403	373	903	3	0
403	375	888	3	0
34	368	902	2	0
403	372	900	1	0
398	371	898	5	0
397	372	896	6	0
430	375	893	0	0
394	372	902	4	0
1176	368	907	0	0
396	368	898	4	0
401	375	898	2	0
32	371	905	4	0
1176	375	905	0	0
404	375	901	5	0
1176	367	892	0	0
397	362	892	7	0
395	362	904	0	0
394	364	900	0	0
398	363	888	7	0
1169	367	888	4	0
430	360	893	0	0
32	364	896	0	0
397	365	905	0	0
395	362	897	7	0
395	366	895	3	0
404	360	889	7	0
405	352	899	5	0
394	357	895	7	0
403	357	888	7	0
1113	359	901	4	0
1176	358	892	0	0
32	352	896	5	0
396	355	902	0	0
405	353	908	7	0
400	355	890	7	0
405	359	896	7	0
32	352	891	5	0
395	355	897	7	0
394	357	906	0	0
399	355	893	7	0
394	351	893	5	0
1176	351	888	0	0
1176	347	903	0	0
394	345	898	5	0
401	347	888	5	0
1176	350	906	0	0
32	348	899	5	0
395	344	903	7	0
32	351	904	7	0
395	350	901	7	0
394	347	894	5	0
1113	349	891	4	0
32	344	908	7	0
395	343	892	5	0
430	342	895	5	0
398	341	892	5	0
32	339	902	7	0
1176	342	905	0	0
397	341	901	7	0
398	342	907	7	0
395	341	888	5	0
32	342	899	5	0
33	340	898	5	0
1	506	33	0	0
191	505	24	0	0
37	509	36	0	0
1	485	31	0	0
70	483	3	0	0
191	504	25	3	0
25	511	35	2	0
63	510	37	4	0
191	507	25	0	0
191	509	20	0	0
191	504	29	7	0
191	509	29	7	0
191	505	27	0	0
191	504	22	0	0
191	510	27	0	0
72	496	36	3	0
191	510	25	5	0
191	507	29	7	0
72	496	38	3	0
191	507	27	6	0
37	505	11	0	0
191	509	22	0	0
37	499	33	0	0
72	497	38	3	0
191	507	22	4	0
72	496	37	3	0
191	507	19	5	0
191	509	24	0	0
191	505	21	0	0
72	497	37	3	0
191	501	25	0	0
1	500	32	0	0
72	493	35	3	0
191	503	27	7	0
70	496	2	0	0
37	499	10	0	0
72	494	35	3	0
1	496	9	0	0
191	501	28	0	0
191	503	24	2	0
192	489	35	0	0
72	493	36	3	0
72	495	35	3	0
72	495	36	3	0
192	493	25	0	0
72	494	37	3	0
72	492	35	3	0
72	495	37	3	0
192	500	38	0	0
72	494	36	3	0
72	495	38	3	0
1	491	30	0	0
70	485	2	0	0
1	493	15	0	0
37	495	28	0	0
1	495	14	0	0
1	485	12	0	0
37	487	32	0	0
1	488	7	0	0
37	492	10	0	0
1	491	20	0	0
1	485	21	0	0
37	516	11	0	0
37	518	29	0	0
37	518	24	0	0
1	516	19	0	0
37	515	20	0	0
1	518	7	0	0
7	518	31	4	0
1	519	22	0	0
7	519	31	4	0
1	515	15	0	0
7	517	31	4	0
7	518	35	6	0
7	516	31	4	0
7	518	37	6	0
7	518	33	0	0
37	516	26	0	0
59	512	24	0	0
9	519	35	2	0
7	513	39	0	0
7	516	33	0	0
9	516	32	0	0
7	516	37	4	0
9	513	38	0	0
7	515	39	0	0
43	513	34	4	0
7	515	32	6	0
3	522	15	0	0
1	523	27	0	0
1	524	17	0	0
1	520	18	1	1
10	520	36	4	0
10	520	38	4	0
274	521	33	2	0
25	517	40	4	0
70	509	44	0	0
70	506	43	0	0
25	511	40	2	0
1	491	43	0	0
1	485	43	0	0
44	513	978	4	0
45	513	980	4	0
45	514	980	4	0
45	515	979	2	0
5	515	976	0	0
281	517	976	0	0
45	515	980	2	0
14	517	981	2	0
14	518	976	4	0
45	512	980	6	0
45	512	979	6	0
45	515	978	2	0
45	512	978	6	0
2	511	982	0	1
2	517	979	1	1
2	516	982	0	1
403	421	891	3	0
405	419	894	5	0
399	423	890	0	0
1163	418	891	0	0
1163	416	890	0	0
1163	418	888	0	0
34	419	888	0	0
398	421	888	0	0
32	417	895	6	0
394	423	894	4	0
582	417	889	0	0
25	515	1923	3	0
6	515	1920	0	0
25	515	1925	1	0
25	513	1923	1	0
25	513	1925	3	0
281	514	1924	7	0
634	709	2364	0	0
646	710	2364	2	0
145	567	1535	2	0
14	560	1526	6	0
25	563	1525	6	0
145	567	1533	2	0
257	567	1532	0	0
25	565	1525	6	0
45	565	1531	6	0
25	565	1529	6	0
45	565	1533	6	0
25	563	1529	6	0
45	565	1532	6	0
45	565	1534	0	0
45	565	1530	6	0
44	563	1531	0	0
45	562	1531	2	0
45	562	1532	2	0
45	562	1534	0	0
45	562	1533	2	0
332	560	1531	5	0
215	567	1536	0	0
27	561	1530	0	0
24	560	1530	0	0
45	562	1530	2	0
299	565	1536	0	0
6	567	1544	0	0
299	560	1536	0	0
217	564	1536	0	0
121	563	1536	6	0
1	566	1530	0	1
565	463	3399	4	0
1152	465	3398	6	0
565	469	3403	6	0
730	468	3398	0	0
730	468	3392	6	0
729	471	3398	7	0
1153	470	3401	0	0
1152	465	3400	6	0
1153	468	3401	0	0
51	485	3391	4	0
1154	473	3398	2	0
729	474	3385	7	0
730	473	3388	4	0
1154	473	3400	2	0
666	486	3389	0	0
730	474	3383	0	0
51	485	3388	0	0
565	475	3399	0	0
51	489	3391	4	0
5	490	3393	0	0
5	490	3389	2	0
51	489	3388	0	0
25	397	3264	0	0
51	397	3286	4	0
27	399	3272	6	0
25	396	3275	0	0
51	407	3282	0	0
270	392	3271	6	0
5	399	3276	0	0
51	399	3282	0	0
51	406	3264	0	0
51	411	3269	4	0
51	407	3268	4	0
51	401	3274	0	0
270	394	3282	0	0
10	395	3274	2	0
25	399	3271	6	0
51	416	3282	2	0
51	414	3276	2	0
25	400	3266	6	0
51	403	3269	0	0
25	392	3269	6	0
51	402	3278	0	0
270	398	3280	0	0
5	423	3285	0	0
51	418	3278	2	0
51	405	3274	0	0
51	412	3280	0	0
51	408	3277	0	0
51	410	3264	0	0
51	418	3283	2	0
51	408	3273	4	0
51	404	3282	4	0
51	423	3268	6	0
51	409	3285	4	0
51	411	3273	4	0
51	403	3285	4	0
270	419	3286	0	0
51	400	3295	0	0
51	403	3294	4	0
51	395	3292	0	0
51	407	3292	0	0
63	397	3294	2	0
51	401	3292	0	0
51	393	3295	0	0
51	392	3292	0	0
51	398	3292	0	0
51	413	3295	0	0
51	416	3295	0	0
51	419	3295	0	0
51	396	3295	0	0
51	395	3294	4	0
51	404	3292	0	0
51	426	3273	6	0
51	429	3292	6	0
51	425	3283	6	0
51	428	3290	0	0
51	424	3278	6	0
270	424	3289	0	0
43	426	3290	4	0
51	419	3298	4	0
51	423	3298	4	0
6	375	1382	0	0
19	369	1381	2	0
110	354	3281	0	0
210	353	3277	0	0
110	356	3276	0	0
58	374	3276	6	0
110	362	3283	0	0
108	373	3281	4	0
110	360	3278	0	0
110	372	3281	4	0
110	375	3280	0	0
41	368	3270	4	0
106	371	3281	4	0
478	648	3272	2	0
51	648	3274	4	0
51	646	3274	4	0
51	644	3274	4	0
477	646	3272	2	0
23	629	3295	6	0
51	631	3293	0	0
486	649	3275	6	0
507	634	3295	4	0
51	632	3293	0	0
474	648	3267	6	0
476	650	3269	4	0
51	650	3265	0	0
485	647	3269	0	0
473	646	3267	6	0
475	644	3269	0	0
483	647	3266	4	0
51	644	3265	0	0
6	60	2328	2	0
433	589	3409	0	0
210	450	3710	1	1
210	452	3712	3	1
210	450	3705	1	1
210	456	3710	1	1
210	456	3706	1	1
210	456	3709	2	1
210	456	3708	1	1
1156	471	3708	0	0
37	126	496	0	0
34	126	495	0	0
1160	441	3702	2	0
210	456	3707	1	1
210	456	3704	1	1
210	452	3703	2	1
1184	440	3709	0	0
1144	462	3703	0	0
1178	472	3698	0	0
210	454	3711	2	1
1162	460	3705	4	0
1178	471	3711	0	0
7	458	3701	4	0
1165	474	3719	6	0
210	454	3704	3	1
210	455	3705	3	1
1032	452	3708	0	0
210	452	3704	0	1
210	454	3704	0	1
210	448	3708	3	1
210	455	3710	2	1
210	450	3708	1	1
210	450	3704	1	1
210	455	3704	0	1
210	451	3711	3	1
210	453	3704	0	1
210	450	3707	1	1
210	450	3706	1	1
210	449	3709	3	1
210	451	3704	2	1
210	450	3709	1	1
210	453	3712	0	1
210	451	3704	0	1
210	454	3712	0	1
210	450	3711	1	1
210	455	3712	0	1
210	453	3712	2	1
210	450	3705	2	1
210	450	3704	0	1
210	450	3712	0	1
210	452	3712	0	1
210	450	3710	3	1
210	448	3707	2	1
210	451	3712	0	1
210	456	3711	1	1
210	449	3706	2	1
210	457	3707	3	1
210	457	3708	2	1
210	456	3705	1	1
1178	463	3713	0	0
1178	474	3701	0	0
1161	458	3702	6	0
931	450	3702	7	0
1178	468	3710	0	0
47	460	3708	4	0
1178	464	3710	0	0
1178	467	3712	0	0
1165	474	3715	2	0
1178	474	3714	0	0
1185	440	3713	0	0
210	453	3703	3	1
1159	446	3698	0	0
1033	440	3718	6	0
1117	441	3705	0	0
1178	467	3706	0	0
47	452	3700	7	0
1158	462	3699	3	0
1178	472	3699	0	0
51	465	3724	0	0
51	467	3722	0	0
1037	471	3722	0	0
212	464	3721	2	1
212	466	3723	2	1
1037	464	3730	0	0
1037	471	3734	0	0
1037	474	3730	0	0
1037	469	3728	0	0
211	456	3728	3	1
1037	466	3739	0	0
1037	460	3737	0	0
1168	462	3739	0	0
5	98	2930	0	0
5	137	2932	0	0
41	515	3370	6	0
42	508	1479	4	0
41	516	1479	0	0
42	516	2423	0	0
51	511	3371	6	0
51	511	3374	6	0
51	503	3373	2	0
51	503	3371	2	0
51	507	3369	0	0
51	505	3369	0	0
51	508	3380	2	0
51	508	3389	2	0
51	508	3385	2	0
51	519	3375	2	0
51	518	3368	0	0
51	519	3388	6	0
51	519	3384	6	0
51	511	3396	4	0
51	508	3394	2	0
51	514	3394	6	0
205	508	3369	0	0
205	507	3375	0	0
205	506	3370	0	0
205	510	3376	0	0
20	510	3386	6	0
20	510	3393	6	0
730	519	3390	0	0
730	515	3391	0	0
730	512	3386	0	0
729	518	3387	0	0
729	515	3386	0	0
729	514	3390	0	0
729	513	3387	0	0
29	509	2421	2	0
29	509	2426	2	0
1	514	2426	1	1
1	514	2424	0	1
29	512	1482	0	0
1	513	1479	0	1
421	690	1424	6	0
1155	446	3374	0	0
6	642	1541	0	0
227	281	3494	0	0
22	283	3495	0	0
22	279	3495	0	0
55	283	3493	0	0
55	284	3494	0	0
1167	427	3707	0	0
1114	426	3703	0	0
1123	423	3701	2	0
558	422	3702	0	0
559	421	3702	0	0
560	420	3702	0	0
561	419	3702	0	0
1124	418	3703	4	0
1125	420	3706	6	0
99	413	3700	0	0
99	406	3702	0	0
99	408	3702	1	0
965	404	3703	4	0
99	408	3705	1	0
99	415	3710	6	0
99	411	3710	6	0
205	413	3707	5	0
729	401	3716	0	0
730	400	3715	0	0
729	400	3718	0	0
730	398	3717	0	0
730	399	3720	0	0
164	394	3714	0	0
164	393	3718	0	0
729	395	3721	0	0
965	390	3719	0	0
99	387	3716	0	0
99	387	3722	0	0
99	391	3723	0	0
730	396	3723	0	0
1145	394	3726	0	0
99	404	3733	0	0
1145	394	3732	4	0
99	394	3735	0	0
673	393	3737	0	0
99	389	3737	0	0
99	388	3740	0	0
99	402	3740	0	0
673	399	3738	0	0
673	407	3729	0	0
99	410	3731	0	0
205	409	3740	0	0
673	407	3736	0	0
205	412	3729	0	0
673	413	3733	0	0
673	417	3733	0	0
1116	414	3725	0	0
673	410	3723	0	0
729	413	3721	0	0
99	418	3720	0	0
205	420	3723	0	0
673	422	3726	0	0
1116	424	3723	0	0
673	423	3719	0	0
673	428	3720	0	0
729	428	3732	0	0
730	423	3732	0	0
99	428	3737	0	0
99	427	3736	0	0
99	427	3738	0	0
1116	411	3737	0	0
729	424	3736	0	0
223	312	3348	0	0
22	272	742	0	0
22	274	742	0	0
22	275	742	0	0
193	273	738	0	0
193	275	739	0	0
193	278	738	0	0
193	279	740	6	0
193	280	740	6	0
193	280	739	6	0
193	279	744	6	0
193	280	744	6	0
193	280	745	6	0
1179	266	741	4	0
1094	269	743	0	0
1093	269	741	0	0
1094	271	744	0	0
1093	271	740	0	0
1102	277	743	6	0
1101	277	741	6	0
1180	278	741	4	0
1181	280	741	4	0
22	247	729	0	0
22	249	729	0	0
22	250	729	0	0
193	249	726	0	0
193	248	724	0	0
193	252	724	0	0
193	255	726	6	0
193	255	727	6	0
193	254	727	6	0
193	246	734	0	0
193	249	732	0	0
193	253	734	0	0
193	255	732	0	0
193	255	731	6	0
193	254	731	6	0
1179	241	728	4	0
1094	244	730	0	0
1093	244	728	0	0
1093	246	727	0	0
1094	246	731	0	0
1101	252	728	6	0
1102	252	730	6	0
1180	253	728	4	0
1181	255	728	4	0
1069	249	759	4	0
1070	247	758	6	0
1070	250	755	2	0
1070	251	762	5	0
1070	254	756	0	0
1070	254	758	5	0
1070	256	758	4	0
1070	259	757	5	0
22	320	742	0	0
22	322	742	0	0
22	323	742	0	0
193	321	738	0	0
193	323	739	0	0
193	326	738	0	0
193	328	739	6	0
193	328	740	6	0
193	327	740	6	0
193	328	745	6	0
193	328	744	6	0
193	327	744	6	0
1179	314	741	4	0
1093	317	741	0	0
1094	317	743	0	0
1093	319	740	0	0
1094	319	744	0	0
1101	325	741	6	0
1102	325	743	6	0
1180	326	741	4	0
1181	328	741	4	0
22	295	729	0	0
22	297	729	0	0
22	298	729	0	0
193	296	724	0	0
193	297	726	0	0
193	300	724	0	0
193	303	726	6	0
193	303	727	6	0
193	302	727	6	0
193	294	734	0	0
193	297	732	0	0
193	301	734	0	0
193	303	732	0	0
193	303	731	6	0
193	302	731	6	0
1179	289	728	4	0
1093	292	728	0	0
1094	292	730	0	0
1093	294	727	0	0
1094	294	731	0	0
1101	300	728	6	0
1102	300	730	6	0
1180	301	728	4	0
1181	303	728	4	0
1069	297	759	4	0
1070	295	758	6	0
1070	298	755	2	0
1070	299	762	5	0
1070	302	756	0	0
1070	302	758	5	0
1070	304	758	4	0
1070	307	757	5	0
145	214	742	6	0
145	214	745	6	0
36	203	723	1	0
399	200	743	1	0
24	209	744	0	0
27	210	747	6	0
205	213	748	0	0
70	215	765	0	0
70	203	763	0	0
70	202	761	0	0
70	205	765	0	0
205	194	748	0	0
205	194	757	0	0
205	197	763	0	0
205	204	757	0	0
205	207	763	0	0
213	581	3553	0	0
\.


--
-- Data for Name: game_objects; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.game_objects (id, name, description, command_one, command_two, type, width, height, modelheight) FROM stdin;
72	Wheat	nice ripe looking wheat	walkto	pick	0	1	1	0
618	glider	i wonder if it flys	fly	examine	1	1	1	0
755	Swamp	That smells horrid	walkto	examine	1	1	1	0
0	Tree	A pointy tree	chop	examine	1	1	1	0
1	Tree	A leafy tree	chop	examine	1	1	1	0
2	Well	The bucket is missing	walkto	examine	1	2	2	0
3	Table	A mighty fine table	walkto	examine	1	1	1	96
4	Treestump	Someone has chopped this tree down!	walkto	examine	1	1	1	0
5	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
6	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
7	Chair	A sturdy looking chair	walkto	examine	1	1	1	0
8	logs	A pile of logs	walkto	examine	1	1	1	0
9	Longtable	It has nice candles	walkto	examine	1	4	1	0
10	Throne	It looks fancy and expensive	walkto	examine	1	1	1	0
11	Range	A hot well stoked range	walkto	examine	1	1	2	0
12	Gravestone	R I P	walkto	examine	1	1	1	0
13	Gravestone	Its covered in moss	walkto	examine	1	1	1	0
14	Bed	Ooh nice blankets	rest	examine	1	2	3	0
15	Bed	Its a bed - wow	rest	examine	1	2	2	0
16	bar	Mmm beer	walkto	examine	1	1	1	0
17	Chest	Perhaps I should search it	search	close	1	1	1	0
18	Chest	I wonder what is inside...	open	examine	1	1	1	0
19	Altar	Its an Altar	recharge at	examine	1	2	1	0
20	Post	What am I examining posts for	walkto	examine	1	1	1	0
21	Support	A wooden pole	walkto	examine	0	1	1	0
22	barrel	Its empty	walkto	examine	1	1	1	0
23	Bench	It doesn't look very comfy	walkto	examine	1	2	1	0
24	Portrait	A painting of our beloved king	walkto	examine	0	1	1	0
25	candles	Candles on a fancy candlestick	walkto	examine	1	1	1	0
26	fountain	The water looks fairly clean	walkto	examine	1	2	2	0
27	landscape	An oil painting	walkto	examine	0	1	1	0
28	Millstones	You can use these to make flour	walkto	examine	1	3	3	0
29	Counter	It's the shop counter	walkto	examine	1	2	1	120
30	Stall	A market stall	walkto	examine	1	2	2	112
31	Target	Coming soon archery practice	practice	examine	1	1	1	0
32	PalmTree	A nice palm tree	walkto	examine	1	1	1	0
33	PalmTree	A shady palm tree	walkto	examine	1	1	1	0
34	Fern	A leafy plant	walkto	examine	0	1	1	0
35	Cactus	It looks very spikey	walkto	examine	1	1	1	0
36	Bullrushes	I wonder why it's called a bullrush	walkto	examine	0	1	1	0
37	Flower	Ooh thats pretty	walkto	examine	0	1	1	0
38	Mushroom	I think it's a poisonous one	walkto	examine	0	1	1	0
39	Coffin	This coffin is closed	open	examine	1	2	2	0
40	Coffin	This coffin is open	search	close	1	2	2	0
41	stairs	These lead upstairs	go up	examine	1	2	3	0
42	stairs	These lead downstairs	go down	examine	1	2	3	0
43	stairs	These lead upstairs	go up	examine	1	2	3	0
44	stairs	These lead downstairs	go down	examine	1	2	3	0
45	railing	nice safety measure	walkto	examine	1	1	1	0
46	pillar	An ornate pillar	walkto	examine	1	1	1	0
47	Bookcase	A large collection of books	walkto	examine	1	1	2	0
48	Sink	Its fairly dirty	walkto	examine	1	1	2	0
49	Dummy	I can practice my fighting here	hit	examine	1	1	1	0
50	anvil	heavy metal	walkto	examine	1	1	1	0
51	Torch	It would be very dark without this	walkto	examine	0	1	1	0
52	hopper	You put grain in here	operate	examine	1	2	2	0
53	chute	Flour comes out here	walkto	examine	1	2	2	40
54	cart	A farm cart	walkto	examine	1	2	3	0
55	sacks	Yep they're sacks	walkto	examine	1	1	1	0
56	cupboard	The cupboard is shut	open	examine	1	1	2	0
57	Gate	The gate is closed	open	examine	2	1	2	0
58	gate	The gate is open	walkto	close	3	1	2	0
59	gate	The gate is open	walkto	close	3	1	2	0
60	gate	The gate is closed	open	examine	2	1	2	0
61	signpost	To Varrock	walkto	examine	1	1	1	0
62	signpost	To the tower of wizards	walkto	examine	1	1	1	0
63	doors	The doors are open	walkto	close	3	1	2	0
64	doors	The doors are shut	open	examine	2	1	2	0
65	signpost	To player owned houses	walkto	examine	1	1	1	0
66	signpost	To Lumbridge Castle	walkto	examine	1	1	1	0
67	bookcase	It's a bookcase	walkto	search	1	1	2	0
68	henge	these look impressive	walkto	examine	1	2	2	0
69	Dolmen	A sort of ancient altar thingy	walkto	examine	1	2	2	0
70	Tree	This tree doesn't look too healthy	walkto	chop	1	1	1	0
71	cupboard	Perhaps I should search it	search	close	1	1	2	0
73	sign	The blue moon inn	walkto	examine	0	1	1	0
74	sails	The windmill's sails	walkto	examine	0	1	3	0
75	sign	estate agent	walkto	examine	0	1	1	0
76	sign	The Jolly boar inn	walkto	examine	0	1	1	0
77	Drain	This drainpipe runs from the kitchen to the sewers	walkto	search	0	1	1	0
78	manhole	A manhole cover	open	examine	0	1	1	0
79	manhole	How dangerous - this manhole has been left open	climb down	close	0	1	1	0
80	pipe	a dirty sewer pipe	walkto	examine	1	1	1	0
81	Chest	Perhaps I should search it	search	close	1	1	1	0
82	Chest	I wonder what is inside...	open	examine	1	1	1	0
83	barrel	It seems to be full of newt's eyes	walkto	examine	1	1	1	0
84	cupboard	The cupboard is shut	open	examine	1	1	2	0
85	cupboard	Perhaps I should search it	search	close	1	1	2	0
86	fountain	I think I see something in the fountain	walkto	search	1	2	2	0
87	signpost	To Draynor Manor	walkto	examine	1	1	1	0
88	Tree	This tree doesn't look too healthy	approach	search	1	1	1	0
89	sign	General Store	walkto	examine	0	1	1	0
90	sign	Lowe's Archery store	walkto	examine	0	1	1	0
91	sign	The Clothes Shop	walkto	examine	0	1	1	0
92	sign	Varrock Swords	walkto	examine	0	1	1	0
93	gate	You can pass through this on the members server	open	examine	2	1	2	0
94	gate	You can pass through this on the members server	open	examine	2	1	2	0
95	sign	Bob's axes	walkto	examine	0	1	1	0
96	sign	The staff shop	walkto	examine	0	1	1	0
97	fire	A strongly burning fire	walkto	examine	0	1	1	0
98	Rock	A rocky outcrop	mine	prospect	1	1	1	0
99	Rock	A rocky outcrop	mine	prospect	1	1	1	0
100	Rock	A rocky outcrop	mine	prospect	1	1	1	0
101	Rock	A rocky outcrop	mine	prospect	1	1	1	0
102	Rock	A rocky outcrop	mine	prospect	1	1	1	0
103	Rock	A rocky outcrop	mine	prospect	1	1	1	0
104	Rock	A rocky outcrop	mine	prospect	1	1	1	0
105	Rock	A rocky outcrop	mine	prospect	1	1	1	0
106	Rock	A rocky outcrop	mine	prospect	1	1	1	0
107	Rock	A rocky outcrop	mine	prospect	1	1	1	0
108	Rock	A rocky outcrop	mine	prospect	1	1	1	0
109	Rock	A rocky outcrop	mine	prospect	1	1	1	0
110	Rock	A rocky outcrop	mine	prospect	1	1	1	0
111	Rock	A rocky outcrop	mine	prospect	1	1	1	0
112	Rock	A rocky outcrop	mine	prospect	1	1	1	0
113	Rock	A rocky outcrop	mine	prospect	1	1	1	0
114	Rock	A rocky outcrop	mine	prospect	1	1	1	0
115	Rock	A rocky outcrop	mine	prospect	1	1	1	0
116	web	A spider's web	walkto	examine	0	1	1	0
117	web	A spider's web	walkto	examine	0	1	1	0
118	furnace	A red hot furnace	walkto	examine	1	2	2	0
119	Cook's Range	A hot well stoked range	walkto	examine	1	1	2	0
120	Machine	I wonder what it's supposed to do	walkto	examine	1	2	2	0
121	Spinning wheel	I can spin wool on this	walkto	examine	1	1	1	0
122	Lever	The lever is up	walkto	examine	0	1	1	0
123	Lever	The lever is down	walkto	examine	0	1	1	0
124	LeverA	It's a lever	pull	inspect	0	1	1	0
125	LeverB	It's a lever	pull	inspect	0	1	1	0
126	LeverC	It's a lever	pull	inspect	0	1	1	0
127	LeverD	It's a lever	pull	inspect	0	1	1	0
128	LeverE	It's a lever	pull	inspect	0	1	1	0
129	LeverF	It's a lever	pull	inspect	0	1	1	0
130	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
131	signpost	To the forge	walkto	examine	1	1	1	0
132	signpost	To the Barbarian's  Village	walkto	examine	1	1	1	0
133	signpost	To Al Kharid	walkto	examine	1	1	1	0
134	Compost Heap	A smelly pile of compost	walkto	search	1	2	2	0
135	Coffin	This coffin is closed	open	examine	1	2	2	0
136	Coffin	This coffin is open	search	close	1	2	2	0
137	gate	You can pass through this on the members server	open	examine	2	1	2	0
138	gate	You can pass through this on the members server	open	examine	2	1	2	0
139	sign	The Bank of runescape	walkto	examine	0	1	1	0
140	cupboard	The cupboard is shut	open	examine	1	1	2	0
141	cupboard	Perhaps I should search it	search	close	1	1	2	0
142	doors	The doors are shut	open	examine	2	1	2	0
143	torch	A scary torch	walkto	examine	0	1	1	0
144	Altar	An altar to the evil God Zamorak	recharge at	examine	1	2	1	0
145	Shield	A display shield	walkto	examine	0	1	1	0
146	Grill	some sort of ventilation	walkto	examine	0	1	1	0
147	Cauldron	A very large pot	walkto	drink from	1	1	1	0
148	Grill	some sort of ventilation	listen	examine	0	1	1	0
149	Mine Cart	It's empty	walkto	examine	1	1	1	0
150	Buffers	Stop the carts falling off the end	walkto	examine	1	1	1	0
151	Track	Train track	walkto	examine	0	2	2	0
152	Track	Train track	walkto	examine	0	2	2	0
153	Track	Train track	walkto	examine	0	1	1	0
154	Hole	I can see a witches cauldron directly below it	walkto	examine	1	1	1	0
155	ship	A ship to Karamja	board	examine	0	5	3	0
156	ship	A ship to Karamja	board	examine	0	2	3	0
157	ship	A ship to Karamja	board	examine	0	5	3	0
158	Emergency escape ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
159	sign	Wydin's grocery	walkto	examine	0	1	1	0
160	sign	The Rusty Anchor	walkto	examine	0	1	1	0
161	ship	A ship to Port Sarim	board	examine	0	5	3	0
162	ship	A ship to Port Sarim	board	examine	0	2	3	0
163	ship	A ship to Port Sarim	board	examine	0	5	3	0
164	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
165	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
166	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
167	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
168	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
169	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
170	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
171	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
172	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
173	hopper	You put grain in here	operate	examine	1	2	2	0
174	cupboard	The cupboard is shut	open	examine	1	1	2	0
175	cupboard	Perhaps I should search it	search	close	1	1	2	0
176	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
177	Doric's anvil	Property of Doric the dwarf	walkto	examine	1	1	1	0
178	pottery oven	I can fire clay pots in this	walkto	examine	1	2	2	0
179	potter's wheel	I can make clay pots using this	walkto	examine	1	1	1	0
180	gate	A gate from Lumbridge to Al Kharid	open	examine	2	1	2	0
181	gate	This gate is open	walkto	examine	2	1	2	0
182	crate	A crate used for storing bananas	walkto	search	1	1	1	0
183	Banana tree	A tree with nice ripe bananas growing on it	walkto	pick banana	1	1	1	0
184	Banana tree	There are no bananas left on the tree	walkto	pick banana	1	1	1	0
185	crate	A crate used for storing bananas	walkto	search	1	1	1	0
186	Chest	A battered old chest	walkto	examine	1	1	1	0
187	Chest	I wonder what is inside...	open	examine	1	1	1	0
188	Flower	Ooh thats pretty	walkto	examine	0	1	1	0
189	sign	Fishing Supplies	walkto	examine	0	1	1	0
190	sign	Jewellers	walkto	examine	0	1	1	0
191	Potato	A potato plant	walkto	pick	0	1	1	0
192	fish	I can see fish swimming in the water	lure	bait	0	1	1	0
193	fish	I can see fish swimming in the water	net	bait	0	1	1	0
194	fish	I can see fish swimming in the water	harpoon	cage	0	1	1	0
195	Rock	A rocky outcrop	mine	prospect	1	1	1	0
196	Rock	A rocky outcrop	mine	prospect	1	1	1	0
197	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
198	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
199	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
200	Monks Altar	Its an Altar	recharge at	examine	1	2	1	0
201	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
202	Coffin	This coffin is closed	open	examine	1	2	2	0
203	Coffin	This coffin is open	search	close	1	2	2	0
204	Smashed table	This table has seen better days	walkto	examine	1	1	1	0
205	Fungus	A creepy looking fungus	walkto	examine	0	1	1	0
206	Smashed chair	This chair is broken	walkto	examine	1	1	1	0
207	Broken pillar	The remains of a pillar	walkto	examine	1	1	1	0
208	Fallen tree	A fallen tree	walkto	examine	1	3	2	0
209	Danger Sign	Danger!	walkto	examine	1	1	1	0
210	Rock	A rocky outcrop	mine	prospect	1	1	1	0
211	Rock	A rocky outcrop	mine	prospect	1	1	1	0
212	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
213	Gravestone	A big impressive gravestone	walkto	examine	1	2	2	0
214	bone	Eep!	walkto	examine	1	1	1	0
215	bone	This would feed a dog for a month	walkto	examine	1	1	1	0
216	carcass	I think it's dead	walkto	examine	1	2	2	0
217	animalskull	I wouldn't like to meet a live one	walkto	examine	1	1	1	0
218	Vine	A creepy creeper	walkto	examine	0	1	1	0
219	Vine	A creepy creeper	walkto	examine	0	1	1	0
220	Vine	A creepy creeper	walkto	examine	0	1	1	0
221	Chest	Perhaps I should search it	walkto	examine	1	1	1	0
222	Chest	I wonder what is inside...	open	examine	1	1	1	0
223	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
224	ship	The Lumbridge Lady	board	examine	0	5	3	0
225	ship	The Lumbridge Lady	board	examine	0	5	3	0
226	hole	This ship isn't much use with that there	walkto	examine	2	1	1	0
227	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
228	Chest	Perhaps I should search it	search	close	1	1	1	0
229	Chest	I wonder what is inside...	open	examine	1	1	1	0
230	Chest	Perhaps I should search it	search	close	1	1	1	0
231	Chest	I wonder what is inside...	open	examine	1	1	1	0
232	hole	This ship isn't much use with that there	walkto	examine	2	1	1	0
233	ship	The Lumbridge Lady	board	examine	0	5	3	0
234	ship	The Lumbridge Lady	board	examine	0	5	3	0
235	Altar of Guthix	A sort of ancient altar thingy	recharge at	examine	1	2	2	0
236	The Cauldron of Thunder	A very large pot	walkto	examine	1	1	1	0
237	Tree	A leafy tree	search	examine	1	1	1	0
238	ship	A ship to Entrana	board	examine	0	5	3	0
239	ship	A ship to Entrana	board	examine	0	2	3	0
240	ship	A ship to Entrana	board	examine	0	5	3	0
241	ship	A ship to Port Sarim	board	examine	0	5	3	0
242	ship	A ship to Port Sarim	board	examine	0	2	3	0
243	ship	A ship to Port Sarim	board	examine	0	5	3	0
244	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
245	Dramen Tree	This tree doesn't look too healthy	chop	examine	1	1	1	0
246	hopper	You put grain in here	operate	examine	1	2	2	0
247	Chest	Perhaps I should search it	walkto	examine	1	1	1	0
248	Chest	I wonder what is inside...	open	examine	1	1	1	0
249	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
250	sign	2-handed swords sold here	walkto	examine	0	1	1	0
251	sign	ye olde herbalist	walkto	examine	0	1	1	0
252	gate	You can pass through this on the members server	open	examine	2	1	2	0
253	gate	You can pass through this on the members server	open	examine	2	1	2	0
254	gate	You can pass through this on the members server	open	examine	2	1	2	0
255	Door mat	If I ever get my boots muddy I know where to come	search	examine	0	1	1	0
256	gate	The gate is closed	open	examine	2	1	2	0
257	Cauldron	A very large pot	walkto	examine	1	1	1	0
258	cupboard	The cupboard is shut	open	examine	1	1	2	0
259	cupboard	Perhaps I should search it	search	close	1	1	2	0
260	gate	The bank vault gate	open	examine	2	1	2	0
261	fish	I can see fish swimming in the water	net	harpoon	0	1	1	0
262	sign	Harry's fishing shack	walkto	examine	0	1	1	0
263	cupboard	The cupboard is shut	open	examine	1	1	2	0
264	cupboard	Perhaps I should search it	search	close	1	1	2	0
265	Chest	Perhaps I should search it	search	close	1	1	1	0
266	Chest	I wonder what is inside...	open	examine	1	1	1	0
267	sign	The shrimp and parrot	walkto	examine	0	1	1	0
268	signpost	Palm Street	walkto	examine	1	1	1	0
269	Rockslide	A pile of rocks blocks your path	mine	prospect	1	1	1	0
270	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
271	fish	I can see fish swimming in the lava!	bait	examine	0	1	1	0
272	barrel	Its got ale in it	walkto	examine	1	1	1	0
273	table	It's a sturdy table	walkto	examine	1	2	1	96
274	Fireplace	It would be very cold without this	walkto	examine	1	2	1	0
275	Egg	Thats one big egg!	walkto	examine	1	1	1	0
276	Eggs	They'd make an impressive omlette	walkto	examine	1	1	1	0
277	Stalagmites	Hmm pointy	walkto	examine	1	1	1	0
278	Stool	A simple three legged stool	walkto	examine	1	1	1	0
279	Bench	It doesn't look to comfortable	walkto	examine	1	1	1	0
280	table	A round table ideal for knights	walkto	examine	1	2	2	0
281	table	A handy little table	walkto	examine	1	1	1	96
282	fountain of heros	Use a dragonstone gem here to increase it's abilties	walkto	examine	1	2	2	0
283	bush	A leafy bush	walkto	examine	1	1	1	0
284	hedge	A carefully trimmed hedge	walkto	examine	1	1	1	0
285	flower	A nice colourful flower	walkto	examine	1	1	1	0
286	plant	Hmm leafy	walkto	examine	1	1	1	0
287	Giant crystal	How unusual a crystal with a wizard trapped in it	walkto	examine	1	3	3	0
288	sign	The dead man's chest	walkto	examine	0	1	1	0
289	sign	The rising sun	walkto	examine	0	1	1	0
290	crate	A large wooden storage box	walkto	search	1	1	1	0
291	crate	A large wooden storage box	walkto	search	1	1	1	0
292	ship	A merchant ship	stow away	examine	0	5	3	0
293	ship	A merchant ship	stow away	examine	0	5	3	0
294	beehive	It's guarded by angry looking bees	walkto	examine	1	1	1	0
295	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
296	Altar	An altar to the evil God Zamorak	recharge at	search	1	2	1	0
297	sign	Hickton's Archery store	walkto	examine	0	1	1	0
298	signpost	To Camelot	walkto	examine	1	1	1	0
299	Archway	A decorative marble arch	walkto	examine	1	2	1	0
300	Obelisk of water	It doesn't look very wet	walkto	examine	1	1	1	0
301	Obelisk of fire	It doesn't look very hot	walkto	examine	1	1	1	0
302	sand pit	I can use a bucket to get sand from here	walkto	search	1	2	2	0
303	Obelisk of air	A tall stone pointy thing	walkto	examine	1	1	1	0
304	Obelisk of earth	A tall stone pointy thing	walkto	examine	1	1	1	0
305	gate	You can pass through this on the members server	open	examine	2	1	2	0
306	Oak Tree	A grand old oak tree	chop	examine	1	2	2	0
307	Willow Tree	A weeping willow	chop	examine	1	2	2	0
308	Maple Tree	It's got nice shaped leaves	chop	examine	1	2	2	0
309	Yew Tree	A tough looking yew tree	chop	examine	1	2	2	0
310	Tree	A magical tree	chop	examine	1	1	1	0
311	gate	A gate guarded by a fierce barbarian	open	examine	2	1	2	0
312	sign	The forester's arms	walkto	examine	0	1	1	0
313	flax	A flax plant	walkto	pick	0	1	1	0
314	Large treestump	Someone has chopped this tree down!	walkto	examine	1	2	2	0
315	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
316	Lever	It's a lever	pull	inspect	0	1	1	0
317	Lever	It's a lever	pull	inspect	0	1	1	0
318	Lever	It's a lever	pull	inspect	0	1	1	0
319	gate	You can pass through this on the members server	open	examine	2	1	2	0
320	ship	A ship bound for Ardougne	board	examine	0	5	3	0
321	ship	A ship bound for Ardougne	board	examine	0	5	3	0
322	Bakers Stall	A market stall	walkto	steal from	1	2	2	112
323	Silk Stall	A market stall	walkto	steal from	1	2	2	112
324	Fur Stall	A market stall	walkto	steal from	1	2	2	112
325	Silver Stall	A market stall	walkto	steal from	1	2	2	112
326	Spices Stall	A market stall	walkto	steal from	1	2	2	112
327	gems Stall	A market stall	walkto	steal from	1	2	2	112
328	crate	A large heavy sealed crate	walkto	search	1	1	1	0
329	crate	A large heavy sealed crate	walkto	search	1	1	1	0
330	sign	RPDT depot	walkto	examine	0	1	1	0
331	stairs	These lead upstairs	go up	search for traps	1	2	3	0
332	Chest	Perhaps I should search it	search	close	1	1	1	0
333	Chest	I wonder what is inside...	open	examine	1	1	1	0
334	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
335	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
336	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
337	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
338	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
339	Chest	someone is stealing something from it	walkto	examine	1	1	1	0
340	Chest	I wonder what is inside...	open	search for traps	1	1	1	0
341	empty stall	A market stall	walkto	examine	1	2	2	112
342	stairs	These lead upstairs	go up	examine	1	2	3	0
343	hopper	You put grain in here	operate	examine	1	2	2	0
344	signpost	Ardougne city zoo	walkto	examine	1	1	1	0
345	sign	The flying horse	walkto	examine	0	1	1	0
346	gate	You can pass through this on the members server	open	examine	2	1	2	0
347	gate	You can pass through this on the members server	open	examine	2	1	2	0
348	Lever	The lever is up	pull	examine	0	1	1	0
349	Lever	The lever is up	pull	examine	0	1	1	0
350	pipe	a dirty sewer pipe	walkto	examine	1	1	1	0
351	fish	I can see fish swimming in the water	bait	examine	0	1	1	0
352	fish	I can see fish swimming in the water	bait	examine	0	1	1	0
353	fish	I can see fish swimming in the water	bait	examine	0	1	1	0
354	fish	I can see fish swimming in the water	bait	examine	0	1	1	0
355	Vine	A creepy creeper	walkto	examine	0	1	1	0
356	gate	The main entrance to McGrubor's wood	open	examine	2	1	2	0
357	gate	The gate is open	walkto	examine	2	1	2	0
358	gate	The gate is closed	open	examine	2	1	2	0
359	stairs	These lead downstairs	go down	examine	1	2	3	0
360	broken cart	A farm cart	walkto	examine	1	2	3	0
361	Lever	It's a lever	pull	searchfortraps	0	1	1	0
362	clock pole blue	A pole - a pole to put cog's on	inspect	examine	0	1	1	0
363	clock pole red	A pole - a pole to put cog's on	inspect	examine	0	1	1	0
364	clock pole purple	A pole - a pole to put cog's on	inspect	examine	0	1	1	0
365	clock pole black	A pole - a pole to put cog's on	inspect	examine	0	1	1	0
366	wallclockface	It's a large clock face	walkto	examine	1	2	2	0
367	Lever Bracket	Theres something missing here	walkto	examine	0	1	1	0
368	Lever	It's a lever	pull	examine	0	1	1	0
369	stairs	These lead upstairs	go up	examine	1	2	3	0
370	stairs	These lead downstairs	go down	examine	1	2	3	0
371	gate	The gate is closed	open	examine	2	1	2	0
372	gate	The gate is open	close	examine	3	1	2	0
373	Lever	The lever is up	pull	examine	0	1	1	0
374	Lever	The lever is up	push	examine	0	1	1	0
375	Foodtrough	It's for feeding the rat's	walkto	examine	1	2	1	0
376	fish	I can see fish swimming in the water	cage	harpoon	0	1	1	0
377	spearwall	It's a defensive battlement	walkto	examine	1	2	1	0
378	hornedskull	A horned dragon skull	walkto	examine	1	2	2	0
379	Chest	I wonder what is inside...	open	picklock	1	1	1	0
380	Chest	I wonder what is inside...	open	picklock	1	1	1	0
381	guardscupboard	The cupboard is shut	open	examine	1	1	2	0
382	guardscupboard	Perhaps I should search it	search	close	1	1	2	0
383	Coal truck	I can use this to transport coal	get coal from	examine	1	1	1	0
384	ship	A ship to Port Birmhaven	board	examine	0	5	3	0
385	ship	A ship to Port Birmhaven	board	examine	0	2	3	0
386	ship	A ship to Port Birmhaven	board	examine	0	5	3	0
387	Tree	It's a tree house	walkto	examine	1	1	1	0
388	Ballista	It's a war machine	fire	examine	1	4	1	0
389	largespear		walkto	examine	1	2	1	0
390	spirit tree	A grand old spirit tree	talk to	examine	1	2	2	0
391	young spirit Tree	Ancestor of the spirit tree	talk to	examine	1	1	1	0
392	gate	The gate is closed	talk through	examine	2	1	2	0
393	wall	A damaged wall	climb	examine	1	3	1	0
394	tree	An exotic looking tree	walkto	examine	1	1	1	0
395	tree	An exotic looking tree	walkto	examine	1	1	1	0
396	Fern	An exotic leafy plant	walkto	examine	0	1	1	0
397	Fern	An exotic leafy plant	walkto	examine	0	1	1	0
398	Fern	An exotic leafy plant	walkto	examine	0	1	1	0
399	Fern	An exotic leafy plant	walkto	examine	0	1	1	0
400	fly trap	A small carnivourous plant	approach	search	0	1	1	0
401	Fern	An exotic leafy plant	walkto	examine	0	1	1	0
402	Fern	An exotic spikey plant	walkto	examine	0	1	1	0
403	plant	What an unusual plant	walkto	examine	0	1	1	0
404	plant	An odd looking plant	walkto	examine	1	1	1	0
405	plant	some nice jungle foliage	walkto	examine	1	1	1	0
406	stone head	It looks like it's been here some time	walkto	examine	1	2	2	0
407	dead Tree	A rotting tree	walkto	examine	1	1	1	0
408	sacks	Yep they're sacks	walkto	prod	1	1	1	0
409	khazard open Chest	Perhaps I should search it	search	close	1	1	1	0
410	khazard shut Chest	I wonder what is inside...	open	examine	1	1	1	0
411	doorframe	It's a stone doorframe	walkto	examine	3	1	2	0
412	Sewer valve	It changes the water flow of the sewer's	turn left	turn right	1	1	1	0
413	Sewer valve 2	It changes the water flow of the sewer's	turn left	turn right	1	1	1	0
414	Sewer valve 3	It changes the water flow of the sewer's	turn left	turn right	1	1	1	0
415	Sewer valve 4	It changes the water flow of the sewer's	turn left	turn right	1	1	1	0
416	Sewer valve 5	It changes the water flow of the sewer's	turn left	turn right	1	1	1	0
417	Cave entrance	I wonder what is inside...	enter	examine	1	2	2	0
418	Log bridge	A tree gnome construction	walkto	examine	0	1	1	0
419	Log bridge	A tree gnome construction	walkto	examine	0	1	1	0
420	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
421	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
422	gate	The gate is open	close	examine	2	1	2	0
423	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
424	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
425	Log bridge	A tree gnome construction	walkto	examine	0	1	1	0
426	Log bridge	A tree gnome construction	walkto	examine	0	1	1	0
427	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
428	tree platform	A tree gnome construction	walkto	examine	0	1	1	0
429	Tribal brew	A very large pot	walkto	drink	1	1	1	0
430	Pineapple tree	A tree with nice ripe pineapples growing on it	walkto	pick pineapple	1	1	1	0
431	Pineapple tree	There are no pineapples left on the tree	walkto	pick pineapple	1	1	1	0
432	log raft	A mighty fine raft	board	examine	0	1	1	96
433	log raft	A mighty fine raft	board	examine	0	1	1	96
434	Tomb of hazeel	A clay shrine to lord hazeel	walkto	examine	1	1	2	96
435	range	A pot of soup slowly cooking	walkto	examine	1	1	2	0
436	Bookcase	A large collection of books	search	examine	1	1	2	0
437	Carnillean Chest	Perhaps I should search it	walkto	close	1	1	1	0
438	Carnillean Chest	I wonder what is inside...	open	examine	1	1	1	0
439	crate	A crate used for storing food	search	examine	1	1	1	0
440	Butlers cupboard	The cupboard is shut	open	examine	1	1	2	0
441	Butlers cupboard	The cupboard is open	search	close	1	1	2	0
442	gate	The gate is open	walkto	examine	2	1	2	0
443	gate	The gate is closed	open	examine	2	1	2	0
444	Cattle furnace	A red hot furnace	walkto	examine	1	2	2	0
445	Ardounge wall	A huge wall seperating east and west ardounge	walkto	examine	1	1	3	0
446	Ardounge wall corner	A huge wall seperating east and west ardounge	walkto	examine	1	1	1	0
447	Dug up soil	A freshly dug pile of mud	walkto	examine	0	1	1	0
448	Pile of mud	Mud caved in from above	climb	examine	1	2	1	0
449	large Sewer pipe	a dirty sewer pipe	enter	examine	1	1	1	0
450	Ardounge wall gateway	A huge set of heavy wooden doors	open	examine	2	1	2	0
451	cupboard	The cupboard is shut	open	examine	1	1	2	0
452	cupboard	The cupboard is open	search	close	1	1	2	0
453	Fishing crane	For hauling in large catches of fish	operate	examine	1	1	2	0
454	Rowboat	A reasonably sea worthy two man boat	walkto	examine	1	2	2	0
455	Damaged Rowboat	A not so sea worthy two man boat	walkto	examine	1	2	2	0
456	barrel	I wonder what's inside	walkto	search	1	1	1	0
457	gate	The gate is closed	open	examine	2	1	2	0
458	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
459	Fishing crane	For hauling in large catches of fish	operate	examine	1	1	1	0
460	Fishing crane	For hauling in large catches of fish	operate	examine	1	1	1	0
461	Waterfall	it's a waterfall	walkto	examine	2	1	2	0
462	leaflessTree	A pointy tree	jump off	jump to next	1	1	1	0
463	leaflessTree	A pointy tree	jump off	jump to next	1	1	1	0
464	log raft	A mighty fine raft	board	examine	0	1	1	96
465	doors	The doors are shut	open	examine	2	1	2	0
466	Well	An oddly placed well	operate	examine	1	2	2	0
467	Tomb of glarial	A stone tomb surrounded by flowers	search	examine	1	2	4	96
468	Waterfall	it's a fast flowing waterfall	jump off	examine	2	1	2	0
469	Waterfall	it's a fast flowing waterfall	jump off	examine	0	1	2	0
470	Bookcase	A large collection of books	search	examine	1	1	2	0
471	doors	The doors are shut	open	examine	2	1	2	0
472	doors	The doors are shut	open	examine	2	1	2	0
473	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
474	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
475	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
476	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
477	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
478	Stone stand	On top is an indent the size of a rune stone	walkto	examine	1	1	1	0
479	Glarial's Gravestone	There is an indent the size of a pebble in the stone's center	read	examine	1	1	1	0
480	gate	The gate is closed	open	examine	2	1	2	0
481	crate	It's a crate	search	examine	1	1	1	0
482	leaflessTree	A pointy tree	jump off	examine	1	1	1	0
483	Statue of glarial	A statue of queen glarial - something's missing	walkto	examine	1	1	1	0
484	Chalice of eternity	A magically elevated chalice full of treasure	walkto	examine	1	1	1	0
485	Chalice of eternity	A magically elevated chalice full of treasure	empty	examine	1	1	1	0
486	doors	The doors are shut	open	examine	2	1	2	0
487	Lever	The lever is up	pull	examine	0	1	1	0
488	Lever	The lever is up	pull	examine	0	1	1	0
489	log raft remains	oops!	walkto	examine	0	1	1	96
490	Tree	A pointy tree	walkto	examine	0	1	1	0
491	 Range	A hot well stoked range	walkto	examine	1	1	2	0
492	crate	It's an old crate	walkto	search	1	1	1	0
493	fish	I can see fish swimming in the water	net	examine	0	1	1	0
494	Watch tower	They're always watching	approach	examine	1	2	2	0
495	signpost	Tourist infomation	walkto	examine	1	1	1	0
496	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
497	doors	The doors are open	walkto	examine	2	1	2	0
498	Rope ladder	A hand made ladder	walkto	examine	1	1	1	0
499	cupboard	The cupboard is shut	open	examine	1	1	2	0
500	cupboard	Perhaps I should search it	search	close	1	1	2	0
501	Rope ladder	A hand made ladder	walkto	examine	1	1	1	0
502	Cooking pot	the mourners are busy enjoying this stew	walkto	examine	1	1	1	0
503	Gallow	Best not hang about!	walkto	examine	1	2	2	0
504	gate	The gate is closed	open	examine	2	1	2	0
505	crate	A crate used for storing confiscated goods	walkto	search	1	1	1	0
506	cupboard	The cupboard is shut	open	examine	1	1	2	0
507	cupboard	Perhaps I should search it	search	close	1	1	2	0
508	gate	You can pass through this on the members server	open	examine	2	1	2	0
509	cupboard	The cupboard is shut	open	examine	1	1	2	0
510	cupboard	Perhaps I should search it	search	close	1	1	2	0
511	sign	Tailors fancy dress	walkto	examine	0	1	1	0
512	grand tree	the grand tree	walkto	examine	0	1	1	0
513	gate	The gate is closed	open	examine	2	1	2	0
514	gate	The gate is open	walkto	close	2	1	2	0
515	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
516	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
517	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
518	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
519	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
520	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
521	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
522	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
523	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
524	Log bridge	A tree gnome construction	walkto	examine	0	1	1	0
525	Watch tower	A tree gnome construction	walkto	examine	1	1	1	0
526	Log bridge	A tree gnome construction	walkto	examine	0	0	0	0
527	climbing rocks	I wonder if I can climb up these	climb	examine	0	1	1	0
528	Ledge	It looks rather thin	balance on	examine	0	1	0	0
529	Ledge	It looks rather thin	balance on	examine	0	1	1	0
530	log	It looks slippery	walkto	examine	0	1	1	0
531	log	It looks slippery	walkto	examine	0	1	1	0
532	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
533	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
534	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
535	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
536	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
537	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
538	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
539	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
540	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
541	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
542	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
543	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
544	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
545	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
546	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
547	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
548	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
549	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
550	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
551	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
552	Rocks	A rocky outcrop	mine	prospect	1	1	1	0
553	PalmTree	A shady palm tree	walkto	search	1	1	1	0
554	Scorched Earth	An area of burnt soil	walkto	search	0	1	1	0
555	Rocks	A moss covered rock	mine	search	1	1	1	0
556	sign	The dancing donkey inn	walkto	examine	0	1	1	0
557	fish	I can see fish swimming in the water	harpoon	cage	0	1	1	0
558	Rocky Walkway	A precarious rocky walkway	balance	examine	1	1	1	0
559	Rocky Walkway	A precarious rocky walkway	balance	examine	1	1	1	0
560	Rocky Walkway	A precarious rocky walkway	balance	examine	1	1	1	0
561	Rocky Walkway	A precarious rocky walkway	balance	examine	1	1	1	0
562	fight Dummy	I can practice my fighting here	hit	examine	1	1	1	0
563	gate	The gate is closed	open	examine	2	1	2	0
564	Jungle Vine	A deep jungle Vine	walkto	search	0	1	1	0
565	statue	hand carved	walkto	examine	1	1	1	0
566	sign	Ye Olde Dragon Inn	walkto	examine	0	1	1	0
567	grand tree	the grand tree	walkto	examine	0	1	1	0
568	grand tree	the grand tree	walkto	examine	0	1	1	0
569	grand tree	the grand tree	walkto	examine	0	1	1	0
570	grand tree	the grand tree	walkto	examine	0	1	1	0
571	grand tree	the grand tree	walkto	examine	0	1	1	0
572	Hillside Entrance	Large doors that seem to lead into the hillside	open	search	2	1	2	0
573	tree	A large exotic looking tree	walkto	search	1	1	1	0
574	Log bridge	A tree gnome construction	walkto	examine	0	0	0	0
575	Tree platform	A tree gnome construction	walkto	examine	0	0	0	0
576	Tree platform	A tree gnome construction	walkto	examine	0	0	0	0
577	Metalic Dungeon Gate	It seems to be closed	open	search	2	1	2	0
578	Log bridge	A tree gnome construction	walkto	examine	0	0	0	0
579	Log bridge	A tree gnome construction	walkto	examine	0	0	0	0
580	Watch tower	A tree gnome construction	walkto	examine	1	0	0	0
581	Watch tower	A tree gnome construction	walkto	examine	1	0	0	0
582	Shallow water	A small opening in the ground with some spots of water	walkto	investigate	1	2	2	0
583	Doors	Perhaps you should give them a push	open	search	2	1	2	0
584	grand tree	the grand tree	walkto	examine	0	1	1	0
585	Tree Ladder	it's a ladder leading upwards	climb-up	examine	0	1	1	0
586	Tree Ladder	it's a ladder leading downwards	climb-down	examine	0	1	1	0
587	blurberrys cocktail bar	the gnome social hot spot	walkto	examine	0	1	1	0
588	Gem Rocks	A rocky outcrop with a vein of semi precious stones	mine	prospect	1	1	1	0
589	Giannes place	Eat green eat gnome cruisine	walkto	examine	0	1	1	0
590	ropeswing	A good place to train agility	walkto	examine	0	1	1	0
591	net	A good place to train agility	walkto	examine	1	2	1	0
592	Frame	A good place to train agility	walkto	examine	1	1	1	0
593	Tree	It has a branch ideal for tying ropes to	walkto	examine	1	2	2	0
594	Tree	I wonder who put that rope there	walkto	examine	1	2	2	0
595	Tree	they look fun to swing on	walkto	examine	1	2	2	0
596	cart	A farm cart	walkto	search	1	2	3	0
597	fence	it doesn't look too strong	walkto	examine	1	1	1	0
598	beam	A plank of wood	walkto	examine	0	1	1	0
599	Sign	read me	walkto	examine	1	1	1	0
600	Sign	Blurberry's cocktail bar	walkto	examine	1	1	1	0
601	Sign	Giannes tree gnome cuisine	walkto	examine	1	1	1	0
602	Sign	Heckel funch's grocery store	walkto	examine	1	1	1	0
603	Sign	Hudo glenfad's grocery store	walkto	examine	1	1	1	0
604	Sign	Rometti's fashion outlet	walkto	examine	1	1	1	0
605	Sign	Tree gnome bank and rometti's fashion outlet	walkto	examine	1	1	1	0
606	Sign	Tree gnome local swamp	walkto	examine	1	1	1	0
607	Sign	Agility training course	walkto	examine	1	1	1	0
608	Sign	To the grand tree	walkto	examine	1	1	1	0
609	Root	To the grand tree	search	examine	1	1	1	0
610	Root	To the grand tree	search	examine	1	1	1	0
611	Metal Gate	The gate is closed	open	examine	2	1	2	0
612	Metal Gate	The gate is open	walkto	close	2	1	2	0
613	A farm cart	It is blocking the entrance to the village	examine	search	1	2	3	0
614	Ledge	It looks rather thin	balance on	examine	0	1	1	0
615	Ledge	It looks rather thin	balance on	examine	0	1	1	0
616	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
617	cage	i don't like the look of that	open	examine	0	1	1	0
619	cupboard	The cupboard is shut	open	examine	1	1	2	0
620	cupboard	Perhaps I should search it	search	close	1	1	2	0
621	stairs	These lead upstairs	go up	examine	1	2	3	0
622	glider	i wonder if it flys	walkto	examine	1	1	1	0
623	gate	The gate is open	walkto	close	1	1	2	0
624	gate	The gate is closed	open	examine	1	1	2	0
625	chaos altar	An altar to the evil God Zamorak	recharge at	examine	1	2	1	0
626	Gnome stronghold gate	The gate is closed	open	examine	2	1	2	0
627	ropeswing	A good place to train agility	swing	examine	0	1	1	0
628	ropeswing	A good place to train agility	swing	examine	0	1	1	0
629	stairs	These lead upstairs	go up	examine	1	2	3	0
630	stairs	These lead downstairs	go down	examine	1	2	3	0
631	Chest	Perhaps I should search it	search	close	1	1	1	0
632	Chest	I wonder what is inside...	open	examine	1	1	1	0
633	Pile of rubble	What a mess	climb	examine	1	2	1	0
634	Stone stand	On top our four indents from left to right	walkto	push down	1	1	1	0
635	Watch tower	A tree gnome construction	climb up	examine	1	1	1	0
636	Pile of rubble	What a mess	climb	examine	1	2	1	0
637	Root	To the grand tree	search	examine	1	1	1	0
638	Root	To the grand tree	push	examine	1	1	1	0
639	Root	To the grand tree	push	examine	1	1	1	0
640	Sign	Home to the Head tree guardian	walkto	examine	1	1	1	0
641	Hammock	They've got to sleep somewhere	lie in	examine	1	1	2	0
642	Goal	You're supposed to throw the ball here	walkto	examine	0	1	1	0
643	stone tile	It looks as if it might move	twist	examine	1	1	1	0
644	Chest	You get a sense of dread from the chest	walkto	examine	1	1	1	0
645	Chest	You get a sense of dread from the chest	open	examine	1	1	1	0
646	Watch tower	A tree gnome construction	walkto	climb down	0	1	1	0
647	net	A good place to train agility	climb	examine	1	2	1	0
648	Watch tower	A tree gnome construction	climb up	examine	1	1	1	0
649	Watch tower	A tree gnome construction	climb down	examine	1	1	1	0
650	ropeswing	A good place to train agility	grab hold of	examine	0	1	1	0
651	Bumpy Dirt	Some disturbed earth	look	search	0	1	1	0
652	pipe	a dirty sewer pipe	walkto	examine	1	1	1	0
653	net	A good place to train agility	climb	examine	1	2	1	0
654	pipe	a dirty sewer pipe	enter	examine	1	1	1	0
655	log	It looks slippery	balance on	examine	0	1	1	0
656	pipe	a dirty sewer pipe	enter	examine	1	1	1	0
657	pipe	a dirty sewer pipe	enter	examine	0	1	1	0
658	Handholds	I wonder if I can climb up these	climb	examine	0	1	1	0
659	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
660	gate	The gate is closed	open	examine	2	1	2	0
661	stronghold spirit Tree	Ancestor of the spirit tree	talk to	examine	1	1	1	0
662	Tree	It has a branch ideal for tying ropes to	walkto	examine	1	2	2	0
663	Tree	I wonder who put that rope there	swing on	examine	1	2	2	0
664	Tree	I wonder who put that rope there	swing on	examine	1	2	2	0
665	Spiked pit	I don't want to go down there	walkto	examine	1	2	2	0
666	Spiked pit	I don't want to go down there	walkto	examine	1	2	2	0
667	Cave	I wonder what is inside...	enter	examine	1	2	2	0
668	stone pebble	Looks like a stone	walkto	examine	1	1	1	0
669	Pile of rubble	Rocks that have caved in	walkto	examine	1	2	1	0
670	Pile of rubble	Rocks that have caved in	walkto	search	1	2	1	0
671	pipe	I might be able to fit through this	enter	examine	1	1	1	0
672	pipe	2	enter	examine	1	1	1	0
673	Stone	Looks like a stone	walkto	examine	1	1	1	0
674	Stone	Looks like a stone	look closer	investigate	1	1	1	0
675	ropeswing	A good place to train agility	swing	examine	0	1	1	0
676	log	It looks slippery	balance on	examine	0	1	1	0
677	net	A good place to train agility	climb up	examine	1	2	1	0
678	Ledge	It looks rather thin	balance on	examine	0	1	1	0
679	Handholds	I wonder if I can climb up these	climb	examine	0	1	1	0
680	log	It looks slippery	balance on	examine	0	1	1	0
681	log	It looks slippery	balance on	examine	0	1	1	0
682	Rotten Gallows	A human corpse hangs from the noose	look	search	1	2	2	0
683	Pile of rubble	Rocks that have caved in	walkto	search	1	2	1	0
684	ropeswing	I wonder what's over here	swing	examine	0	1	1	0
685	ropeswing	I wonder what's over here	swing	examine	0	1	1	0
686	ocks	A moss covered rock	balance	examine	1	1	1	0
687	Tree	This tree doesn't look too healthy	walkto	balance	1	1	1	0
688	Well stacked rocks	Rocks that have been stacked at regular intervals	investigate	search	1	1	1	0
689	Tomb Dolmen	An ancient construct for displaying the bones of the deceased	look	search	1	2	2	0
690	Handholds	I wonder if I can climb up these	climb	examine	0	1	1	0
691	Bridge Blockade	A crudely constructed fence to stop you going further	investigate	jump	1	1	1	0
692	Log Bridge	A slippery log that is a make-do bridge	balance on	examine	0	1	1	0
693	Handholds	I wonder if I can climb up these	climb	examine	0	1	1	0
694	Tree	they look fun to swing on	swing on	examine	1	2	2	0
695	Tree	they look fun to swing on	swing on	examine	1	2	2	0
696	Wet rocks	A rocky outcrop	look	search	1	1	1	0
697	Smashed table	This table has seen better days	examine	craft	1	1	1	0
698	Crude Raft	A crudely constructed raft	disembark	examine	0	1	1	96
699	Daconia rock	Piles of daconia rock	mine	prospect	1	1	1	0
700	statue	A statue to mark Taie Bwo Wannai sacred grounds	walkto	examine	1	1	1	0
701	Stepping stones	A rocky outcrop	balance	jump onto	1	1	1	0
702	gate	The gate is closed	open	examine	2	1	2	0
703	gate	Enter to balance into an agility area	open	examine	2	1	2	0
704	gate	Enter to balance into an agility area	open	examine	2	1	2	0
705	pipe	It looks a tight squeeze	enter	examine	1	1	1	0
706	ropeswing	A good place to train agility	swing	examine	0	1	1	0
707	Stone	Looks like a stone	balance on	examine	1	1	1	0
708	Ledge	It doesn't look stable	balance on	examine	0	1	1	0
709	Vine	A creepy creeper	climb up	examine	0	1	1	0
710	Rocks	A rocky outcrop	walkto	climb	1	1	1	0
711	Wooden Gate	The gate is open	close	close	2	1	2	0
712	Wooden Gate	The gate is closed	open	examine	2	1	2	0
713	Stone bridge	An ancient stone construction	walkto	examine	0	1	1	0
714	Stone bridge	An ancient stone construction	walkto	examine	0	1	1	0
715	Stone bridge	An ancient stone construction	walkto	examine	0	1	1	0
716	Stone bridge	An ancient stone construction	walkto	examine	0	1	1	0
717	Stone platform	An ancient stone construction	walkto	examine	0	1	1	0
718	fence	it doesn't look too strong	walkto	examine	1	1	1	0
719	Rocks	A rocky outcrop	climb	climb	1	1	1	0
720	Stone bridge	The bridge has partly collapsed	cross	examine	1	3	1	0
721	Stone bridge	The bridge has partly collapsed	cross	examine	1	3	1	0
722	Gate of Iban	It doesn't look very inviting	open	examine	1	3	1	0
723	Wooden Door	It doesn't look very inviting	cross	examine	1	3	1	0
724	Tomb Dolmen	An ancient construct for displaying the bones of the deceased	look	search	1	2	2	0
725	Cave entrance	It doesn't look very inviting	enter	examine	1	3	1	0
726	Old bridge	That's been there a while	walkto	examine	1	3	1	0
727	Old bridge	That's been there a while	cross	examine	1	3	1	0
728	Crumbled rock	climb up to above ground	climb	examine	1	2	1	0
729	stalagmite	Formed over thousands of years	walkto	examine	1	1	1	0
730	stalagmite	Formed over thousands of years	walkto	examine	1	1	1	0
731	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
732	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
733	Lever	Seems to be some sort of winch	pull	examine	1	1	1	0
734	stalactite	Formed over thousands of years	walkto	examine	0	1	1	0
735	stalactite	Formed over thousands of years	walkto	examine	0	1	1	0
736	stalactite	Formed over thousands of years	climb	examine	0	1	1	0
737	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
738	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
739	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
740	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
741	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
742	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
743	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
744	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
745	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
746	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
747	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
748	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
749	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
750	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
751	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
752	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
753	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
754	Swamp	That smells horrid	step over	examine	1	1	1	0
756	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
757	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
758	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
759	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
760	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
761	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
762	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
763	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
764	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
765	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
766	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
767	Pile of mud	Mud caved in from above	climb	examine	1	2	1	0
768	Travel Cart	A sturdy cart for travelling in	board	look	1	2	3	0
769	Travel Cart	A sturdy cart for travelling in	board	look	1	2	3	0
770	Rocks	A small rocky outcrop	mine	examine	1	1	1	0
771	stalactite	Formed over thousands of years	walkto	examine	0	1	1	0
772	Rocks	You should be able to move these	clear	examine	1	2	2	0
773	Rocks	You should be able to move these	walkto	examine	1	2	2	0
774	Rocks	You should be able to move these	walkto	examine	1	1	3	0
775	sign	The Paramaya Hostel	walkto	examine	0	1	1	0
776	Ladder	A ladder that leads to the dormitory - a ticket is needed	climb-up	examine	1	1	1	0
777	Grill	They looks suspicious	walk over	examine	1	1	1	0
778	Spiked pit	They looks suspicious	walk over	examine	1	1	1	0
779	signpost	To the Furnace	walkto	examine	1	1	1	0
780	Ship	A sea faring ship called 'Lady Of The Waves'	board	examine	0	5	3	0
781	Ship	A sea faring ship called 'Lady Of The Waves'	board	examine	0	5	3	0
782	Grill	They looks suspicious	walk over	examine	1	1	1	0
783	sacks	Yep they're sacks	walkto	search	1	1	1	0
784	Zamorakian Temple	Scary!	walkto	examine	0	1	1	0
785	Grill	They looks suspicious	walk over	examine	1	1	1	0
786	Grill	They looks suspicious	walk over	examine	1	1	1	0
787	Grill	They looks suspicious	walk over	examine	1	1	1	0
788	Grill	They looks suspicious	walk over	examine	1	1	1	0
789	Grill	They looks suspicious	walk over	examine	1	1	1	0
790	Grill	They looks suspicious	walk over	examine	1	1	1	0
791	Grill	They looks suspicious	walk over	examine	1	1	1	0
792	Rocks	A small rocky outcrop	walk here	examine	1	1	1	0
793	Rocks	A small rocky outcrop	walk here	examine	1	1	1	0
794	Tomb Doors	Ornately carved wooden doors depicting skeletal warriors	open	search	2	1	2	0
795	Swamp	That smells horrid	step over	examine	1	1	1	0
796	Rocks	You should be able to move these	clear	examine	1	2	1	0
797	Rocks	You should be able to move these	clear	examine	1	2	1	0
798	stalactite	Formed over thousands of years	walkto	examine	0	1	1	0
799	stalactite	Formed over thousands of years	walkto	examine	0	1	1	0
800	Spiked pit	They looks suspicious	walk over	examine	1	1	1	0
801	Lever	Seems to be some sort of winch	pull	examine	1	1	1	0
802	Cage	Seems to be mechanical 	walkto	examine	1	1	1	0
803	Cage	Seems to be mechanical 	walkto	examine	1	1	1	0
804	Rocks	More rocks!	step over	search for traps	1	1	1	0
805	Spear trap	Ouch!	walkto	examine	1	1	1	0
806	Rocks	More rocks!	step over	search	1	1	1	0
807	Rocks	More rocks!	step over	search	1	1	1	0
808	Rocks	More rocks!	step over	search	1	1	1	0
809	Rocks	More rocks!	step over	search	1	1	1	0
810	Rocks	More rocks!	step over	search	1	1	1	0
811	Rocks	More rocks!	step over	search	1	1	1	0
812	Ledge	I might be able to climb that	drop down	examine	1	1	1	0
813	Furnace	Charred bones are slowly burning inside	walkto	examine	1	1	1	0
814	Well	The remains of a warrior slump over the strange construction	drop down	examine	1	1	1	0
815	Passage	A strange metal grill covers the passage	walk down	examine	1	2	1	0
816	Passage	The passage way has swung down to a vertical position	climb up	examine	1	2	1	0
817	Passage	The passage way has swung down to a vertical position	climb up rope	examine	1	2	1	0
818	stalagmite	Formed over thousands of years	search	examine	1	1	1	0
819	Rocks	You should be able to move these	clear	search	1	2	2	0
820	Rocks	You should be able to move these	clear	search	1	2	2	0
821	Rocks	You should be able to move these	clear	search	1	2	2	0
822	Rocks	You should be able to move these	clear	search	1	2	2	0
823	Rocks	You should be able to move these	clear	search	1	2	2	0
824	Rocks	You should be able to move these	clear	search	1	2	2	0
825	Passage	Looks suspicous!	walk here	search	1	1	1	0
826	snap trap	aaaarghh	walkto	examine	1	1	1	0
827	Wooden planks	You can walk across these	walkto	examine	1	1	1	0
828	Passage	Looks suspicous!	walk here	search	1	1	1	0
829	Passage	Looks suspicous!	walk here	search	1	1	1	0
830	Flames of zamorak	Careful	search	examine	1	2	2	0
831	Platform	An ancient construction	walkto	examine	1	1	1	0
832	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
833	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
834	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
835	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
836	wall grill	It seems to filter the rotten air through the caverns	climb up	examine	1	1	1	0
837	Ledge	I might be able to make to the other side	jump off	climb up	1	1	1	0
838	wall grill	It seems to filter the rotten air through the caverns	climb up	examine	1	1	1	0
839	Dug up soil	A freshly dug pile of mud	search	examine	0	1	1	0
840	Dug up soil	A freshly dug pile of mud	search	examine	0	1	1	0
841	Pile of mud	Mud caved in from above	climb	examine	1	2	1	0
842	stalagmite	Formed over thousands of years	walkto	examine	0	1	1	0
843	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
844	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
845	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
846	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
847	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
848	Spiked pit	I don't want to go down there	walkto	examine	0	1	1	0
849	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
850	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
851	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
852	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
853	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
854	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
855	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
856	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
857	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
858	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
859	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
860	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
861	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
862	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
863	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
864	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
865	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
866	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
867	Boulder	Could be dangerous!	walkto	examine	1	2	2	0
868	crate	Someone or something has been here before us	walkto	search	1	1	1	0
869	Door	Spooky!	open	examine	1	1	1	0
870	Platform	An ancient construction	walkto	examine	1	1	1	0
871	Cage remains	Poor unicorn!	walkto	search	1	1	1	0
872	Ledge	I might be able to climb that	climb up	examine	1	1	1	0
873	Passage	Looks suspicous!	walk here	examine	1	1	1	0
874	Passage	Looks suspicous!	walk here	examine	1	1	1	0
875	Gate of Zamorak	It doesn't look very inviting	open	examine	1	3	1	0
876	Rocks	A small rocky outcrop	climb over	examine	1	1	1	0
877	Bridge support	An ancient construction	walkto	examine	1	1	1	0
878	Tomb of Iban	A clay shrine to lord iban	open	examine	1	1	2	96
879	Claws of Iban	claws of iban	walkto	examine	1	1	1	96
880	barrel	Its stinks of alcohol	empty	examine	1	1	1	0
881	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
882	Rocks	More rocks	step over	search for traps	1	1	1	0
883	Rocks	More rocks	step over	search for traps	1	1	1	0
884	Swamp	That smells horrid	walkto	examine	1	1	1	0
885	Chest	Perhaps I should search it	search	examine	1	1	1	0
886	Stone bridge	An ancient stone construction	walkto	examine	0	1	1	0
887	cage	That's no way to live	search	examine	0	1	1	0
888	cage	That's no way to live	search	examine	0	1	1	0
889	Stone steps	They lead into the darkness	walk down	examine	0	1	1	0
890	Pile of mud	Mud and rocks piled up	climb	examine	1	2	1	0
891	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
892	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
893	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
894	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
895	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
896	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
897	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
898	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
899	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
900	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
901	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
902	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
903	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
904	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
905	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
906	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
907	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
908	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
909	Stone bridge	The bridge has partly collapsed	cross	examine	1	1	3	0
910	Stone bridge	The bridge has partly collapsed	jump over	examine	1	1	3	0
911	Chest	Perhaps I should search it	search	close	1	1	1	0
912	Chest	I wonder what is inside...	open	examine	1	1	1	0
913	Pit of the Damned	The son of zamoracks alter...	walkto	examine	1	1	1	0
914	Open Door	Spooky!	open	examine	1	1	1	0
915	signpost	Observatory reception	walkto	examine	1	1	1	0
916	Stone Gate	A mystical looking object	go through	look	1	2	2	0
917	Chest	Perhaps there is something inside	search	close	1	1	1	0
918	Zodiac	A map of the twelve signs of the zodiac	walkto	examine	0	3	3	0
919	Chest	Perhaps I should search it	search	close	1	1	1	0
920	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
921	Stone steps	They lead into the darkness	walk down	examine	0	1	1	0
922	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
923	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
924	Rock	Scripture has been carved into the rock	read	examine	1	1	1	0
925	Telescope	A device for viewing the heavens	use	examine	1	1	1	0
926	Gate	The entrance to the dungeon jail	open	examine	2	1	2	0
927	sacks	These sacks feels lumpy!	search	examine	1	1	1	0
928	Ladder	the ladder goes down into a dark area	climb-down	examine	1	1	1	0
929	Chest	All these chests look the same!	open	examine	1	1	1	0
930	Chest	All these chests look the same!	open	examine	1	1	1	0
931	Bookcase	A very roughly constructed bookcase.	walkto	search	1	1	2	0
932	Iron Gate	A well wrought iron gate - it's locked.	open	search	2	1	2	0
933	Ladder	the ladder down to the cavern	climb-down	examine	1	1	1	0
934	Chest	Perhaps there is something inside	search	close	1	1	1	0
935	Chest	All these chests look the same!	open	examine	1	1	1	0
936	Chest	Perhaps there is something inside	search	close	1	1	1	0
937	Chest	All these chests look the same!	open	examine	1	1	1	0
938	Rockslide	A pile of rocks blocks your path	walkto	examine	1	1	1	0
939	Altar	An altar to the evil God Zamorak	recharge at	examine	1	2	1	0
940	column	Formed over thousands of years	walkto	examine	1	1	1	0
941	Grave of Scorpius	Here lies Scorpius: dread follower of zamorak	read	examine	1	1	3	0
942	Bank Chest	Allows you to access your bank.	open	examine	1	1	1	0
943	dwarf multicannon	fires metal balls	fire	pick up	0	1	1	0
944	Disturbed sand	Footprints in the sand show signs of a struggle	look	search	0	1	1	0
945	Disturbed sand	Footprints in the sand show signs of a struggle	look	search	0	1	1	0
946	dwarf multicannon base	bang	pick up	examine	0	1	1	0
947	dwarf multicannon stand	bang	pick up	examine	0	1	1	0
948	dwarf multicannon barrels	bang	pick up	examine	0	1	1	0
949	Cave	I wonder what's inside!	enter	examine	1	1	1	0
950	Cave	I wonder what's inside!	enter	examine	1	1	1	0
951	fence	These bridges seem hastily put up	walkto	examine	0	1	1	0
952	signpost	a signpost	read	examine	1	1	1	0
953	Rocks	I wonder if I can climb up these	climb	examine	0	1	1	0
954	Rocks	I wonder if I can climb up these	climb	examine	0	1	1	0
955	Cave entrance	A noxious smell emanates from the cave...	enter	examine	1	3	1	0
956	Chest	Perhaps I should search it	search	close	1	1	1	0
957	Chest	I wouldn't like to think where the owner is now	search	close	1	1	1	0
958	Wooden Doors	Large oak doors constantly watched by guards	open	watch	2	1	2	0
959	Pedestal	something fits on here	walkto	examine	1	1	1	96
960	bush	A leafy bush	search	examine	1	1	1	0
961	bush	A leafy bush	search	examine	1	1	1	0
962	Standard	A standard with a human skull on it	walkto	examine	1	1	1	0
963	Mining Cave	A gaping hole that leads to another section of the mine	enter	examine	1	3	1	0
964	Mining Cave	A gaping hole that leads to another section of the mine	enter	examine	1	3	1	0
965	Rocks	A small rocky outcrop	walkto	examine	1	1	1	0
966	Lift	To brings mined rocks to the surface	operate	examine	1	1	2	0
967	Mining Barrel	For loading up mined stone from below ground	walkto	search	1	1	1	0
968	Hole	I wonder where this leads...	enter	examine	1	1	1	0
969	Hole	I wonder where this leads...	enter	examine	1	1	1	0
970	Cave	I wonder what's inside!	enter	examine	1	1	1	0
971	Cave	I wonder what's inside!	enter	examine	1	1	1	0
972	Cave	I wonder what's inside!	enter	examine	1	1	1	0
973	Counter	An ogre is selling items here	steal from	examine	1	1	1	0
974	Track	Train track	look	examine	1	2	2	0
975	Cave	I wonder what's inside!	enter	examine	1	1	1	0
976	Mine Cart	A heavily constructed and often used mining cart.	look	search	1	1	1	0
977	Lift Platform	A wooden lift that is operated from the surface.	use	search	1	1	1	0
978	Chest	I wonder what is inside...	open	examine	1	1	1	0
979	Chest	I wonder what is inside...	close	examine	1	1	1	0
980	Watch tower	Constructed by the dwarven black guard	walkto	examine	0	2	2	0
981	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
982	Cave entrance	I wonder what is inside...	enter	examine	1	2	2	0
983	Pile of mud	Mud caved in from above	climb	examine	1	2	1	0
984	Cave	I wonder what's inside!	enter	examine	1	1	1	0
985	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
986	crate	A crate	search	examine	1	1	1	0
987	crate	A crate	search	examine	1	1	1	0
988	Gate	This gate barrs your way into gu'tanoth	open	examine	2	1	2	0
989	Gate	This gate barrs your way into gu'tanoth	open	examine	2	1	2	0
990	bush	A leafy bush	search	examine	1	1	1	0
991	bush	A leafy bush	search	examine	1	1	1	0
992	bush	A leafy bush	search	examine	1	1	1	0
993	bush	A leafy bush	search	examine	1	1	1	0
994	multicannon	fires metal balls	inspect	examine	1	1	1	0
995	Rocks	Some rocks are close to the egde	jump over	look at	1	1	1	0
996	Rocks	Some rocks are close to the edge	jump over	look at	1	1	1	0
997	Ladder	it's a ladder leading downwards	climb-down	examine	0	1	1	0
998	Cave entrance	I wonder what is inside...	enter	examine	1	1	1	0
999	Counter	An ogre is selling cakes here	steal from	examine	1	1	1	0
1000	Chest	Perhaps I should search it	search	close	1	1	1	0
1001	Chest	I wonder what is inside...	open	examine	1	1	1	0
1002	Chest	Perhaps I should search it	search	close	1	1	1	0
1003	Chest	I wonder what is inside...	open	examine	1	1	1	0
1004	Bookcase	A large collection of books	look	search	1	1	2	0
1005	Captains Chest	I wonder what is inside...	open	examine	1	1	1	0
1006	Experimental Anvil	An experimental anvil - for developing new techniques in forging	use	examine	1	1	1	0
1007	Rocks	A small pile of stones	search	examine	1	1	1	0
1008	Rocks	A small rocky outcrop	search	examine	1	1	1	0
1009	Column	Created by ancient mages	walkto	examine	1	1	1	0
1010	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1011	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1012	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1013	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1014	Lever	The lever is up	pull	examine	0	1	1	0
1015	Lever	The lever is down	pull	examine	0	1	1	0
1016	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1017	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
1018	Wall	Created by ancient mages	walkto	examine	1	1	1	0
1019	Gate	The gate is closed	open	examine	2	1	2	0
1020	Gate	The gate is closed	open	examine	2	1	2	0
1021	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
1022	shock	cosmic energy	walkto	examine	1	1	1	0
1023	Desk	A very strong looking table with some locked drawers.	walkto	search	1	2	1	120
1024	Cave	I wonder what's inside!	enter	examine	1	1	1	0
1025	Mining Cart	A sturdy well built mining cart with barrels full of rock on the back.	walkto	search	1	2	3	0
1026	Rock of Dalgroth	A mysterious boulder of the ogres	mine	prospect	1	2	2	0
1027	entrance	Created by ancient mages	walk through	examine	1	1	1	0
1028	Dried Cactus	It looks very spikey	walkto	examine	1	1	1	0
1029	climbing rocks	I wonder if I can climb up these	climb	examine	0	1	1	0
1030	Rocks	Strange rocks - who knows why they're wanted?	mine	prospect	1	1	1	0
1031	lightning	blimey!	walkto	examine	1	1	1	0
1032	Crude Desk	A very roughly constructed desk	walkto	search	1	2	1	120
1033	Heavy Metal Gate	This is an immense and very heavy looking gate made out of thick wrought metal	look	push	2	1	2	0
1034	Counter	An ogre is selling cakes here	steal from	examine	1	1	1	0
1035	Crude bed	A flea infested sleeping experience	rest	examine	1	1	2	0
1036	flames	looks hot!	walkto	examine	1	1	1	0
1037	Carved Rock	An ornately carved rock with a pointed recepticle	walkto	search	1	1	1	120
1038	USE	FREE SLOT PLEASE USE	walkto	examine	1	1	1	0
1039	crate	A crate used for storing materials	search	examine	1	1	1	0
1040	crate	A crate used for storing materials	search	examine	1	1	1	0
1041	barrel	Its shut	search	examine	1	1	1	0
1042	Brick	A stone brick	walkto	examine	1	1	1	0
1043	Brick	A stone brick	walkto	examine	1	1	1	0
1044	Brick	A stone brick	walkto	examine	1	3	1	0
1045	Brick	A stone brick	walkto	examine	1	3	1	0
1046	Brick	A stone brick	walkto	examine	1	2	2	0
1047	Brick	A stone brick	walkto	examine	1	1	1	0
1048	Barrier	this section is roped off	walkto	examine	1	1	1	0
1049	buried skeleton	I hope I don't meet any of these	search	examine	1	1	1	0
1050	Brick	A stone brick	walkto	examine	1	1	1	0
1051	Brick	A stone brick	walkto	examine	1	1	1	0
1052	Specimen tray	A pile of sifted earth	walkto	search	1	2	2	0
1053	winch	This winches earth from the dig hole	operate	examine	1	1	2	0
1054	crate	A crate	search	examine	1	1	1	0
1055	crate	A crate	search	examine	1	1	1	0
1056	Urn	A large ornamental urn	walkto	examine	1	1	1	0
1057	buried skeleton	I'm glad this isn't around now	search	examine	1	1	1	0
1058	panning point	a shallow where I can pan for gold	look	examine	0	1	1	0
1059	Rocks	A small rocky outcrop	walkto	examine	1	1	1	0
1060	signpost	a signpost	read	examine	1	1	1	0
1061	signpost	a signpost	read	examine	1	1	1	0
1062	signpost	a signpost	read	examine	1	1	1	0
1063	signpost	a signpost	read	examine	1	1	1	0
1064	signpost	Digsite educational centre	walkto	examine	1	1	1	0
1065	soil	soil	search	examine	0	1	1	0
1066	soil	soil	search	examine	0	1	1	0
1067	soil	soil	search	examine	0	1	1	0
1068	Gate	The gate has closed	open	examine	2	1	2	0
1069	ship	The ship is sinking	walkto	examine	2	1	2	0
1070	barrel	The ship is sinking	climb on	examine	2	1	2	0
1071	Leak	The ship is sinking	fill	examine	0	1	1	0
1072	bush	A leafy bush	search	examine	1	1	1	0
1073	bush	A leafy bush	search	examine	1	1	1	0
1074	cupboard	The cupboard is shut	open	examine	1	1	2	0
1075	sacks	Yep they're sacks	search	examine	1	1	1	0
1076	sacks	Yep they're sacks	search	examine	1	1	1	0
1077	Leak	The ship is sinking	fill	examine	0	1	1	0
1078	cupboard	The cupboard is shut	search	examine	1	1	2	0
1079	Wrought Mithril Gates	Magnificent wrought mithril gates giving access to the Legends Guild	open	search	2	1	2	0
1080	Legends Hall Doors	Solid Oak doors leading to the Hall of Legends	open	search	2	1	2	0
1081	Camp bed	Not comfortable but useful nonetheless	walkto	examine	1	1	2	0
1082	barrel	It has a lid on it - I need something to lever it off	walkto	examine	1	1	1	0
1083	barrel	I wonder what is inside...	search	examine	1	1	1	0
1084	Chest	Perhaps I should search it	search	close	1	1	1	0
1085	Chest	I wonder what is inside...	open	examine	1	1	1	0
1086	Dense Jungle Tree	Thick vegetation	chop	examine	1	1	1	0
1087	Jungle tree stump	A chopped down jungle tree	walk	examine	1	1	1	0
1088	signpost	To the digsite	walkto	examine	1	1	1	0
1089	gate	You can pass through this on the members server	open	examine	2	1	2	0
1090	Bookcase	A large collection of books	search	examine	1	1	2	0
1091	Dense Jungle Tree	An exotic looking tree	chop	examine	1	1	1	0
1092	Dense Jungle Tree	An exotic looking tree	chop	examine	1	1	1	0
1093	Spray	There's a strong wind	walkto	examine	1	1	1	0
1094	Spray	There's a strong wind	walkto	examine	1	1	1	0
1095	winch	This winches earth from the dig hole	operate	examine	1	1	2	0
1096	Brick	It seems these were put here deliberately	search	examine	1	1	1	0
1097	Rope	it's a rope leading upwards	climb-up	examine	1	1	1	0
1098	Rope	it's a rope leading upwards	climb-up	examine	1	1	1	0
1099	Dense Jungle Palm	A hardy palm tree with dense wood	chop	examine	1	1	1	0
1100	Dense Jungle Palm	A hardy palm tree with dense wood	chop	examine	1	1	1	0
1101	Trawler net	A huge net to catch little fish	inspect	examine	1	1	1	0
1102	Trawler net	A huge net to catch little fish	inspect	examine	1	1	1	0
1103	Brick	The bricks are covered in the strange compound	walkto	examine	1	1	1	0
1104	Chest	I wonder what is inside ?	open	examine	1	1	1	0
1105	Chest	Perhaps I should search it	search	close	1	1	1	0
1106	Trawler catch	Smells like fish!	search	examine	1	1	1	0
1107	Yommi Tree	An adolescent rare and mystical looking tree in	walkto	examine	1	2	2	0
1108	Grown Yommi Tree	A fully grown rare and mystical looking tree	walkto	examine	1	2	2	0
1109	Chopped Yommi Tree	A mystical looking tree that has recently been felled	walkto	examine	1	2	2	0
1110	Trimmed Yommi Tree	The trunk of the yommi tree.	walkto	examine	1	2	2	0
1111	Totem Pole	A nicely crafted wooden totem pole.	lift	examine	1	2	2	0
1112	Baby Yommi Tree	A baby Yommi tree - with a mystical aura	walkto	examine	1	2	2	0
1113	Fertile earth	A very fertile patch of earth	walkto	examine	0	2	2	0
1114	Rock Hewn Stairs	steps cut out of the living rock	climb	examine	1	2	3	0
1115	Hanging rope	A rope hangs from the ceiling	walkto	examine	1	1	1	0
1116	Rocks	A large boulder blocking the stream	move	examine	1	2	2	0
1117	Boulder	A large boulder blocking the way	walkto	smash to pieces	1	2	2	0
1118	dwarf multicannon	fires metal balls	fire	pick up	1	1	1	0
1119	dwarf multicannon base	bang	pick up	examine	1	1	1	0
1120	dwarf multicannon stand	bang	pick up	examine	1	1	1	0
1121	dwarf multicannon barrels	bang	pick up	examine	1	1	1	0
1122	rock	A rocky outcrop	climb over	examine	1	1	1	0
1123	Rock Hewn Stairs	steps cut out of the living rock	climb	examine	1	2	3	0
1124	Rock Hewn Stairs	steps cut out of the living rock	climb	examine	1	2	3	0
1125	Rock Hewn Stairs	steps cut out of the living rock	climb	examine	1	2	3	0
1126	Compost Heap	The family gardeners' compost heap	walkto	investigate	1	2	2	0
1127	beehive	An old looking beehive	walkto	investigate	1	1	1	0
1128	Drain	This drainpipe runs from the kitchen to the sewers	walkto	investigate	0	1	1	0
1129	web	An old thick spider's web	walkto	investigate	0	1	1	0
1130	fountain	There seems to be a lot of insects here	walkto	investigate	1	2	2	0
1131	Sinclair Crest	The Sinclair family crest	walkto	investigate	0	1	1	0
1132	barrel	Annas stuff - There seems to be something shiny at the bottom	walkto	search	1	1	1	0
1133	barrel	Bobs things - There seems to be something shiny at the bottom	walkto	search	1	1	1	0
1134	barrel	Carols belongings - there seems to be something shiny at the bottom	walkto	search	1	1	1	0
1135	barrel	Davids equipment - there seems to be something shiny at the bottom	walkto	search	1	1	1	0
1136	barrel	Elizabeths clothes - theres something shiny at the bottom	walkto	search	1	1	1	0
1137	barrel	Franks barrel seems to have something shiny at the bottom	walkto	search	1	1	1	0
1138	Flour Barrel	Its full of flour	walkto	take from	1	1	1	0
1139	sacks	Full of various gardening tools	walkto	investigate	1	1	1	0
1140	gate	A sturdy and secure wooden gate	walkto	investigate	2	1	2	0
1141	Dead Yommi Tree	A dead Yommi Tree - it looks like a tough axe will be needed to fell this	walkto	inspect	1	2	2	0
1142	clawspell	forces of guthix	walkto	examine	1	1	1	0
1143	Rocks	The remains of a large rock	walkto	examine	1	2	2	0
1144	crate	A crate of some kind	walkto	search	1	1	1	70
1145	Cavernous Opening	A dark and mysterious cavern	enter	search	1	3	1	0
1146	Ancient Lava Furnace	A badly damaged furnace fueled by red hot Lava - it looks ancient	look	search	1	2	2	0
1147	Spellcharge	forces of guthix	walkto	examine	1	1	1	0
1148	Rocks	A small rocky outcrop	walkto	search	1	1	1	0
1149	cupboard	The cupboard is shut	open	examine	1	1	2	0
1150	sacks	Yep they're sacks	search	examine	1	1	1	0
1151	Rock	A rocky outcrop	walkto	search	1	1	1	0
1152	Saradomin stone	A faith stone	chant to	examine	1	1	1	0
1153	Guthix stone	A faith stone	chant to	examine	1	1	1	0
1154	Zamorak stone	A faith stone	chant to	examine	1	1	1	0
1155	Magical pool	A cosmic portal	step into	examine	1	2	2	0
1156	Wooden Beam	Some sort of support - perhaps used with ropes to lower people over the hole	walkto	search	0	1	1	0
1157	Rope down into darkness	A scarey downwards trip into possible doom.	walkto	use	0	1	1	0
1158	Cave entrance	A dark cave entrance leading to the surface.	enter	examine	1	3	1	0
1159	Cave entrance	A small tunnel that leads to a large room beyond.	enter	examine	1	2	2	0
1160	Ancient Wooden Doors	The doors are locked shut	open	pick lock	2	1	2	0
1161	Table	An old rickety table	walkto	search	1	1	1	96
1162	Crude bed	Barely a bed at all	rest	search	1	1	2	0
1163	Tall Reeds	A tall plant with a tube for a stem.	walkto	search	0	1	1	0
1164	Goblin foot prints	They seem to be heading south east	walkto	examine	0	1	1	0
1165	Dark Metal Gate	A dark metalic gate which seems to be fused with the rock	open	search	2	1	2	0
1166	Magical pool	A cosmic portal	step into	examine	1	2	2	0
1167	Rope Up	A welcome rope back up and out of this dark place.	climb	examine	0	1	1	0
1168	Half buried remains	Some poor unfortunate soul	walkto	search	1	1	1	0
1169	Totem Pole	A carved and decorated totem pole	look	examine	1	1	1	0
1170	Totem Pole	A carved and decorated totem pole	look	examine	1	1	1	0
1171	Comfy bed	Its a bed - wow	rest	examine	1	2	2	0
1172	Rotten Yommi Tree	A decomposing fully grown Yommi Tree	walkto	inspect	1	2	2	0
1173	Rotten Yommi Tree	A decomposing felled Yommi Tree	walkto	inspect	1	2	2	0
1174	Rotten Yommi Tree	A decomposing Yommi Tree Trunk	walkto	inspect	1	2	2	0
1175	Rotten Totem Pole	A decomposing Totem Pole	walkto	inspect	1	2	2	0
1176	Leafy Palm Tree	A shady palm tree	walkto	shake	1	1	1	0
1177	Grand Viziers Desk	A very elegant desk - you could knock it to get the Grand Viziers attention.	walkto	knock on table	1	2	1	120
1178	Strange Barrel	It might have something inside of it.	smash	examine	1	1	1	0
1179	ship	A sturdy sailing ship	walkto	examine	0	5	3	0
1180	ship	A sturdy sailing ship	walkto	examine	0	2	3	0
1181	ship	A sturdy sailing ship	walkto	examine	0	5	3	0
1182	digsite bed	Not comfortable but useful nonetheless	sleep	examine	1	1	2	0
1183	Tea stall	A stall selling oriental infusions	walkto	steal from	1	2	2	112
1184	Boulder	A large boulder blocking the way	walkto	smash to pieces	1	2	2	0
1185	Boulder	A large boulder blocking the way	walkto	smash to pieces	1	2	2	0
1186	Damaged Earth	Disturbed earth - it will heal itself in time	walkto	examine	0	1	1	0
1187	Ladder	it's a ladder leading upwards	climb-up	examine	1	1	1	0
1188	Ladder	it's a ladder leading downwards	climb-down	examine	1	1	1	0
\.


--
-- Data for Name: item_locations; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.item_locations (id, x, y, amount, respawn) FROM stdin;
471	133	211	1	5
471	141	214	1	5
469	362	602	1	25
469	361	603	1	25
41	180	200	1	126
170	220	356	1	40
32	329	146	5	42
410	191	315	1	5
190	224	295	2	42
121	166	254	1	607
28	119	631	1	60
60	183	318	1	30
411	320	299	1	20
410	327	299	1	5
411	63	325	1	15
411	65	320	1	15
411	62	320	1	15
19	119	603	1	25
190	229	290	2	42
67	68	326	1	622
319	273	289	1	25
320	272	289	1	25
59	107	1476	1	5
16	217	453	1	37
337	170	3532	1	15
337	171	3532	1	15
21	158	3523	1	15
19	126	3511	1	25
133	586	614	1	30
262	265	102	1	30
637	129	112	3	37
637	131	112	2	37
637	131	114	3	37
33	117	669	1	30
35	138	668	1	37
11	140	648	1	37
241	160	621	1	6
241	158	621	1	6
241	160	623	1	6
241	158	623	1	6
21	121	607	1	25
156	231	508	1	25
13	130	667	1	37
375	588	523	1	10
548	585	519	1	15
13	522	462	1	30
603	440	484	1	6
18	248	596	1	123
18	248	598	1	123
18	248	600	1	123
18	248	602	1	123
18	248	604	1	123
18	248	606	1	123
18	252	596	1	123
18	252	598	1	123
18	252	600	1	123
18	252	602	1	123
18	252	604	1	123
18	252	606	1	123
18	256	596	1	123
18	256	598	1	123
18	256	600	1	123
18	256	602	1	123
18	256	604	1	123
18	256	606	1	123
18	260	596	1	123
18	260	598	1	123
18	260	600	1	123
18	260	602	1	123
18	260	604	1	123
18	260	606	1	123
621	430	559	1	75
341	430	1502	1	30
377	403	567	1	25
210	408	547	1	63
190	277	1571	1	30
190	278	1578	1	30
18	327	627	1	123
18	329	627	1	123
18	331	627	1	123
18	327	629	1	123
18	329	629	1	123
18	331	629	1	123
241	331	633	1	6
241	329	633	1	6
241	327	633	1	6
140	345	605	1	25
294	344	606	1	30
167	343	606	1	25
168	343	608	1	25
295	346	1545	1	37
386	345	1545	1	37
144	344	1545	1	25
182	278	657	1	23
194	110	518	1	26
15	107	526	1	63
14	110	526	1	30
14	110	527	1	30
135	119	523	1	25
18	119	530	1	123
18	119	532	1	123
18	119	534	1	123
18	118	530	1	123
18	118	532	1	123
18	118	534	1	123
211	127	1461	1	37
13	124	1458	1	37
10	150	3330	3	60
186	148	3337	1	60
186	149	3337	1	60
152	150	3336	1	610
152	152	3335	1	610
10	147	3337	10	60
10	149	3335	42	60
10	149	3339	10	60
10	150	3335	66	60
10	150	3339	10	60
10	152	3337	10	60
251	121	465	1	37
21	120	456	1	30
10	118	452	4	61
10	121	449	7	61
189	377	1446	1	620
16	306	522	1	37
251	296	3329	1	37
135	234	507	1	25
228	231	548	1	30
228	230	548	1	30
228	229	548	1	30
228	231	549	1	30
228	230	549	1	30
228	229	549	1	30
228	231	550	1	30
228	230	550	1	30
228	229	550	1	30
19	122	603	1	25
41	361	2458	1	126
140	133	661	1	25
135	134	660	1	25
18	141	3489	1	123
17	141	3488	1	37
62	129	1601	1	37
14	139	2537	1	25
14	140	2537	1	25
14	140	2538	1	25
14	139	2536	1	25
36	77	1618	1	60
17	64	677	1	37
21	60	672	1	60
18	79	693	1	123
140	78	693	1	25
31	67	587	3	37
32	74	583	4	37
33	183	649	1	37
0	602	699	1	60
20	606	700	1	60
20	606	703	1	60
5	602	704	1	60
20	600	702	1	60
12	609	703	1	600
20	611	701	1	60
20	614	703	1	60
20	603	711	1	60
20	607	711	1	60
20	608	709	1	60
0	608	707	1	60
20	615	707	1	60
20	604	706	1	60
20	601	707	1	60
601	639	737	1	6
10	470	672	3	60
10	470	673	3	60
10	471	674	3	60
10	468	674	3	60
10	470	675	3	60
799	574	583	1	60
799	574	585	1	60
799	574	586	1	60
757	597	543	1	60
251	432	1424	1	37
465	429	1434	1	60
319	223	624	1	25
320	222	624	1	25
337	176	480	1	60
135	177	2370	1	25
252	175	2370	1	60
140	176	2370	1	25
143	181	2374	1	30
338	181	1424	1	60
341	180	1424	1	60
251	178	1424	1	37
252	177	1424	1	60
252	176	1424	1	60
596	367	3372	1	60
236	90	544	1	6
236	91	551	1	6
55	86	542	1	6
236	79	540	1	6
55	83	548	1	6
388	264	1402	1	63
389	265	1402	1	23
193	231	500	1	60
193	233	500	1	60
193	232	497	1	60
132	234	500	1	60
132	231	497	1	60
10	223	3280	1	60
10	223	3279	4	60
20	216	3289	1	25
20	212	3292	1	25
20	216	3299	1	24
20	205	3288	1	24
20	210	3288	1	24
20	207	3283	1	24
20	209	3283	1	24
11	197	3277	3	37
11	194	3267	1	37
219	201	3234	1	127
219	204	3232	1	127
219	209	3236	1	127
219	208	3240	1	127
186	204	3295	1	60
99	198	3307	1	63
7	106	489	1	607
6	106	487	1	613
128	106	485	1	613
1	105	484	1	613
0	103	484	1	613
132	307	1475	1	180
183	310	1479	1	25
21	35	542	1	60
156	33	539	1	60
140	13	541	1	25
1111	15	541	1	60
21	17	536	1	60
21	23	527	1	60
21	23	526	1	60
21	16	524	1	60
20	16	522	1	60
20	14	519	1	60
21	26	518	1	60
21	27	518	1	60
20	25	506	1	60
21	13	509	1	60
1284	11	3375	1	60
1284	9	3373	1	60
1284	5	3377	1	60
27	8	3401	1	60
27	17	3401	1	60
33	339	703	3	60
31	408	3532	1	30
31	405	3537	1	30
31	409	3543	1	30
31	414	3540	1	30
31	413	3534	1	30
219	404	3518	1	127
10	415	3476	100	60
10	419	3477	100	60
10	420	3484	100	60
10	415	3483	100	60
13	509	506	1	60
4	139	437	1	63
34	146	438	1	37
501	333	434	1	61
413	123	263	1	15
413	130	259	1	15
342	543	3273	1	25
725	540	3273	1	15
118	236	180	1	120
36	232	177	10	37
20	233	176	1	30
714	580	3524	1	37
168	544	576	1	25
167	543	576	1	25
16	116	710	1	37
164	171	161	1	603
46	333	145	3	610
40	67	177	3	120
40	69	178	4	120
11	212	367	4	20
11	215	373	3	20
11	222	370	4	20
11	217	366	2	20
11	213	364	2	20
11	213	367	2	20
190	222	298	2	42
190	228	294	2	42
190	233	296	2	42
41	263	345	1	126
32	255	347	1	40
34	259	347	1	40
33	257	347	1	40
35	261	347	1	37
36	263	347	1	37
31	253	347	1	40
18	252	453	1	123
18	251	453	1	123
18	250	453	1	123
18	249	453	1	123
18	252	455	1	123
18	251	455	1	123
18	250	455	1	123
18	249	455	1	123
18	252	457	1	123
18	251	457	1	123
18	250	457	1	123
18	249	457	1	123
622	437	536	1	74
622	437	538	1	74
622	437	540	1	74
622	437	543	1	74
10	91	3301	1	60
35	129	3292	1	37
36	129	3290	1	37
34	146	3296	6	37
10	150	3300	7	60
10	155	3301	5	60
219	160	3299	1	127
12	161	3295	1	600
10	158	3292	8	60
10	151	3290	7	60
18	137	612	1	123
18	139	612	1	123
18	141	612	1	123
18	143	612	1	123
18	145	612	1	123
18	147	612	1	123
18	149	612	1	123
18	151	612	1	123
18	153	612	1	123
18	137	609	1	123
18	139	609	1	123
18	141	609	1	123
18	143	609	1	123
18	145	609	1	123
18	147	609	1	123
18	149	609	1	123
18	151	609	1	123
18	153	609	1	123
18	137	606	1	123
18	139	606	1	123
18	141	606	1	123
18	143	606	1	123
18	145	606	1	123
18	147	606	1	123
18	149	606	1	123
18	151	606	1	123
18	153	606	1	123
18	137	603	1	123
18	139	603	1	123
18	141	603	1	123
18	143	603	1	123
18	145	603	1	123
18	147	603	1	123
18	149	603	1	123
18	151	603	1	123
18	153	603	1	123
18	137	600	1	123
18	139	600	1	123
18	141	600	1	123
18	143	600	1	123
18	145	600	1	123
18	147	600	1	123
18	149	600	1	123
18	151	600	1	123
18	153	600	1	123
20	568	3328	1	15
20	566	3326	1	15
20	569	3326	1	15
20	566	3324	1	15
20	569	3322	1	15
20	566	3319	1	15
20	569	3318	1	15
10	566	3320	1	60
10	567	3321	1	60
10	566	3318	1	60
10	568	3317	1	60
164	567	3318	1	603
283	565	3317	1	60
27	218	3521	1	30
104	206	3380	1	60
208	232	3382	1	30
176	212	1497	1	6
213	208	546	1	30
139	209	2442	1	30
166	208	2438	1	25
181	212	2441	1	30
181	214	2439	1	30
177	221	546	1	6
211	197	554	1	37
144	200	551	1	25
21	201	552	1	30
17	104	1476	1	37
59	105	1476	1	5
62	107	1477	1	37
16	107	1478	1	37
537	362	1437	1	100
539	351	491	1	12
1093	58	2327	1	30
1100	167	524	1	1
1100	167	522	1	1
1100	165	523	1	1
1100	165	521	1	1
1100	163	522	1	1
1046	617	1435	1	60
21	616	640	1	30
716	619	3500	1	30
20	622	3499	1	30
20	622	3500	1	30
20	616	3497	1	30
20	613	3498	1	30
20	612	3498	1	30
20	612	3497	1	30
783	166	676	6	30
783	166	674	6	30
783	164	679	6	30
783	163	679	6	30
783	163	675	6	30
783	173	679	6	30
783	171	676	6	30
783	168	679	6	30
783	176	674	6	30
783	167	681	6	30
783	165	683	6	30
783	163	683	6	30
783	174	682	6	30
783	170	681	6	30
783	169	684	6	30
783	169	684	6	30
783	164	691	6	30
783	163	694	6	30
783	160	690	6	30
783	174	689	6	30
783	170	689	6	30
783	169	693	6	30
783	156	679	6	30
783	156	678	6	30
783	152	675	6	30
783	157	686	6	30
783	155	680	6	30
783	154	685	6	30
783	154	682	6	30
783	153	683	6	30
783	153	683	1	30
783	152	687	6	30
783	159	694	6	30
783	157	690	6	30
783	155	693	6	30
783	155	688	6	30
783	154	694	6	30
783	152	691	6	30
783	152	689	6	30
783	156	698	6	30
783	160	697	6	30
783	151	687	6	30
783	150	686	6	30
783	150	682	6	30
783	150	691	6	30
783	150	691	6	30
783	148	695	6	30
783	146	688	6	30
783	150	699	6	30
135	176	657	1	25
769	487	614	1	30
769	487	610	1	30
778	485	614	1	30
769	485	612	1	30
769	484	612	1	30
776	482	615	1	30
769	495	615	1	30
769	493	614	1	30
769	492	615	1	30
769	490	615	1	30
769	490	611	1	30
769	490	610	1	30
769	489	614	1	30
769	489	611	1	30
769	489	610	1	30
769	489	609	1	30
776	488	609	1	30
769	495	616	1	30
769	492	619	1	30
769	491	620	1	39
769	488	619	1	30
769	498	615	1	30
769	497	614	1	30
769	496	614	1	30
769	487	1564	1	30
769	487	1563	1	30
769	487	1561	1	30
769	495	1559	1	39
769	494	1564	1	30
769	493	1563	1	30
769	492	1564	1	30
769	491	1564	1	30
769	490	1563	1	30
769	490	1562	1	30
769	490	1561	1	30
769	489	1564	1	30
769	489	1563	1	30
769	489	1562	1	30
769	488	1564	1	30
769	488	1562	1	30
410	488	540	1	5
410	486	540	1	5
410	486	539	1	5
410	485	540	1	5
1025	715	680	1	30
20	696	3511	1	30
20	698	3512	1	30
20	697	3527	1	30
20	711	3511	1	30
20	711	3515	1	30
20	709	3516	1	30
168	707	3520	1	25
20	714	3505	1	30
20	712	3510	1	30
20	713	3518	1	30
20	713	3513	1	30
20	673	3511	1	30
20	682	3534	1	30
20	678	3527	1	30
20	690	3536	1	30
20	676	3550	1	30
20	718	3543	1	30
20	706	3538	1	30
20	690	3516	1	30
20	692	3514	1	30
568	710	695	1	30
469	358	615	1	25
765	544	453	1	25
765	544	460	1	25
765	566	457	1	25
765	564	443	1	25
767	615	577	1	30
211	614	577	1	37
138	644	565	1	30
516	622	631	1	30
467	368	3352	1	25
467	372	3354	1	25
467	368	3356	1	25
467	368	3350	1	25
706	563	600	1	30
1205	485	389	1	1
1204	484	388	1	1
731	592	3483	1	25
7	589	3481	1	30
730	579	3474	1	30
728	609	3465	1	30
727	583	3456	1	30
729	616	3484	1	30
738	388	24	1	30
738	426	15	1	30
738	421	976	1	30
746	418	1924	1	30
738	517	976	1	30
746	514	1924	1	30
1285	89	516	1	30
14	221	622	1	15
82	316	1608	1	60
14	326	670	1	15
14	325	670	1	15
11	323	667	1	15
11	334	1509	2	60
87	310	1465	1	45
21	320	442	1	25
113	298	3344	1	60
14	70	446	1	15
14	67	441	1	15
14	64	446	1	15
14	57	444	1	15
14	62	448	1	15
14	59	451	1	15
14	62	435	1	15
14	56	438	1	15
14	71	436	1	15
14	69	433	1	15
14	68	436	1	15
14	64	433	1	15
14	54	433	1	15
14	51	441	1	15
412	310	410	1	6
335	519	663	1	10
10	367	3354	1	60
10	370	3351	1	60
10	371	3356	1	60
738	522	15	1	30
10	350	3322	4	60
10	351	3336	2	60
31	460	1385	2	15
18	497	401	1	123
18	497	400	1	123
18	497	398	1	123
18	497	397	1	123
621	558	483	1	60
252	561	549	1	5
341	547	531	1	5
379	588	519	1	15
427	612	2488	1	120
10	617	3483	1	61
20	614	3483	1	12
13	612	3482	1	10
604	606	3484	1	15
20	640	648	1	30
20	647	653	1	30
20	654	653	1	30
20	643	651	1	30
20	642	649	1	30
20	638	644	1	30
20	637	644	1	30
20	635	641	1	30
20	634	639	1	30
20	643	640	1	30
20	644	642	1	30
20	646	645	1	30
20	667	646	1	30
20	667	648	1	30
20	663	649	1	30
20	665	643	1	30
20	664	641	1	30
20	699	651	1	30
20	699	652	1	30
20	698	648	1	30
20	701	650	1	30
20	700	647	1	30
20	691	656	1	30
20	690	649	1	30
20	694	652	1	30
20	691	651	1	30
20	691	652	1	30
20	690	652	1	30
20	690	655	1	30
20	691	654	1	30
20	692	653	1	30
20	688	655	1	30
20	687	657	1	30
20	642	645	1	30
20	643	631	1	30
20	642	630	1	30
20	650	640	1	30
20	646	639	1	30
20	646	638	1	30
20	649	638	1	30
20	646	634	1	30
4	650	650	1	15
4	657	650	1	15
15	645	650	1	63
5	665	647	1	20
134	691	640	1	10
412	693	647	1	15
412	692	652	1	15
412	691	650	1	15
412	689	652	1	15
412	698	649	1	15
412	701	648	1	15
412	700	651	1	15
625	696	678	1	20
20	661	853	1	30
20	661	852	1	30
20	660	853	1	30
20	659	855	1	30
38	583	856	1	30
20	607	800	1	30
20	608	804	1	30
20	607	798	1	30
20	607	792	1	30
20	605	798	1	30
20	605	797	1	30
20	604	795	1	30
20	608	796	1	30
411	69	329	1	15
411	63	318	1	15
410	71	328	1	5
410	59	324	1	5
410	63	322	1	5
410	68	325	1	5
410	67	320	1	5
410	66	317	1	5
35	71	284	6	37
413	96	279	1	15
413	116	268	1	15
413	131	278	1	60
412	123	263	1	15
412	137	261	1	15
63	63	256	1	622
36	113	385	2	37
36	109	391	2	37
36	109	399	2	37
10	113	405	2	60
60	73	281	1	15
152	173	250	1	605
423	161	201	1	120
10	133	203	8	60
10	131	201	6	60
465	151	184	1	5
465	156	190	1	5
411	304	295	1	5
411	315	294	1	5
411	315	292	1	5
411	314	292	1	5
411	329	293	1	5
411	307	302	1	5
411	306	303	1	5
411	318	299	1	5
411	312	298	1	5
411	312	298	1	5
411	323	299	1	5
411	319	304	1	5
411	315	304	1	5
411	325	309	1	5
411	320	311	1	5
411	318	313	1	5
411	297	309	1	5
380	317	291	1	25
380	318	301	1	25
376	315	293	1	25
338	324	303	1	25
410	313	294	1	5
410	330	292	1	5
410	304	302	1	5
410	324	302	1	5
410	322	298	1	5
410	320	303	1	5
410	310	306	1	5
410	319	309	1	5
410	314	305	1	5
410	326	307	1	5
410	322	307	1	5
410	316	312	1	5
410	326	312	1	5
410	303	299	1	5
410	303	299	1	5
410	303	297	1	5
410	301	304	1	5
10	321	309	5	60
10	317	310	4	60
10	323	313	4	60
10	326	311	5	60
20	321	300	1	15
20	328	298	1	15
412	312	412	1	6
6	242	198	1	84
410	191	312	1	5
410	187	321	1	5
410	171	321	1	5
410	169	322	1	5
410	160	314	1	5
410	166	315	1	5
410	161	320	1	5
411	177	316	1	15
411	183	324	1	15
411	175	324	1	15
5	170	310	1	84
63	181	315	1	622
2	120	222	1	620
156	319	659	1	25
152	495	659	1	60
920	714	1418	1	10
801	650	564	1	10
801	642	565	1	10
167	631	593	1	25
13	631	595	1	45
181	634	598	1	30
181	633	600	1	30
181	635	601	1	30
411	678	569	1	30
411	687	570	1	30
411	688	567	1	30
411	280	2963	1	30
476	276	2961	1	60
444	269	2967	1	45
936	658	760	1	37
936	660	762	1	37
936	663	755	1	37
936	658	757	1	37
20	596	695	1	60
20	665	760	1	30
20	664	762	1	30
20	660	735	1	30
20	659	734	1	30
20	664	735	1	30
722	543	3374	1	30
724	551	3327	1	45
723	563	3354	1	60
723	560	3352	1	60
732	533	3302	1	10
94	51	718	1	60
1099	65	727	1	37
141	85	800	1	60
342	84	1744	1	60
156	75	3626	1	30
156	75	3625	1	30
156	75	3624	1	30
156	75	3622	1	30
986	63	3641	1	60
986	62	3641	1	60
211	291	524	1	37
907	759	3449	1	30
143	761	3448	1	30
60	760	3451	1	60
138	760	3448	1	30
20	747	3432	1	60
95	763	3451	1	60
60	775	3527	1	60
60	768	3533	1	60
5	768	3523	1	60
21	756	667	1	30
1003	804	3509	1	60
346	742	619	1	30
132	743	625	1	120
357	745	626	1	30
133	678	3445	1	120
132	678	3451	1	120
410	708	3449	1	5
410	705	3450	1	5
410	726	3442	1	5
996	726	3448	1	60
991	728	3434	1	60
993	752	3440	1	60
992	755	3442	1	60
994	708	3473	1	60
994	715	3482	1	60
237	721	3471	1	30
909	721	3493	1	60
20	752	3430	1	60
20	748	3426	1	60
20	743	3432	1	60
20	742	3429	1	60
20	745	3427	1	60
20	740	3427	1	60
237	605	801	1	60
20	604	799	1	45
43	78	451	1	124
702	703	654	1	63
288	245	176	1	45
1283	300	758	1	15
1283	299	755	1	15
1283	299	760	1	15
1283	306	758	1	15
1283	304	754	1	15
1283	305	761	1	15
1283	252	758	1	15
1283	251	755	1	15
1283	251	760	1	15
1283	258	758	1	15
1283	256	754	1	15
1283	257	761	1	15
36	149	454	2	40
83	136	279	1	45
27	131	280	1	45
538	362	1439	1	30
20	591	690	1	60
8	615	2603	1	600
12	616	2604	1	600
3	617	2605	1	600
0	616	2606	1	600
1	617	2606	1	600
6	620	2606	1	600
20	651	3594	1	45
20	648	3596	1	45
1086	650	3592	1	45
10	628	3595	1	45
10	629	3595	1	45
20	628	3596	1	45
20	627	3593	1	45
20	628	3589	1	45
10	631	3594	1	45
20	652	3557	1	45
20	648	3557	1	45
20	651	3559	1	45
1086	650	3560	1	45
20	652	3559	1	45
1036	635	3566	1	60
165	642	3564	1	60
165	642	3563	1	60
20	626	3568	1	45
20	629	3569	1	45
20	630	3570	1	45
20	629	3570	1	45
20	628	3569	1	45
20	630	3568	1	45
20	628	3567	1	45
319	631	3583	1	25
20	661	3615	1	25
20	661	3619	1	25
20	662	3617	1	25
20	662	3618	1	25
20	663	3618	1	25
1087	664	3614	1	25
20	663	3620	1	25
20	663	3634	1	25
20	663	3637	1	25
20	664	3636	1	25
20	665	3635	1	25
20	666	3632	1	25
20	666	3637	1	25
20	667	3634	1	25
20	667	3636	1	25
20	668	3634	1	25
20	669	3633	1	25
10	358	3626	10	60
1284	11	3327	1	30
1284	9	3325	1	30
1284	5	3329	1	30
1174	11	3348	1	30
27	8	3353	1	30
27	8	3354	1	30
27	6	3355	1	30
27	17	3353	1	30
27	17	3355	1	30
156	25	3347	1	60
156	28	3347	1	60
156	30	3345	1	60
156	30	3346	1	60
21	22	3343	1	60
21	23	3342	1	60
21	24	3341	1	60
1174	11	3396	1	30
413	700	648	1	60
362	614	3564	1	25
1264	426	3708	1	30
895	721	440	1	75
897	723	440	1	75
897	727	440	1	75
897	720	439	1	75
895	725	439	1	75
895	733	439	1	75
895	721	438	1	75
897	720	437	1	75
895	728	437	1	75
895	735	437	1	75
895	720	436	1	75
897	724	436	1	75
895	726	436	1	75
897	721	435	1	75
897	724	435	1	75
897	725	435	1	75
895	721	434	1	75
895	729	434	1	75
897	732	434	1	75
895	733	434	1	75
895	735	434	1	75
897	724	433	1	75
897	729	433	1	75
897	734	433	1	75
14	220	622	1	25
671	207	3206	1	37
0	536	601	1	350
109	450	1621	1	622
109	450	1621	1	622
833	693	1403	1	68
833	693	1403	1	68
833	693	1404	1	68
833	693	1404	1	68
833	691	1401	1	68
833	691	1401	1	68
20	664	3616	1	25
20	664	3619	1	25
20	665	3618	1	25
20	665	3620	1	25
20	666	3617	1	25
\.


--
-- Data for Name: item_wieldable; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.item_wieldable (id, sprite, type, armour_points, magic_points, prayer_points, range_points, weapon_aim_points, weapon_power_points, pos, femaleonly) FROM stdin;
59	107	16	0	0	0	24	0	0	4	f
60	107	16	0	0	0	9	0	0	4	f
15	46	64	4	0	0	0	0	0	6	f
104	70	32	3	0	0	0	0	0	5	f
105	72	32	7	0	0	0	0	0	5	f
106	73	32	10	0	0	0	0	0	5	f
107	74	32	14	0	0	0	0	0	5	f
108	13	33	4	0	0	0	0	0	0	f
109	15	33	9	0	0	0	0	0	0	f
110	16	33	13	0	0	0	0	0	0	f
111	17	33	19	0	0	0	0	0	0	f
24	80	1024	0	0	0	0	0	0	8	f
45	80	1024	0	0	0	0	0	0	8	f
112	18	33	30	0	0	0	0	0	0	f
113	21	64	9	0	0	0	0	0	6	f
114	23	64	19	0	0	0	0	0	6	f
115	24	64	27	0	0	0	0	0	6	f
116	25	64	38	0	0	0	0	0	6	f
117	28	322	14	0	0	0	0	0	1	f
118	30	322	31	0	0	0	0	0	1	f
119	31	322	44	0	0	0	0	0	1	f
120	32	322	63	0	0	0	0	0	1	f
121	39	644	17	0	0	0	0	0	2	f
122	40	644	22	0	0	0	0	0	2	f
123	41	644	31	0	0	0	0	0	2	f
124	98	8	5	0	0	0	0	0	3	f
125	100	8	12	0	0	0	0	0	3	f
126	101	8	17	0	0	0	0	0	3	f
127	102	8	24	0	0	0	0	0	3	f
128	98	8	6	0	0	0	0	0	3	f
129	100	8	14	0	0	0	0	0	3	f
130	101	8	20	0	0	0	0	0	3	f
131	102	8	29	0	0	0	0	0	3	f
183	63	2048	2	0	0	0	0	0	11	f
196	33	322	40	0	0	0	0	0	1	f
206	37	644	7	0	0	0	0	0	2	f
209	64	2048	2	0	0	0	0	0	11	f
214	92	640	7	0	0	0	0	0	7	f
215	93	640	10	0	0	0	0	0	7	f
225	94	640	17	0	0	0	0	0	7	f
226	95	640	22	0	0	0	0	0	7	f
227	96	640	31	0	0	0	0	0	7	f
229	65	2048	2	0	0	0	0	0	11	f
230	19	33	12	0	0	0	0	0	0	f
248	43	644	20	0	0	0	0	0	2	f
308	55	322	14	0	0	0	0	0	1	t
309	57	322	31	0	0	0	0	0	1	t
310	58	322	44	0	0	0	0	0	1	t
0	118	16	0	0	0	0	7	5	4	f
1	49	16	0	0	0	0	7	7	4	f
12	110	16	0	0	0	0	5	7	4	f
28	49	16	0	0	0	0	5	5	4	f
52	49	16	0	0	0	0	9	9	4	f
62	48	16	0	0	0	0	4	4	4	f
63	50	16	0	0	0	0	8	8	4	f
64	51	16	0	0	0	0	11	11	4	f
65	52	16	0	0	0	0	15	15	4	f
66	48	16	0	0	0	0	6	6	4	f
67	50	16	0	0	0	0	13	13	4	f
68	51	16	0	0	0	0	18	18	4	f
69	52	16	0	0	0	0	24	24	4	f
70	48	16	0	0	0	0	7	7	4	f
71	49	16	0	0	0	0	9	9	4	f
72	50	16	0	0	0	0	16	16	4	f
73	51	16	0	0	0	0	22	22	4	f
74	52	16	0	0	0	0	31	31	4	f
75	53	16	0	0	0	0	49	49	4	f
184	77	64	0	5	0	0	0	0	6	f
185	78	32	0	3	0	0	0	0	5	f
199	79	32	0	3	0	0	0	0	5	f
314	81	1024	0	6	0	0	0	0	8	f
197	123	16	0	20	0	0	7	3	4	f
188	108	24	0	0	0	7	0	0	3	f
189	108	24	0	0	0	4	0	0	3	f
311	59	322	63	0	0	0	0	0	1	t
312	56	322	20	0	0	0	0	0	1	t
313	61	322	40	0	0	0	0	0	1	t
315	81	1024	7	0	0	0	0	0	8	f
399	75	32	22	0	0	0	0	0	5	f
400	26	64	65	0	0	0	0	0	6	f
401	34	322	80	0	0	0	0	0	1	f
402	42	644	49	0	0	0	0	0	2	f
403	103	8	38	0	0	0	0	0	3	f
404	103	8	46	0	0	0	0	0	3	f
406	97	640	49	0	0	0	0	0	7	f
182	10	1024	0	0	0	0	0	0	8	f
407	60	322	80	0	0	0	0	0	1	t
186	81	1024	0	0	0	0	0	0	8	f
187	82	128	0	0	0	0	0	0	7	f
191	11	1024	0	0	0	0	0	0	8	f
192	9	32	0	0	0	0	0	0	5	f
194	90	128	0	0	0	0	0	0	7	f
195	89	128	0	0	0	0	0	0	7	f
420	105	8	3	0	0	0	0	0	3	f
431	27	64	24	0	0	0	0	0	6	f
432	104	8	15	0	0	0	0	0	3	f
433	104	8	17	0	0	0	0	0	3	f
434	159	640	20	0	0	0	0	0	7	f
216	83	64	0	0	0	0	0	0	6	f
470	76	32	9	0	0	0	0	0	5	f
511	66	2048	2	0	0	0	0	0	11	f
512	67	2048	2	0	0	0	0	0	11	f
513	68	2048	2	0	0	0	0	0	11	f
514	69	2048	2	0	0	0	0	0	11	f
288	81	1024	0	0	0	0	0	0	8	f
289	81	1024	0	0	0	0	0	0	8	f
290	81	1024	0	0	0	0	0	0	8	f
291	81	1024	0	0	0	0	0	0	8	f
292	81	1024	0	0	0	0	0	0	8	f
301	81	1024	0	0	0	0	0	0	8	f
302	81	1024	0	0	0	0	0	0	8	f
303	81	1024	0	0	0	0	0	0	8	f
304	81	1024	0	0	0	0	0	0	8	f
305	81	1024	0	0	0	0	0	0	8	f
556	156	256	3	0	0	0	0	0	10	f
198	123	16	0	15	0	0	12	7	4	f
235	80	1024	0	0	0	0	3	0	8	f
203	112	16	0	0	0	0	12	16	4	f
204	113	16	0	0	0	0	17	23	4	f
205	109	16	0	0	0	0	7	9	4	f
217	48	16	0	0	0	0	4	4	4	f
265	49	16	0	0	0	0	9	9	4	f
307	117	16	0	0	0	0	7	5	4	f
648	108	24	0	0	0	12	0	0	3	f
649	108	24	0	0	0	9	0	0	3	f
650	108	24	0	0	0	17	0	0	3	f
651	108	24	0	0	0	14	0	0	3	f
316	81	1024	0	0	0	0	0	10	8	f
594	162	16	0	0	0	0	69	75	4	f
544	81	1024	0	0	0	0	0	0	8	f
576	150	32	0	0	0	0	0	0	5	f
577	151	32	0	0	0	0	0	0	5	f
578	152	32	0	0	0	0	0	0	5	f
579	153	32	0	0	0	0	0	0	5	f
580	154	32	0	0	0	0	0	0	5	f
581	155	32	0	0	0	0	0	0	5	f
609	164	32	0	0	0	0	0	0	5	f
610	81	1024	0	0	0	0	0	0	8	f
385	80	1024	0	0	7	0	0	0	8	f
388	84	64	0	0	5	0	0	0	6	f
389	88	128	0	0	4	0	0	0	7	f
607	86	64	0	0	5	0	0	0	6	f
608	87	128	0	0	4	0	0	0	7	f
396	53	16	0	0	0	0	25	25	4	f
397	53	16	0	0	0	0	40	40	4	f
398	53	16	0	0	0	0	44	44	4	f
405	114	16	0	0	0	0	26	36	4	f
423	54	16	0	0	0	0	10	8	4	f
424	54	16	0	0	0	0	15	12	4	f
425	54	16	0	0	0	0	19	15	4	f
426	54	8216	0	0	0	0	28	22	4	f
427	54	16	0	0	0	0	18	14	4	f
428	115	16	0	0	0	0	10	14	4	f
429	115	16	0	0	0	0	18	25	4	f
430	122	16	0	0	0	0	14	64	4	f
509	123	16	0	0	0	0	7	3	4	f
559	49	16	0	0	0	0	5	5	4	f
560	48	16	0	0	0	0	4	4	4	f
561	50	16	0	0	0	0	8	8	4	f
562	51	16	0	0	0	0	11	11	4	f
563	53	16	0	0	0	0	25	25	4	f
564	52	16	0	0	0	0	15	15	4	f
565	54	16	0	0	0	0	10	8	4	f
593	163	16	0	0	0	0	71	71	4	f
693	81	1024	0	0	0	0	0	0	8	f
652	108	24	0	0	0	22	0	0	3	f
653	108	24	0	0	0	19	0	0	3	f
721	172	1024	0	0	0	0	0	0	8	f
654	108	24	0	0	0	27	0	0	3	f
726	80	1024	0	0	0	0	0	0	8	f
655	108	24	0	0	0	24	0	0	3	f
656	108	24	0	0	0	32	0	0	3	f
657	108	24	0	0	0	29	0	0	3	f
760	46	64	0	0	0	0	0	0	6	f
761	177	644	0	0	0	0	0	0	2	f
766	178	32	0	0	0	0	0	0	5	f
782	80	1024	0	0	0	0	0	0	8	f
1013	0	16	0	0	0	34	0	0	4	f
802	86	64	0	0	0	0	0	0	6	f
826	81	1024	0	0	0	0	0	0	8	f
827	181	16	0	0	0	0	0	0	4	f
828	182	32	0	0	0	0	0	0	5	f
831	185	32	0	0	0	0	0	0	5	f
832	186	32	0	0	0	0	0	0	5	f
836	187	128	0	0	0	0	0	0	7	f
837	188	128	0	0	0	0	0	0	7	f
838	189	128	0	0	0	0	0	0	7	f
839	190	128	0	0	0	0	0	0	7	f
840	191	128	0	0	0	0	0	0	7	f
841	192	32	0	0	0	0	0	0	5	f
842	193	32	0	0	0	0	0	0	5	f
843	194	32	0	0	0	0	0	0	5	f
844	195	32	0	0	0	0	0	0	5	f
845	196	32	0	0	0	0	0	0	5	f
846	197	64	0	0	0	0	0	0	6	f
847	198	64	0	0	0	0	0	0	6	f
848	199	64	0	0	0	0	0	0	6	f
849	200	64	0	0	0	0	0	0	6	f
850	201	64	0	0	0	0	0	0	6	f
852	81	1024	0	0	0	0	0	0	8	f
1015	0	16	0	0	0	34	0	0	4	f
971	209	32	0	3	0	0	0	0	5	f
682	123	16	0	20	0	0	35	35	4	f
683	123	16	0	20	0	0	35	35	4	f
684	123	16	0	20	0	0	35	35	4	f
685	123	16	0	20	0	0	35	35	4	f
1000	210	8216	0	20	0	0	50	50	4	f
1009	81	1024	0	0	0	0	0	0	8	f
1010	81	1024	0	0	0	0	0	0	8	f
1011	81	1024	0	0	0	0	0	0	8	f
1019	87	128	0	0	0	0	0	0	7	f
1020	86	64	0	0	0	0	0	0	6	f
725	123	16	0	0	0	0	7	3	4	f
754	110	16	0	0	0	0	8	12	4	f
702	85	64	0	0	5	0	0	0	6	f
703	91	128	0	0	4	0	0	0	7	f
807	228	64	0	0	5	0	0	0	6	f
808	159	128	0	0	4	0	0	0	7	f
1022	215	128	0	0	0	0	0	0	7	f
1023	214	64	0	0	0	0	0	0	6	f
1028	80	1024	0	0	0	0	0	0	8	f
1264	78	32	0	3	0	0	0	0	5	f
1024	0	16	0	0	0	34	0	0	4	f
1088	181	16	0	0	0	0	0	0	4	f
1089	181	16	0	0	0	0	0	0	4	f
1090	181	16	0	0	0	0	0	0	4	f
1091	181	16	0	0	0	0	0	0	4	f
1092	220	16	0	0	0	0	0	0	4	f
1068	0	16	0	0	0	34	0	0	4	f
1069	0	16	0	0	0	34	0	0	4	f
1070	0	16	0	0	0	34	0	0	4	f
1122	0	16	0	0	0	34	0	0	4	f
1123	0	16	0	0	0	34	0	0	4	f
1124	0	16	0	0	0	34	0	0	4	f
1125	0	16	0	0	0	34	0	0	4	f
1135	181	16	0	0	0	0	0	0	4	f
1136	181	16	0	0	0	0	0	0	4	f
1137	181	16	0	0	0	0	0	0	4	f
1138	181	16	0	0	0	0	0	0	4	f
1139	181	16	0	0	0	0	0	0	4	f
1140	220	16	0	0	0	0	0	0	4	f
1156	218	32	0	0	0	0	0	0	5	f
1126	0	16	0	0	0	34	0	0	4	f
1194	80	1024	0	0	0	0	0	0	8	f
1127	0	16	0	0	0	34	0	0	4	f
1075	49	16	0	0	0	39	5	5	4	f
1076	48	16	0	0	0	39	4	4	4	f
1077	50	16	0	0	0	39	8	8	4	f
1078	51	16	0	0	0	39	11	11	4	f
1079	52	16	0	0	0	39	15	15	4	f
2	99	8	9	0	0	0	0	0	3	f
1224	80	1024	0	0	0	0	0	0	8	f
3	99	8	8	0	0	0	0	0	3	f
1234	77	64	0	0	0	0	0	0	6	f
4	106	8	3	0	0	0	0	0	3	f
1237	164	32	0	0	0	0	0	0	5	f
1239	0	16	0	0	0	0	0	0	4	f
5	71	32	4	0	0	0	0	0	5	f
6	14	33	6	0	0	0	0	0	0	f
7	22	64	12	0	0	0	0	0	6	f
8	29	322	20	0	0	0	0	0	1	f
9	38	644	10	0	0	0	0	0	2	f
1029	80	1024	0	0	7	0	0	0	8	f
16	47	256	2	0	0	0	0	0	10	f
17	12	512	2	0	0	0	0	0	9	f
722	227	512	2	0	0	0	0	0	9	f
733	174	32	4	0	0	0	0	0	5	f
734	175	64	12	0	0	0	0	0	6	f
966	204	512	2	0	0	0	0	0	9	f
967	205	512	2	0	0	0	0	0	9	f
968	206	512	2	0	0	0	0	0	9	f
969	207	512	2	0	0	0	0	0	9	f
970	208	512	2	0	0	0	0	0	9	f
990	212	512	2	0	0	0	0	0	9	f
1288	226	2048	7	0	0	0	0	0	11	f
795	179	32	34	0	0	0	0	0	5	f
1278	225	8	50	0	0	0	0	0	3	f
1215	66	2048	0	10	0	0	0	0	11	f
1214	65	2048	0	10	0	0	0	0	11	f
1213	63	2048	0	10	0	0	0	0	11	f
744	81	1024	13	0	0	0	0	0	8	f
76	48	8216	0	0	0	0	10	10	4	f
77	49	8216	0	0	0	0	14	14	4	f
78	50	8216	0	0	0	0	22	22	4	f
79	51	8216	0	0	0	0	31	31	4	f
80	52	8216	0	0	0	0	44	44	4	f
81	53	8216	0	0	0	0	70	70	4	f
82	48	16	0	0	0	0	6	6	4	f
83	49	16	0	0	0	0	9	9	4	f
84	50	16	0	0	0	0	14	14	4	f
85	51	16	0	0	0	0	20	20	4	f
86	52	16	0	0	0	0	28	28	4	f
87	109	16	0	0	0	0	4	5	4	f
88	111	16	0	0	0	0	8	11	4	f
89	110	16	0	0	0	0	8	12	4	f
90	111	16	0	0	0	0	15	20	4	f
91	112	16	0	0	0	0	21	29	4	f
92	113	16	0	0	0	0	30	41	4	f
93	114	16	0	0	0	0	47	64	4	f
94	116	16	0	0	0	0	5	4	4	f
95	118	16	0	0	0	0	12	8	4	f
96	119	16	0	0	0	0	16	11	4	f
97	120	16	0	0	0	0	24	18	4	f
98	121	16	0	0	0	0	38	28	4	f
606	49	16	0	0	0	0	9	9	4	f
1172	49	16	0	0	0	0	6	6	4	f
1205	49	16	0	0	0	0	5	5	4	f
1230	49	16	0	0	0	0	5	5	4	f
1236	49	16	0	0	0	0	4	4	4	f
1255	49	16	0	0	0	0	9	9	4	f
1256	49	16	0	0	0	0	9	9	4	f
1289	229	16	0	0	0	0	999	999	4	f
1080	53	16	0	0	0	39	25	25	4	f
101	123	16	0	20	0	0	7	3	4	f
102	123	16	0	20	0	0	7	3	4	f
614	123	16	0	20	0	0	35	35	4	f
597	81	1024	3	3	3	3	10	6	8	f
522	81	1024	3	3	3	3	10	6	8	f
1081	54	16	0	0	0	39	10	8	4	f
615	123	16	0	20	0	0	35	35	4	f
1128	48	16	0	0	0	39	4	4	4	f
1129	49	16	0	0	0	39	5	5	4	f
103	123	16	0	20	0	0	7	3	4	f
616	123	16	0	20	0	0	35	35	4	f
617	123	16	0	20	0	0	35	35	4	f
618	123	16	0	20	0	0	35	35	4	f
100	123	16	0	6	0	0	6	2	4	f
1216	210	16	0	6	0	0	6	2	4	f
1006	156	256	8	0	0	0	2	2	10	f
701	156	256	8	0	0	0	2	2	10	f
1217	123	16	0	6	0	0	6	2	4	f
1218	219	16	0	6	0	0	6	2	4	f
1130	50	16	0	0	0	39	8	8	4	f
1131	51	16	0	0	0	39	11	11	4	f
1132	54	16	0	0	0	39	10	8	4	f
1133	52	16	0	0	0	39	15	15	4	f
1134	53	16	0	0	0	39	25	25	4	f
700	156	256	8	0	0	0	2	2	10	f
699	156	256	8	0	0	0	2	2	10	f
698	156	256	8	0	0	0	2	2	10	f
317	81	1024	6	6	0	0	6	6	8	f
\.


--
-- Data for Name: item_wieldable_requirements; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.item_wieldable_requirements (id, skillindex, level) FROM stdin;
63	0	5
64	0	20
65	0	30
67	0	5
68	0	20
69	0	30
72	0	5
73	0	20
74	0	30
75	0	40
78	0	5
79	0	20
80	0	30
81	0	40
84	0	5
85	0	20
86	0	30
88	0	5
90	0	5
91	0	20
92	0	30
93	0	40
95	0	5
96	0	20
97	0	30
98	0	40
105	1	5
106	1	20
107	1	30
109	1	5
110	1	20
111	1	30
112	1	40
114	1	5
115	1	20
116	1	30
118	1	5
119	1	20
120	1	30
121	1	5
122	1	20
123	1	30
125	1	5
126	1	20
127	1	30
129	1	5
130	1	20
131	1	30
196	1	10
203	0	20
204	0	30
225	1	5
226	1	20
227	1	30
230	1	10
248	1	10
309	1	5
310	1	20
311	1	30
313	1	10
396	0	40
397	0	40
398	0	40
399	1	40
400	1	40
401	1	40
401	6	30
402	1	40
403	1	40
404	1	40
405	0	40
406	1	40
407	1	40
407	6	30
423	0	10
424	0	10
425	0	10
426	0	10
427	0	10
428	0	10
429	0	10
430	0	10
431	1	10
432	1	10
433	1	10
434	1	10
561	0	5
562	0	20
563	0	40
564	0	30
593	0	60
594	0	60
594	5	40
648	4	10
649	4	5
650	4	20
651	4	15
652	4	30
653	4	25
654	4	40
655	4	35
656	4	50
657	4	45
795	1	60
1024	0	5
1068	0	20
1069	0	30
1070	0	40
1077	0	5
1078	0	20
1079	0	30
1080	0	40
1081	0	10
1089	0	5
1090	0	20
1091	0	30
1092	0	40
1124	0	5
1125	0	20
1126	0	30
1127	0	40
1130	0	5
1131	0	20
1133	0	30
1134	0	40
1137	0	5
1138	0	20
1139	0	30
1140	0	40
1278	1	60
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.items (id, name, description, command, base_price, stackable, special, members) FROM stdin;
0	Iron Mace	A spiky mace		63	f	f	f
1	Iron Short Sword	A razor sharp sword		91	f	f	f
2	Iron Kite Shield	A large metal shield		238	f	f	f
3	Iron Square Shield	A medium metal shield		168	f	f	f
4	Wooden Shield	A solid wooden shield		20	f	f	f
5	Medium Iron Helmet	A medium sized helmet		84	f	f	f
6	Large Iron Helmet	A full face helmet		154	f	f	f
7	Iron Chain Mail Body	A series of connected metal rings		210	f	f	f
8	Iron Plate Mail Body	Provides excellent protection		560	f	f	f
9	Iron Plate Mail Legs	These look pretty heavy		280	f	f	f
10	Coins	Lovely money!		1	t	f	f
11	Bronze Arrows	Arrows with bronze heads		2	t	f	f
12	Iron Axe	A woodcutters axe		56	f	f	f
13	Knife	A dangerous looking knife		6	f	f	f
14	Logs	A number of wooden logs		4	f	f	f
15	Leather Armour	Better than no armour!		21	f	f	f
16	Leather Gloves	These will keep my hands warm!		6	f	f	f
17	Boots	Comfortable leather boots		6	f	f	f
18	Cabbage	Yuck I don't like cabbage	Eat	1	f	f	f
19	Egg	A nice fresh egg		4	f	f	f
20	Bones	Ew it's a pile of bones	Bury	1	f	f	f
21	Bucket	It's a wooden bucket		2	f	f	f
22	Milk	It's a bucket of milk		6	f	f	f
23	Flour	A little heap of flour		2	f	t	f
24	Amulet of GhostSpeak	It lets me talk to ghosts		35	f	t	f
25	Silverlight key 1	A key given to me by Wizard Traiborn		1	f	t	f
26	Silverlight key 2	A key given to me by Captain Rovin		1	f	t	f
27	skull	A spooky looking skull		1	f	t	f
28	Iron dagger	Short but pointy		35	f	f	f
29	grain	Some wheat heads		2	f	f	f
30	Book		read	1	f	t	f
31	Fire-Rune	One of the 4 basic elemental runes		4	t	f	f
32	Water-Rune	One of the 4 basic elemental runes		4	t	f	f
33	Air-Rune	One of the 4 basic elemental runes		4	t	f	f
34	Earth-Rune	One of the 4 basic elemental runes		4	t	f	f
35	Mind-Rune	Used for low level missile spells		3	t	f	f
36	Body-Rune	Used for curse spells		3	t	f	f
37	Life-Rune	Used for summon spells		1	t	f	t
38	Death-Rune	Used for high level missile spells		20	t	f	f
39	Needle	Used with a thread to make clothes		1	t	f	f
40	Nature-Rune	Used for alchemy spells		7	t	f	f
41	Chaos-Rune	Used for mid level missile spells		10	t	f	f
42	Law-Rune	Used for teleport spells		12	t	f	f
43	Thread	Used with a needle to make clothes		1	t	f	f
44	Holy Symbol of saradomin	This needs a string putting on it		200	f	f	f
45	Unblessed Holy Symbol	This needs blessing		200	f	f	f
46	Cosmic-Rune	Used for enchant spells		15	t	f	f
47	key	The key to get into the phoenix gang		1	f	t	f
48	key	The key to the phoenix gang's weapons store		1	f	f	f
49	scroll	An intelligence Report		5	f	f	f
50	Water	It's a bucket of water		6	f	f	f
51	Silverlight key 3	A key I found in a drain		1	f	t	f
52	Silverlight	A magic sword		50	f	f	f
53	Broken shield	Half of the shield of Arrav		1	f	f	f
54	Broken shield	Half of the shield of Arrav		1	f	f	f
55	Cadavaberries	Poisonous berries		1	f	f	f
56	message	A message from Juliet to Romeo		1	f	f	f
57	Cadava	I'm meant to give this to Juliet		1	f	t	f
58	potion	this is meant to be good for spots		1	f	f	f
59	Phoenix Crossbow	Former property of the phoenix gang		4	f	f	f
60	Crossbow	This fires crossbow bolts		70	f	f	f
61	Certificate	I can use this to claim a reward from the king		1	f	f	f
62	bronze dagger	Short but pointy		10	f	f	f
63	Steel dagger	Short but pointy		125	f	f	f
64	Mithril dagger	Short but pointy		325	f	f	f
65	Adamantite dagger	Short but pointy		800	f	f	f
66	Bronze Short Sword	A razor sharp sword		26	f	f	f
67	Steel Short Sword	A razor sharp sword		325	f	f	f
68	Mithril Short Sword	A razor sharp sword		845	f	f	f
69	Adamantite Short Sword	A razor sharp sword		2080	f	f	f
70	Bronze Long Sword	A razor sharp sword		40	f	f	f
71	Iron Long Sword	A razor sharp sword		140	f	f	f
72	Steel Long Sword	A razor sharp sword		500	f	f	f
73	Mithril Long Sword	A razor sharp sword		1300	f	f	f
74	Adamantite Long Sword	A razor sharp sword		3200	f	f	f
75	Rune long sword	A razor sharp sword		32000	f	f	f
76	Bronze 2-handed Sword	A very large sword		80	f	f	f
77	Iron 2-handed Sword	A very large sword		280	f	f	f
78	Steel 2-handed Sword	A very large sword		1000	f	f	f
79	Mithril 2-handed Sword	A very large sword		2600	f	f	f
80	Adamantite 2-handed Sword	A very large sword		6400	f	f	f
81	rune 2-handed Sword	A very large sword		64000	f	f	f
82	Bronze Scimitar	A vicious curved sword		32	f	f	f
83	Iron Scimitar	A vicious curved sword		112	f	f	f
84	Steel Scimitar	A vicious curved sword		400	f	f	f
85	Mithril Scimitar	A vicious curved sword		1040	f	f	f
86	Adamantite Scimitar	A vicious curved sword		2560	f	f	f
87	bronze Axe	A woodcutters axe		16	f	f	f
88	Steel Axe	A woodcutters axe		200	f	f	f
89	Iron battle Axe	A vicious looking axe		182	f	f	f
90	Steel battle Axe	A vicious looking axe		650	f	f	f
91	Mithril battle Axe	A vicious looking axe		1690	f	f	f
92	Adamantite battle Axe	A vicious looking axe		4160	f	f	f
93	Rune battle Axe	A vicious looking axe		41600	f	f	f
94	Bronze Mace	A spiky mace		18	f	f	f
95	Steel Mace	A spiky mace		225	f	f	f
96	Mithril Mace	A spiky mace		585	f	f	f
97	Adamantite Mace	A spiky mace		1440	f	f	f
98	Rune Mace	A spiky mace		14400	f	f	f
99	Brass key	I wonder what this is the key to		1	f	f	f
100	staff	It's a slightly magical stick		15	f	f	f
101	Staff of Air	A Magical staff		1500	f	f	f
102	Staff of water	A Magical staff		1500	f	f	f
103	Staff of earth	A Magical staff		1500	f	f	f
104	Medium Bronze Helmet	A medium sized helmet		24	f	f	f
105	Medium Steel Helmet	A medium sized helmet		300	f	f	f
106	Medium Mithril Helmet	A medium sized helmet		780	f	f	f
107	Medium Adamantite Helmet	A medium sized helmet		1920	f	f	f
108	Large Bronze Helmet	A full face helmet		44	f	f	f
109	Large Steel Helmet	A full face helmet		550	f	f	f
110	Large Mithril Helmet	A full face helmet		1430	f	f	f
111	Large Adamantite Helmet	A full face helmet		3520	f	f	f
112	Large Rune Helmet	A full face helmet		35200	f	f	f
113	Bronze Chain Mail Body	A series of connected metal rings		60	f	f	f
114	Steel Chain Mail Body	A series of connected metal rings		750	f	f	f
115	Mithril Chain Mail Body	A series of connected metal rings		1950	f	f	f
116	Adamantite Chain Mail Body	A series of connected metal rings		4800	f	f	f
117	Bronze Plate Mail Body	Provides excellent protection		160	f	f	f
118	Steel Plate Mail Body	Provides excellent protection		2000	f	f	f
119	Mithril Plate Mail Body	Provides excellent protection		5200	f	f	f
120	Adamantite Plate Mail Body	Provides excellent protection		12800	f	f	f
121	Steel Plate Mail Legs	These look pretty heavy		1000	f	f	f
122	Mithril Plate Mail Legs	These look pretty heavy		2600	f	f	f
123	Adamantite Plate Mail Legs	These look pretty heavy		6400	f	f	f
124	Bronze Square Shield	A medium metal shield		48	f	f	f
125	Steel Square Shield	A medium metal shield		600	f	f	f
126	Mithril Square Shield	A medium metal shield		1560	f	f	f
127	Adamantite Square Shield	A medium metal shield		3840	f	f	f
128	Bronze Kite Shield	A large metal shield		68	f	f	f
129	Steel Kite Shield	A large metal shield		850	f	f	f
130	Mithril Kite Shield	A large metal shield		2210	f	f	f
131	Adamantite Kite Shield	A large metal shield		5440	f	f	f
132	cookedmeat	Mmm this looks tasty	Eat	4	f	f	f
133	raw chicken	I need to cook this first		1	f	f	f
134	burntmeat	Oh dear		1	f	f	f
135	pot	This pot is empty		1	f	f	f
136	flour	There is flour in this pot		10	f	f	f
137	bread dough	Some uncooked dough		1	f	f	f
138	bread	Nice crispy bread	Eat	12	f	f	f
139	burntbread	This bread is ruined!		1	f	f	f
140	jug	This jug is empty		1	f	f	f
141	water	It's full of water		1	f	f	f
142	wine	It's full of wine	Drink	1	f	f	f
143	grapes	Good grapes for wine making		1	f	f	f
144	shears	For shearing sheep		1	f	f	f
145	wool	I think this came from a sheep		1	f	f	f
146	fur	This would make warm clothing		10	f	f	f
147	cow hide	I should take this to the tannery		1	f	f	f
148	leather	It's a piece of leather		1	f	f	f
149	clay	Some hard dry clay		1	f	f	f
150	copper ore	this needs refining		3	f	f	f
151	iron ore	this needs refining		17	f	f	f
152	gold	this needs refining		150	f	f	f
153	mithril ore	this needs refining		162	f	f	f
154	adamantite ore	this needs refining		400	f	f	f
155	coal	hmm a non-renewable energy source!		45	f	f	f
156	Bronze Pickaxe	Used for mining		1	f	f	f
157	uncut diamond	this would be worth more cut		200	f	f	f
158	uncut ruby	this would be worth more cut		100	f	f	f
159	uncut emerald	this would be worth more cut		50	f	f	f
160	uncut sapphire	this would be worth more cut		25	f	f	f
161	diamond	this looks valuable		2000	f	f	f
162	ruby	this looks valuable		1000	f	f	f
163	emerald	this looks valuable		500	f	f	f
164	sapphire	this looks valuable		250	f	f	f
165	Herb	I need a closer look to identify this	Identify	1	f	f	t
166	tinderbox	useful for lighting a fire		1	f	f	f
167	chisel	good for detailed crafting		1	f	f	f
168	hammer	good for hitting things!		1	f	f	f
169	bronze bar	it's a bar of bronze		8	f	f	f
170	iron bar	it's a bar of iron		28	f	f	f
171	steel bar	it's a bar of steel		100	f	f	f
172	gold bar	this looks valuable		300	f	f	f
173	mithril bar	it's a bar of mithril		300	f	f	f
174	adamantite bar	it's a bar of adamantite		640	f	f	f
175	Pressure gauge	It looks like part of a machine		1	f	t	f
176	Fish Food	Keeps  your pet fish strong and healthy		1	f	f	f
177	Poison	This stuff looks nasty		1	f	f	f
178	Poisoned fish food	Doesn't seem very nice to the poor fishes		1	f	t	f
179	spinach roll	A home made spinach thing	Eat	1	f	f	f
180	Bad wine	Oh dear	Drink	1	f	f	f
181	Ashes	A heap of ashes		2	f	f	f
182	Apron	A mostly clean apron		2	f	f	f
183	Cape	A bright red cape		2	f	f	f
184	Wizards robe	I can do magic better in this		15	f	f	f
185	wizardshat	A silly pointed hat		2	f	f	f
186	Brass necklace	I'd prefer a gold one		30	f	f	f
187	skirt	A ladies skirt		2	f	f	f
188	Longbow	A Nice sturdy bow		80	f	f	f
189	Shortbow	Short but effective		50	f	f	f
190	Crossbow bolts	Good if you have a crossbow!		3	t	f	f
191	Apron	this will help keep my clothes clean		2	f	f	f
192	Chef's hat	What a silly hat		2	f	f	f
193	Beer	A glass of frothy ale	drink	2	f	f	f
194	skirt	A ladies skirt		2	f	f	f
195	skirt	A ladies skirt		2	f	f	f
196	Black Plate Mail Body	Provides excellent protection		3840	f	f	f
197	Staff of fire	A Magical staff		1500	f	f	f
198	Magic Staff	A Magical staff		200	f	f	f
199	wizardshat	A silly pointed hat		2	f	f	f
200	silk	It's a sheet of silk		30	f	f	f
201	flier	Get your axes from Bob's axes		1	f	f	f
202	tin ore	this needs refining		3	f	f	f
203	Mithril Axe	A powerful axe		520	f	f	f
204	Adamantite Axe	A powerful axe		1280	f	f	f
205	bronze battle Axe	A vicious looking axe		52	f	f	f
206	Bronze Plate Mail Legs	These look pretty heavy		80	f	f	f
207	Ball of wool	Spun from sheeps wool		2	f	f	f
208	Oil can	Its pretty full		3	f	t	f
209	Cape	A warm black cape		7	f	f	f
210	Kebab	A meaty Kebab	eat	3	f	f	f
211	Spade	A fairly small spade	Dig	3	f	f	f
212	Closet Key	A slightly smelly key		1	f	t	f
213	rubber tube	Its slightly charred		3	f	t	f
214	Bronze Plated Skirt	Designer leg protection		80	f	f	f
215	Iron Plated Skirt	Designer leg protection		280	f	f	f
216	Black robe	I can do magic better in this		13	f	f	f
217	stake	A very pointy stick		8	f	t	f
218	Garlic	A clove of garlic		3	f	f	f
219	Red spiders eggs	eewww		7	f	f	f
220	Limpwurt root	the root of a limpwurt plant		7	f	f	f
221	Strength Potion	4 doses of strength potion	Drink	14	f	f	f
222	Strength Potion	3 doses of strength potion	Drink	13	f	f	f
223	Strength Potion	2 doses of strength potion	Drink	13	f	f	f
224	Strength Potion	1 dose of strength potion	Drink	11	f	f	f
225	Steel Plated skirt	designer leg protection		1000	f	f	f
226	Mithril Plated skirt	Designer Leg protection		2600	f	f	f
227	Adamantite Plated skirt	Designer leg protection		6400	f	f	f
228	Cabbage	Yuck I don't like cabbage	Eat	1	f	f	f
229	Cape	A thick blue cape		32	f	f	f
230	Large Black Helmet	A full face helmet		1056	f	f	f
231	Red Bead	A small round red bead		4	f	f	f
232	Yellow Bead	A small round yellow bead		4	f	f	f
233	Black Bead	A small round black bead		4	f	f	f
234	White Bead	A small round white bead		4	f	f	f
235	Amulet of accuracy	It increases my aim		100	f	f	f
236	Redberries	Very bright red berries		3	f	f	f
237	Rope	A Coil of rope		18	f	f	f
238	Reddye	A little bottle of dye		5	f	f	f
239	Yellowdye	A little bottle of dye		5	f	f	f
240	Paste	A bottle off skin coloured paste		5	f	t	f
241	Onion	A strong smelling onion		3	f	f	f
242	Bronze key	A heavy key		1	f	t	f
243	Soft Clay	Clay that's ready to be used		2	f	f	f
244	wig	A blonde wig		2	f	t	f
245	wig	A wig made from wool		2	f	t	f
246	Half full wine jug	It's half full of wine	Drink	1	f	f	f
247	Keyprint	An imprint of a key in a lump of clay		2	f	t	f
248	Black Plate Mail Legs	These look pretty heavy		1920	f	f	f
249	banana	Mmm this looks tasty	Eat	2	f	f	f
250	pastry dough	Some uncooked dough		1	f	f	f
251	Pie dish	For making pies in		3	f	f	f
252	cooking apple	I wonder what i can make with this		1	f	f	f
253	pie shell	I need to find a filling for this pie		1	f	f	f
254	Uncooked apple pie	I need to cook this first		1	f	f	f
255	Uncooked meat pie	I need to cook this first		1	f	f	f
256	Uncooked redberry pie	I need to cook this first		1	f	f	f
257	apple pie	Mmm Apple pie	eat	30	f	f	f
258	Redberry pie	Looks tasty	eat	12	f	f	f
259	meat pie	Mighty and meaty	eat	15	f	f	f
260	burntpie	Oops	empty dish	1	f	f	f
261	Half a meat pie	Mighty and meaty	eat	10	f	f	f
262	Half a Redberry pie	Looks tasty	eat	4	f	f	f
263	Half an apple pie	Mmm Apple pie	eat	5	f	f	f
264	Portrait	It's a picture of a knight		3	f	t	f
265	Faladian Knight's sword	A razor sharp sword		200	f	t	f
266	blurite ore	What Strange stuff		3	f	t	f
267	Asgarnian Ale	A glass of frothy ale	drink	2	f	f	f
268	Wizard's Mind Bomb	It's got strange bubbles in it	drink	2	f	f	f
269	Dwarven Stout	A Pint of thick dark beer	drink	2	f	f	f
270	Eye of newt	It seems to be looking at me		3	f	f	f
271	Rat's tail	A bit of rat		3	f	t	f
272	Bluedye	A little bottle of dye		5	f	f	f
273	Goblin Armour	Armour Designed to fit Goblins		40	f	f	f
274	Goblin Armour	Armour Designed to fit Goblins		40	f	t	f
275	Goblin Armour	Armour Designed to fit Goblins		40	f	t	f
276	unstrung Longbow	I need to find a string for this		60	f	f	t
277	unstrung shortbow	I need to find a string for this		23	f	f	t
278	Unfired Pie dish	I need to put this in a pottery oven		3	f	f	f
279	unfired pot	I need to put this in a pottery oven		1	f	f	f
280	arrow shafts	I need to attach feathers to these		1	t	f	t
281	Woad Leaf	slightly bluish leaves		1	t	f	f
282	Orangedye	A little bottle of dye		5	f	f	f
283	Gold ring	A valuable ring		350	f	f	f
284	Sapphire ring	A valuable ring		900	f	f	f
285	Emerald ring	A valuable ring		1275	f	f	f
286	Ruby ring	A valuable ring		2025	f	f	f
287	Diamond ring	A valuable ring		3525	f	f	f
288	Gold necklace	I wonder if this is valuable		450	f	f	f
289	Sapphire necklace	I wonder if this is valuable		1050	f	f	f
290	Emerald necklace	I wonder if this is valuable		1425	f	f	f
291	Ruby necklace	I wonder if this is valuable		2175	f	f	f
292	Diamond necklace	I wonder if this is valuable		3675	f	f	f
293	ring mould	Used to make gold rings		5	f	f	f
294	Amulet mould	Used to make gold amulets		5	f	f	f
295	Necklace mould	Used to make gold necklaces		5	f	f	f
296	Gold Amulet	It needs a string so I can wear it		350	f	f	f
297	Sapphire Amulet	It needs a string so I can wear it		900	f	f	f
298	Emerald Amulet	It needs a string so I can wear it		1275	f	f	f
299	Ruby Amulet	It needs a string so I can make wear it		2025	f	f	f
300	Diamond Amulet	It needs a string so I can wear it		3525	f	f	f
301	Gold Amulet	I wonder if I can get this enchanted		350	f	f	f
302	Sapphire Amulet	I wonder if I can get this enchanted		900	f	f	f
303	Emerald Amulet	I wonder if I can get this enchanted		1275	f	f	f
304	Ruby Amulet	I wonder if I can get this enchanted		2025	f	f	f
305	Diamond Amulet	I wonder if I can get this enchanted		3525	f	f	f
306	superchisel	I wonder if I can get this enchanted	twiddle	3525	f	f	t
307	Mace of Zamorak	This mace gives me the creeps		4500	f	f	f
308	Bronze Plate Mail top	Armour designed for females		160	f	f	f
309	Steel Plate Mail top	Armour designed for females		2000	f	f	f
310	Mithril Plate Mail top	Armour designed for females		5200	f	f	f
311	Adamantite Plate Mail top	Armour designed for females		12800	f	f	f
312	Iron Plate Mail top	Armour designed for females		560	f	f	f
313	Black Plate Mail top	Armour designed for females		3840	f	f	f
314	Sapphire Amulet of magic	It improves my magic		900	f	f	f
315	Emerald Amulet of protection	It improves my defense		1275	f	f	f
316	Ruby Amulet of strength	It improves my damage		2025	f	f	f
317	Diamond Amulet of power	A powerful amulet		3525	f	f	f
318	Karamja Rum	A very strong spirit brewed in Karamja		30	f	t	f
319	Cheese	It's got holes in it	Eat	4	f	f	f
320	Tomato	This would make good ketchup	Eat	4	f	f	f
321	Pizza Base	I need to add some tomato next		4	f	f	f
322	Burnt Pizza	Oh dear!		1	f	f	f
323	Incomplete Pizza	I need to add some cheese next		10	f	f	f
324	Uncooked Pizza	This needs cooking		25	f	f	f
325	Plain Pizza	A cheese and tomato pizza	Eat	40	f	f	f
326	Meat Pizza	A pizza with bits of meat on it	Eat	50	f	f	f
327	Anchovie Pizza	A Pizza with Anchovies	Eat	60	f	f	f
328	Half Meat Pizza	Half of this pizza has been eaten	Eat	25	f	f	f
329	Half Anchovie Pizza	Half of this pizza has been eaten	Eat	30	f	f	f
330	Cake	A plain sponge cake	Eat	50	f	f	f
331	Burnt Cake	Argh what a mess!		1	f	f	f
332	Chocolate Cake	This looks very tasty!	Eat	70	f	f	f
333	Partial Cake	Someone has eaten a big chunk of this cake	Eat	30	f	f	f
334	Partial Chocolate Cake	Someone has eaten a big chunk of this cake	Eat	50	f	f	f
335	Slice of Cake	I'd rather have a whole cake!	Eat	10	f	f	f
336	Chocolate Slice	A slice of chocolate cake	Eat	30	f	f	f
337	Chocolate Bar	It's a bar of chocolate	Eat	10	f	f	f
338	Cake Tin	Useful for baking cakes		10	f	f	f
339	Uncooked cake	Now all I need to do is cook it		20	f	f	f
340	Unfired bowl	I need to put this in a pottery oven		2	f	f	f
341	Bowl	Useful for mixing things		4	f	f	f
342	Bowl of water	It's a bowl of water		3	f	f	f
343	Incomplete stew	I need to add some meat too		4	f	f	f
344	Incomplete stew	I need to add some potato too		4	f	f	f
345	Uncooked stew	I need to cook this		10	f	f	f
346	Stew	It's a meat and potato stew	Eat	20	f	f	f
347	Burnt Stew	Eew it's horribly burnt	Empty	1	f	f	f
348	Potato	Can be used to make stew		1	f	f	f
349	Raw Shrimp	I should try cooking this		5	f	f	f
350	Shrimp	Some nicely cooked fish	Eat	5	f	f	f
351	Raw Anchovies	I should try cooking this		15	f	f	f
352	Anchovies	Some nicely cooked fish	Eat	15	f	f	f
353	Burnt fish	Oops!		1	f	f	f
354	Raw Sardine	I should try cooking this		10	f	f	f
355	Sardine	Some nicely cooked fish	Eat	10	f	f	f
356	Raw Salmon	I should try cooking this		50	f	f	f
357	Salmon	Some nicely cooked fish	Eat	50	f	f	f
358	Raw Trout	I should try cooking this		20	f	f	f
359	Trout	Some nicely cooked fish	Eat	20	f	f	f
360	Burnt fish	Oops!		1	f	f	f
361	Raw Herring	I should try cooking this		15	f	f	f
362	Herring	Some nicely cooked fish	Eat	15	f	f	f
363	Raw Pike	I should try cooking this		25	f	f	f
364	Pike	Some nicely cooked fish	Eat	25	f	f	f
365	Burnt fish	Oops!		1	f	f	f
366	Raw Tuna	I should try cooking this		100	f	f	f
367	Tuna	Wow this is a big fish	Eat	100	f	f	f
368	Burnt fish	Oops!		1	f	f	f
369	Raw Swordfish	I should try cooking this		200	f	f	f
370	Swordfish	I'd better be careful eating this!	Eat	200	f	f	f
371	Burnt Swordfish	Oops!		1	f	f	f
372	Raw Lobster	I should try cooking this		150	f	f	f
373	Lobster	This looks tricky to eat	Eat	150	f	f	f
374	Burnt Lobster	Oops!		1	f	f	f
375	Lobster Pot	Useful for catching lobsters		20	f	f	f
376	Net	Useful for catching small fish		5	f	f	f
377	Fishing Rod	Useful for catching sardine or herring		5	f	f	f
378	Fly Fishing Rod	Useful for catching salmon or trout		5	f	f	f
379	Harpoon	Useful for catching really big fish		5	f	f	f
380	Fishing Bait	For use with a fishing rod		3	t	f	f
381	Feather	Used for fly-fishing		2	t	f	f
382	Chest key	A key to One eyed Hector's chest		1	f	t	f
383	Silver	this needs refining		75	f	f	f
384	silver bar	this looks valuable		150	f	f	f
385	Holy Symbol of saradomin	This improves my prayer		300	f	f	f
386	Holy symbol mould	Used to make Holy Symbols		5	f	f	f
387	Disk of Returning	Used to get out of Thordur's blackhole	spin	12	f	f	f
388	Monks robe	I feel closer to the God's when I am wearing this		40	f	f	f
389	Monks robe	Keeps a monk's legs nice and warm		30	f	f	f
390	Red key	A painted key		1	f	t	f
391	Orange Key	A painted key		1	f	t	f
392	yellow key	A painted key		1	f	t	f
393	Blue key	A painted key		1	f	t	f
394	Magenta key	A painted key		1	f	t	f
395	black key	A painted key		1	f	t	f
396	rune dagger	Short but pointy		8000	f	f	f
397	Rune short sword	A razor sharp sword		20800	f	f	f
398	rune Scimitar	A vicious curved sword		25600	f	f	f
399	Medium Rune Helmet	A medium sized helmet		19200	f	f	f
400	Rune Chain Mail Body	A series of connected metal rings		50000	f	f	f
401	Rune Plate Mail Body	Provides excellent protection		65000	f	f	f
402	Rune Plate Mail Legs	These look pretty heavy		64000	f	f	f
403	Rune Square Shield	A medium metal shield		38400	f	f	f
404	Rune Kite Shield	A large metal shield		54400	f	f	f
405	rune Axe	A powerful axe		12800	f	f	f
406	Rune skirt	Designer leg protection		64000	f	f	f
407	Rune Plate Mail top	Armour designed for females		65000	f	f	f
408	Runite bar	it's a bar of runite		5000	f	f	f
409	runite ore	this needs refining		3200	f	f	f
410	Plank	This doesn't look very useful		1	f	f	f
411	Tile	This doesn't look very useful		1	f	f	f
412	skull	A spooky looking skull		1	f	f	f
413	Big Bones	Ew it's a pile of bones	Bury	1	f	f	f
414	Muddy key	It looks like a key to a chest		1	f	f	f
415	Map	A map showing the way to the Isle of Crandor		1	f	t	f
416	Map Piece	I need some more of the map for this to be useful		1	f	t	f
417	Map Piece	I need some more of the map for this to be useful		1	f	t	f
418	Map Piece	I need some more of the map for this to be useful		1	f	t	f
419	Nails	Nails made from steel		3	t	f	f
420	Anti dragon breath Shield	Helps prevent damage from dragons		20	f	f	f
421	Maze key	The key to the entrance of Melzar's maze		1	f	f	f
422	Pumpkin	Happy halloween	eat	30	f	f	f
423	Black dagger	Short but pointy		240	f	f	f
424	Black Short Sword	A razor sharp sword		624	f	f	f
425	Black Long Sword	A razor sharp sword		960	f	f	f
426	Black 2-handed Sword	A very large sword		1920	f	f	f
427	Black Scimitar	A vicious curved sword		768	f	f	f
428	Black Axe	A sinister looking axe		384	f	f	f
429	Black battle Axe	A vicious looking axe		1248	f	f	f
430	Black Mace	A spikey mace		432	f	f	f
431	Black Chain Mail Body	A series of connected metal rings		1440	f	f	f
432	Black Square Shield	A medium metal shield		1152	f	f	f
433	Black Kite Shield	A large metal shield		1632	f	f	f
434	Black Plated skirt	designer leg protection		1920	f	f	f
435	Herb	I need a closer look to identify this	Identify	1	f	f	t
436	Herb	I need a closer look to identify this	Identify	1	f	f	t
437	Herb	I need a closer look to identify this	Identify	1	f	f	t
438	Herb	I need a closer look to identify this	Identify	1	f	f	t
439	Herb	I need a closer look to identify this	Identify	1	f	f	t
440	Herb	I need a closer look to identify this	Identify	1	f	f	t
441	Herb	I need a closer look to identify this	Identify	1	f	f	t
442	Herb	I need a closer look to identify this	Identify	1	f	f	t
443	Herb	I need a closer look to identify this	Identify	1	f	f	t
444	Guam leaf	A herb used in attack potion making		3	f	f	t
445	Marrentill	A herb used in poison cures		5	f	f	t
446	Tarromin	A useful herb		11	f	f	t
447	Harralander	A useful herb		20	f	f	t
448	Ranarr Weed	A useful herb		25	f	f	t
449	Irit Leaf	A useful herb		40	f	f	t
450	Avantoe	A useful herb		48	f	f	t
451	Kwuarm	A powerful herb		54	f	f	t
452	Cadantine	A powerful herb		65	f	f	t
453	Dwarf Weed	A powerful herb		70	f	f	t
454	Unfinished potion	I need another ingredient to finish this Guam potion		3	f	f	t
455	Unfinished potion	I need another ingredient to finish this Marrentill potion		5	f	f	t
456	Unfinished potion	I need another ingredient to finish this Tarromin potion		11	f	f	t
457	Unfinished potion	I need another ingredient to finish this Harralander potion		20	f	f	t
458	Unfinished potion	I need another ingredient to finish this Ranarr potion		25	f	f	t
459	Unfinished potion	I need another ingredient to finish this Irit potion		40	f	f	t
460	Unfinished potion	I need another ingredient to finish this Avantoe potion		48	f	f	t
461	Unfinished potion	I need another ingredient to finish this Kwuarm potion		54	f	f	t
462	Unfinished potion	I need another ingredient to finish this Cadantine potion		65	f	f	t
463	Unfinished potion	I need another ingredient to finish this Dwarfweed potion		70	f	f	t
464	Vial	It's full of water		2	f	f	f
465	Vial	This vial is empty		2	f	f	f
466	Unicorn horn	Poor unicorn		20	f	f	t
467	Blue dragon scale	A large shiny scale		50	f	f	t
468	Pestle and mortar	I can grind things for potions in this		4	f	f	t
469	Snape grass	Strange spikey grass		10	f	f	t
470	Medium black Helmet	A medium sized helmet		576	f	f	f
471	White berries	Poisonous berries		10	f	f	t
472	Ground blue dragon scale	This stuff isn't good for you		40	f	f	t
473	Ground unicorn horn	A useful potion ingredient		20	f	f	t
474	attack Potion	3 doses of attack potion	Drink	12	f	f	t
475	attack Potion	2 doses of attack potion	Drink	9	f	f	t
476	attack Potion	1 dose of attack potion	Drink	6	f	f	t
477	stat restoration Potion	3 doses of stat restoration potion	Drink	88	f	f	t
478	stat restoration Potion	2 doses of stat restoration potion	Drink	66	f	f	t
479	stat restoration Potion	1 dose of stat restoration potion	Drink	44	f	f	t
480	defense Potion	3 doses of defense potion	Drink	120	f	f	t
481	defense Potion	2 doses of defense potion	Drink	90	f	f	t
482	defense Potion	1 dose of defense potion	Drink	60	f	f	t
483	restore prayer Potion	3 doses of restore prayer potion	Drink	152	f	f	t
484	restore prayer Potion	2 doses of restore prayer potion	Drink	114	f	f	t
485	restore prayer Potion	1 dose of restore prayer potion	Drink	76	f	f	t
486	Super attack Potion	3 doses of attack potion	Drink	180	f	f	t
487	Super attack Potion	2 doses of attack potion	Drink	135	f	f	t
488	Super attack Potion	1 dose of attack potion	Drink	90	f	f	t
489	fishing Potion	3 doses of fishing potion	Drink	200	f	f	t
490	fishing Potion	2 doses of fishing potion	Drink	150	f	f	t
491	fishing Potion	1 dose of fishing potion	Drink	100	f	f	t
492	Super strength Potion	3 doses of strength potion	Drink	220	f	f	t
493	Super strength Potion	2 doses of strength potion	Drink	165	f	f	t
494	Super strength Potion	1 dose of strength potion	Drink	110	f	f	t
495	Super defense Potion	3 doses of defense potion	Drink	264	f	f	t
496	Super defense Potion	2 doses of defense potion	Drink	198	f	f	t
497	Super defense Potion	1 dose of defense potion	Drink	132	f	f	t
498	ranging Potion	3 doses of ranging potion	Drink	288	f	f	t
499	ranging Potion	2 doses of ranging potion	Drink	216	f	f	t
500	ranging Potion	1 dose of ranging potion	Drink	144	f	f	t
501	wine of Zamorak	It's full of wine	Drink	1	f	f	f
502	raw bear meat	I need to cook this first		1	f	f	f
503	raw rat meat	I need to cook this first		1	f	f	f
504	raw beef	I need to cook this first		1	f	f	f
505	enchanted bear meat	I don't fancy eating this now		1	f	t	t
506	enchanted rat meat	I don't fancy eating this now		1	f	t	t
507	enchanted beef	I don't fancy eating this now		1	f	t	t
508	enchanted chicken meat	I don't fancy eating this now		1	f	t	t
509	Dramen Staff	A magical staff cut from the dramen tree		15	f	t	t
510	Dramen Branch	I need to make this into a staff		15	f	t	t
511	Cape	A thick Green cape		32	f	f	f
512	Cape	A thick yellow cape		32	f	f	f
513	Cape	A thick Orange cape		32	f	f	f
514	Cape	A thick purple cape		32	f	f	f
515	Greendye	A little bottle of dye		5	f	f	f
516	Purpledye	A little bottle of dye		5	f	f	f
517	Iron ore certificate	Each certificate exchangable at draynor market for 5 iron ore		10	t	f	f
518	Coal certificate	Each certificate exchangable at draynor market for 5 coal		20	t	f	f
519	Mithril ore certificate	Each certificate exchangable at draynor market for 5 mithril ore		30	t	f	f
520	silver certificate	Each certificate exchangable at draynor market for 5 silver nuggets		15	t	f	f
521	Gold certificate	Each certificate exchangable at draynor market for 5 gold nuggets		25	t	f	f
522	Dragonstone Amulet	A very powerful amulet		17625	f	f	t
523	Dragonstone	This looks very valuable		10000	f	f	t
524	Dragonstone Amulet	It needs a string so I can wear it		17625	f	f	t
525	Crystal key	A very shiny key		1	f	f	t
526	Half of a key	A very shiny key		1	f	f	t
527	Half of a key	A very shiny key		1	f	f	t
528	Iron bar certificate	Each certificate exchangable at draynor market for 5 iron bars		10	t	f	f
529	steel bar certificate	Each certificate exchangable at draynor market for 5 steel bars		20	t	f	f
530	Mithril bar certificate	Each certificate exchangable at draynor market for 5 mithril bars		30	t	f	f
531	silver bar certificate	Each certificate exchangable at draynor market for 5 silver bars		15	t	f	f
532	Gold bar certificate	Each certificate exchangable at draynor market for 5 gold bars		25	t	f	f
533	Lobster certificate	Each certificate exchangable at draynor market for 5 lobsters		10	t	f	f
534	Raw lobster certificate	Each certificate exchangable at draynor market for 5 raw lobsters		10	t	f	f
535	Swordfish certificate	Each certificate exchangable at draynor market for 5 swordfish		10	t	f	f
536	Raw swordfish certificate	Each certificate exchangable at draynor market for 5 raw swordfish		10	t	f	f
537	Diary	Property of Nora.T.Hag	read	1	f	f	t
538	Front door key	A house key		1	f	t	t
539	Ball	A child's ball		1	f	t	t
540	magnet	A very attractive magnet		3	f	t	t
541	Grey wolf fur	This would make warm clothing		50	f	f	t
542	uncut dragonstone	this would be worth more cut		1000	f	f	t
543	Dragonstone ring	A valuable ring		17625	f	f	t
544	Dragonstone necklace	I wonder if this is valuable		18375	f	f	t
545	Raw Shark	I should try cooking this		300	f	f	t
546	Shark	I'd better be careful eating this!	Eat	300	f	f	t
547	Burnt Shark	Oops!		1	f	f	t
548	Big Net	Useful for catching lots of fish		20	f	f	t
549	Casket	I hope there is treasure in it	open	50	f	f	t
550	Raw cod	I should try cooking this		25	f	f	t
551	Cod	Some nicely cooked fish	Eat	25	f	f	t
552	Raw Mackerel	I should try cooking this		17	f	f	t
553	Mackerel	Some nicely cooked fish	Eat	17	f	f	t
554	Raw Bass	I should try cooking this		120	f	f	t
555	Bass	Wow this is a big fish	Eat	120	f	f	t
556	Ice Gloves	These will keep my hands cold!		6	f	t	t
557	Firebird Feather	A red hot feather		2	f	t	t
558	Firebird Feather	This is cool enough to hold now		2	f	t	t
559	Poisoned Iron dagger	Short but pointy		35	f	f	t
560	Poisoned bronze dagger	Short but pointy		10	f	f	t
561	Poisoned Steel dagger	Short but pointy		125	f	f	t
562	Poisoned Mithril dagger	Short but pointy		325	f	f	t
563	Poisoned Rune dagger	Short but pointy		8000	f	f	t
564	Poisoned Adamantite dagger	Short but pointy		800	f	f	t
565	Poisoned Black dagger	Short but pointy		240	f	f	t
566	Cure poison Potion	3 doses of cure poison potion	Drink	288	f	f	t
567	Cure poison Potion	2 doses of cure poison potion	Drink	216	f	f	t
568	Cure poison Potion	1 dose of cure poison potion	Drink	144	f	f	t
569	Poison antidote	3 doses of anti poison potion	Drink	288	f	f	t
570	Poison antidote	2 doses of anti poison potion	Drink	216	f	f	t
571	Poison antidote	1 dose of anti poison potion	Drink	144	f	f	t
572	weapon poison	For use on daggers and arrows		144	f	f	t
573	ID Paper	ID of Hartigen the black knight		1	f	f	t
574	Poison Bronze Arrows	Venomous looking arrows		2	t	f	t
575	Christmas cracker	Use on another player to pull it		1	f	f	f
576	Party Hat	Party!!!		2	f	f	f
577	Party Hat	Party!!!		2	f	f	f
578	Party Hat	Party!!!		2	f	f	f
579	Party Hat	Party!!!		2	f	f	f
580	Party Hat	Party!!!		2	f	f	f
581	Party Hat	Party!!!		2	f	f	f
582	Miscellaneous key	I wonder what this unlocks		1	f	f	t
583	Bunch of keys	Some keys on a keyring		2	f	f	t
584	Whisky	A bottle of Draynor Malt	drink	5	f	f	t
585	Candlestick	A valuable candlestick		5	f	f	t
586	Master thief armband	This denotes a great act of thievery		2	f	t	t
587	Blamish snail slime	Yuck		5	f	t	t
588	Blamish oil	made from the finest snail slime		10	f	t	t
589	Oily Fishing Rod	A rod covered in Blamish oil		15	f	t	t
590	lava eel	Strange it looks cooler now it's been cooked	eat	150	f	t	t
591	Raw lava eel	A very strange eel		150	f	t	t
592	Poison Crossbow bolts	Good if you have a crossbow!		3	t	f	t
593	Dragon sword	A Razor sharp sword		100000	f	f	t
594	Dragon axe	A vicious looking axe		200000	f	f	t
595	Jail keys	Keys to the black knight jail		2	f	t	t
596	Dusty Key	A key given to me by Velrak		1	f	t	t
597	Charged Dragonstone Amulet	A very powerful amulet	rub	17625	f	f	t
598	Grog	A murky glass of some sort of drink	drink	3	f	f	t
599	Candle	An unlit candle		3	f	t	t
600	black Candle	A spooky but unlit candle		3	f	t	t
601	Candle	A small slowly burning candle		3	f	t	t
602	black Candle	A spooky candle		3	f	t	t
603	insect repellant	Drives away all known 6 legged creatures		3	f	t	t
604	Bat bones	Ew it's a pile of bones	Bury	1	f	f	t
605	wax Bucket	It's a wooden bucket		2	f	t	t
606	Excalibur	This used to belong to king Arthur		200	f	t	t
607	Druids robe	I feel closer to the Gods when I am wearing this		40	f	f	t
608	Druids robe	Keeps a druids's knees nice and warm		30	f	f	t
609	Eye patch	It makes me look very piratical		2	f	f	t
610	Unenchanted Dragonstone Amulet	I wonder if I can get this enchanted		17625	f	f	t
611	Unpowered orb	I'd prefer it if it was powered		100	f	f	t
612	Fire orb	A magic glowing orb		300	f	f	t
613	Water orb	A magic glowing orb		300	f	f	t
614	Battlestaff	It's a slightly magical stick		7000	f	f	t
615	Battlestaff of fire	A Magical staff		15500	f	f	t
616	Battlestaff of water	A Magical staff		15500	f	f	t
617	Battlestaff of air	A Magical staff		15500	f	f	t
618	Battlestaff of earth	A Magical staff		15500	f	f	t
619	Blood-Rune	Used for high level missile spells		25	t	f	t
620	Beer glass	I need to fill this with beer		2	f	f	f
621	glassblowing pipe	Use on molten glass to make things		2	f	f	t
622	seaweed	slightly damp seaweed		2	f	f	t
623	molten glass	hot glass ready to be blown		2	f	f	t
624	soda ash	one of the ingredients for making glass		2	f	f	t
625	sand	one of the ingredients for making glass		2	f	f	t
626	air orb	A magic glowing orb		300	f	f	t
627	earth orb	A magic glowing orb		300	f	f	t
628	bass certificate	Each certificate exchangable at Catherby for 5 bass		10	t	f	t
629	Raw bass certificate	Each certificate exchangable at Catherby for 5 raw bass		10	t	f	t
630	shark certificate	Each certificate exchangable at Catherby for 5 shark		10	t	f	t
631	Raw shark certificate	Each certificate exchangable at Catherby for 5 raw shark		10	t	f	t
632	Oak Logs	Logs cut from an oak tree		20	f	f	t
633	Willow Logs	Logs cut from a willow tree		40	f	f	t
634	Maple Logs	Logs cut from a maple tree		80	f	f	t
635	Yew Logs	Logs cut from a yew tree		160	f	f	t
636	Magic Logs	Logs made from magical wood		320	f	f	t
637	Headless Arrows	I need to attach arrow heads to these		1	t	f	t
638	Iron Arrows	Arrows with iron heads		6	t	f	t
639	Poison Iron Arrows	Venomous looking arrows		6	t	f	t
640	Steel Arrows	Arrows with steel heads		24	t	f	t
641	Poison Steel Arrows	Venomous looking arrows		24	t	f	t
642	Mithril Arrows	Arrows with mithril heads		64	t	f	t
643	Poison Mithril Arrows	Venomous looking arrows		64	t	f	t
644	Adamantite Arrows	Arrows with adamantite heads		160	t	f	t
645	Poison Adamantite Arrows	Venomous looking arrows		160	t	f	t
646	Rune Arrows	Arrows with rune heads		800	t	f	t
647	Poison Rune Arrows	Venomous looking arrows		800	t	f	t
648	Oak Longbow	A Nice sturdy bow		160	f	f	t
649	Oak Shortbow	Short but effective		100	f	f	t
650	Willow Longbow	A Nice sturdy bow		320	f	f	t
651	Willow Shortbow	Short but effective		200	f	f	t
652	Maple Longbow	A Nice sturdy bow		640	f	f	t
653	Maple Shortbow	Short but effective		400	f	f	t
654	Yew Longbow	A Nice sturdy bow		1280	f	f	t
655	Yew Shortbow	Short but effective		800	f	f	t
656	Magic Longbow	A Nice sturdy bow		2560	f	f	t
657	Magic Shortbow	Short but effective		1600	f	f	t
658	unstrung Oak Longbow	I need to find a string for this		80	f	f	t
659	unstrung Oak Shortbow	I need to find a string for this		50	f	f	t
660	unstrung Willow Longbow	I need to find a string for this		160	f	f	t
661	unstrung Willow Shortbow	I need to find a string for this		100	f	f	t
662	unstrung Maple Longbow	I need to find a string for this		320	f	f	t
663	unstrung Maple Shortbow	I need to find a string for this		200	f	f	t
664	unstrung Yew Longbow	I need to find a string for this		640	f	f	t
665	unstrung Yew Shortbow	I need to find a string for this		400	f	f	t
666	unstrung Magic Longbow	I need to find a string for this		1280	f	f	t
667	unstrung Magic Shortbow	I need to find a string for this		800	f	f	t
668	barcrawl card	The official Alfred Grimhand barcrawl	read	10	f	t	t
669	bronze arrow heads	Not much use without the rest of the arrow!		1	t	f	t
670	iron arrow heads	Not much use without the rest of the arrow!		3	t	f	t
671	steel arrow heads	Not much use without the rest of the arrow!		12	t	f	t
672	mithril arrow heads	Not much use without the rest of the arrow!		32	t	f	t
673	adamantite arrow heads	Not much use without the rest of the arrow!		80	t	f	t
674	rune arrow heads	Not much use without the rest of the arrow!		400	t	f	t
675	flax	I should use this with a spinning wheel		5	f	f	t
676	bow string	I need a bow handle to attach this too		10	f	f	t
677	Easter egg	Happy Easter	eat	10	f	f	f
678	scorpion cage	I need to catch some scorpions in this		10	f	t	t
679	scorpion cage	It has 1 scorpion in it		10	f	t	t
680	scorpion cage	It has 2 scorpions in it		10	f	t	t
681	scorpion cage	It has 3 scorpions in it		10	f	t	t
682	Enchanted Battlestaff of fire	A Magical staff		42500	f	f	t
683	Enchanted Battlestaff of water	A Magical staff		42500	f	f	t
684	Enchanted Battlestaff of air	A Magical staff		42500	f	f	t
685	Enchanted Battlestaff of earth	A Magical staff		42500	f	f	t
686	scorpion cage	It has 1 scorpion in it		10	f	t	t
687	scorpion cage	It has 1 scorpion in it		10	f	t	t
688	scorpion cage	It has 2 scorpions in it		10	f	t	t
689	scorpion cage	It has 2 scorpions in it		10	f	t	t
690	gold	this needs refining		150	f	t	t
691	gold bar	this looks valuable		300	f	t	t
692	Ruby ring	A valuable ring		2025	f	t	t
693	Ruby necklace	I wonder if this is valuable		2175	f	t	t
694	Family crest	The crest of a varrocian noble family		10	f	t	t
695	Crest fragment	Part of the Fitzharmon family crest		10	f	t	t
696	Crest fragment	Part of the Fitzharmon family crest		10	f	t	t
697	Crest fragment	Part of the Fitzharmon family crest		10	f	t	t
698	Steel gauntlets	Very handy armour		6	f	t	t
699	gauntlets of goldsmithing	metal gloves for gold making		6	f	t	t
700	gauntlets of cooking	Used for cooking fish		6	f	t	t
701	gauntlets of chaos	improves bolt spells		6	f	t	t
702	robe of Zamorak	A robe worn by worshippers of Zamorak		40	f	f	t
703	robe of Zamorak	A robe worn by worshippers of Zamorak		30	f	f	t
704	Address Label	To lord Handelmort- Handelmort mansion		10	f	t	t
705	Tribal totem	It represents some sort of tribal god		10	f	t	t
706	tourist guide	Your definitive guide to Ardougne	read	1	f	f	t
707	spice	Put it in uncooked stew to make curry		230	f	f	t
708	Uncooked curry	I need to cook this		10	f	f	t
709	curry	It's a spicey hot curry	Eat	20	f	f	t
710	Burnt curry	Eew it's horribly burnt	Empty	1	f	f	t
711	yew logs certificate	Each certificate exchangable at Ardougne for 5 yew logs		10	t	f	t
712	maple logs certificate	Each certificate exchangable at Ardougne for 5 maple logs		20	t	f	t
713	willow logs certificate	Each certificate exchangable at Ardougne for 5 willow logs		30	t	f	t
714	lockpick	It makes picking some locks easier		20	f	f	t
715	Red vine worms	Strange little red worms		3	t	t	t
716	Blanket	A child's blanket		5	f	t	t
717	Raw giant carp	I should try cooking this		50	f	t	t
718	giant Carp	Some nicely cooked fish	Eat	50	f	t	t
719	Fishing competition Pass	Admits one to the Hemenster fishing competition		10	f	t	t
720	Hemenster fishing trophy	Hurrah you won a fishing competition		20	f	t	t
721	Pendant of Lucien	Gets me through the chamber of fear		12	f	t	t
722	Boots of lightfootedness	Wearing these makes me feel like I am floating		6	f	t	t
723	Ice Arrows	Can only be fired with yew or magic bows		2	t	t	t
724	Lever	This was once attached to something		20	f	t	t
725	Staff of Armadyl	A Magical staff		15	f	t	t
726	Pendant of Armadyl	Allows me to fight Lucien		12	f	t	t
727	Large cog	 A large old cog		10	f	t	t
728	Large cog	 A large old cog		10	f	t	t
729	Large cog	 A large old cog		10	f	t	t
730	Large cog	 A large old cog		10	f	t	t
731	Rat Poison	This stuff looks nasty		1	f	f	t
732	shiny Key	Quite a small key		1	f	t	t
733	khazard Helmet	A medium sized helmet		10	f	t	t
734	khazard chainmail	A series of connected metal rings		10	f	t	t
735	khali brew	A bottle of khazard's worst brew	drink	5	f	f	t
736	khazard cell keys	Keys for General Khazard's cells		1	f	t	t
737	Poison chalice	A strange looking drink	drink	20	f	t	t
738	magic whistle	A small tin whistle	blow	10	f	t	t
739	Cup of tea	A nice cup of tea	drink	10	f	f	t
740	orb of protection	a strange glowing green orb		1	f	t	t
741	orbs of protection	two strange glowing green orbs		1	f	t	t
742	Holy table napkin	a cloth given to me by sir Galahad		10	f	t	t
743	bell	I wonder what happens when i ring it	ring	1	f	t	t
744	Gnome Emerald Amulet of protection	It improves my defense		0	f	t	t
745	magic golden feather	It will point the way for me	blow on	2	f	t	t
746	Holy grail	A holy and powerful artifact		1	f	t	t
747	Script of Hazeel	An old scroll with strange ancient text		1	f	t	t
748	Pineapple	It can be cut up with a knife		1	f	f	t
749	Pineapple ring	Exotic fruit	eat	1	f	f	t
750	Pineapple Pizza	A tropicana pizza	Eat	100	f	f	t
751	Half pineapple Pizza	Half of this pizza has been eaten	Eat	50	f	f	t
752	Magic scroll	Maybe I should read it	read	1	f	t	t
753	Mark of Hazeel	A large metal amulet		0	f	t	t
754	bloody axe of zamorak	A vicious looking axe		5000	f	t	t
755	carnillean armour	the carnillean family armour		65	f	t	t
756	Carnillean Key	An old rusty key		1	f	t	t
757	Cattle prod	An old cattle prod		15	f	t	t
758	Plagued sheep remains	These sheep remains are infected		0	f	t	t
759	Poisoned animal feed	This looks nasty		0	f	t	t
760	Protective jacket	A thick heavy leather top		50	f	t	t
761	Protective trousers	A thick pair of leather trousers		50	f	t	t
762	Plagued sheep remains	These sheep remains are infected		0	f	t	t
763	Plagued sheep remains	These sheep remains are infected		0	f	t	t
764	Plagued sheep remains	These sheep remains are infected		0	f	t	t
765	dwellberries	some rather pretty blue berries	eat	4	f	f	t
766	Gasmask	Stops me breathing nasty stuff		2	f	t	t
767	picture	A picture of a lady called Elena		2	f	t	t
768	Book	Turnip growing for beginners	read	1	f	t	t
769	Seaslug	a rather nasty looking crustacean		4	f	t	t
770	chocolaty milk	Milk with chocolate in it	drink	2	f	t	t
771	Hangover cure	It doesn't look very tasty		2	f	t	t
772	Chocolate dust	I prefer it in a bar shape		2	f	f	t
773	Torch	A unlit home made torch		4	f	t	t
774	Torch	A lit home made torch		4	f	t	t
775	warrant	A search warrant for a house in Ardougne		5	f	t	t
776	Damp sticks	Some damp wooden sticks		0	f	t	t
777	Dry sticks	Some dry wooden sticks	rub together	0	f	t	t
778	Broken glass	Glass from a broken window pane		0	f	t	t
779	oyster pearls	I could work wonders with these and a chisel		1400	f	f	t
780	little key	Quite a small key		1	f	t	t
781	Scruffy note	It seems to say hongorer lure	read	2	f	f	t
782	Glarial's amulet	A bright green gem set in a necklace		1	f	t	t
783	Swamp tar	A foul smelling thick tar like substance		1	t	f	t
784	Uncooked Swamp paste	A thick tar like substance mixed with flour		1	t	f	t
785	Swamp paste	A tar like substance mixed with flour and warmed		30	t	f	t
786	Oyster pearl bolts	Great if you have a crossbow!		110	t	f	t
787	Glarials pebble	A small pebble with elven inscription		1	f	t	t
788	book on baxtorian	A book on elven history in north runescape	read	2	f	t	t
789	large key	I wonder what this is the key to		1	f	t	t
790	Oyster pearl bolt tips	Can be used to improve crossbow bolts		56	t	f	t
791	oyster	It's empty		5	f	f	t
792	oyster pearls	I could work wonders with these and a chisel		112	f	f	t
793	oyster	It's a rare oyster	open	200	f	f	t
794	Soil	It's a bucket of fine soil		2	f	t	t
795	Dragon medium Helmet	A medium sized helmet		100000	f	f	t
796	Mithril seed	Magical seeds in a mithril case	open	200	t	t	t
797	An old key	A door key		1	f	t	t
798	pigeon cage	It's for holding pigeons		1	f	t	t
799	Messenger pigeons	some very plump birds	release	1	f	t	t
800	Bird feed	A selection of mixed seeds		1	f	t	t
801	Rotten apples	Yuck!	eat	1	f	t	t
802	Doctors gown	I do feel clever wearing this		40	f	t	t
803	Bronze key	A heavy key		1	f	t	t
804	Distillator	It's for seperating compounds		1	f	t	t
805	Glarial's urn	An urn containing glarials ashes		1	f	t	f
806	Glarial's urn	An empty metal urn		1	f	t	f
807	Priest robe	I feel closer to saradomin in this		5	f	f	f
808	Priest gown	I feel closer to saradomin in this		5	f	f	f
809	Liquid Honey	This isn't worth much		0	f	t	t
810	Ethenea	An expensive colourless liquid		10	f	t	t
811	Sulphuric Broline	it's highly poisonous		1	f	t	t
812	Plague sample	An air tight tin container		1	f	t	t
813	Touch paper	For scientific testing		1	f	t	t
814	Dragon Bones	Ew it's a pile of bones	Bury	1	f	f	t
815	Herb	I need a closer look to identify this	Identify	1	f	t	t
816	Snake Weed	A very rare jungle herb		5	f	t	t
817	Herb	I need a closer look to identify this	Identify	1	f	t	t
818	Ardrigal	An interesting		5	f	t	t
819	Herb	I need a closer look to identify this	Identify	1	f	t	t
820	Sito Foil	An rare species of jungle herb		5	f	t	t
821	Herb	I need a closer look to identify this	Identify	1	f	t	t
822	Volencia Moss	A very rare species of jungle herb		5	f	t	t
823	Herb	I need a closer look to identify this	Identify	1	f	t	t
824	Rogues Purse	 A rare species of jungle herb		5	f	t	t
825	Soul-Rune	Used for high level curse spells		2500	t	f	t
826	king lathas Amulet	The amulet is red		10	f	t	t
827	Bronze Spear	A bronze tipped spear		4	f	f	t
828	halloween mask	aaaarrrghhh ... i'm a monster		15	f	f	f
829	Dragon bitter	A glass of frothy ale	drink	2	f	f	t
830	Greenmans ale	A glass of frothy ale	drink	2	f	f	t
831	halloween mask	aaaarrrghhh ... i'm a monster		15	f	f	f
832	halloween mask	aaaarrrghhh ... i'm a monster		15	f	f	f
833	cocktail glass	For sipping cocktails		0	f	f	t
834	cocktail shaker	For mixing cocktails	pour	2	f	f	t
835	Bone Key	A key delicately carved key made from a single piece of bone	Look	1	f	t	t
836	gnome robe	A high fashion robe		180	f	f	t
837	gnome robe	A high fashion robe		180	f	f	t
838	gnome robe	A high fashion robe		180	f	f	t
839	gnome robe	A high fashion robe		180	f	f	t
840	gnome robe	A high fashion robe		180	f	f	t
841	gnomeshat	A silly pointed hat		160	f	f	t
842	gnomeshat	A silly pointed hat		160	f	f	t
843	gnomeshat	A silly pointed hat		160	f	f	t
844	gnomeshat	A silly pointed hat		160	f	f	t
845	gnomeshat	A silly pointed hat		160	f	f	t
846	gnome top	rometti - the ultimate in gnome design		180	f	f	t
847	gnome top	rometti - the only name in gnome fashion!		180	f	f	t
848	gnome top	rometti - the only name in gnome fashion!		180	f	f	t
849	gnome top	rometti - the only name in gnome fashion!		180	f	f	t
850	gnome top	rometti - the only name in gnome fashion!		180	f	f	t
851	gnome cocktail guide	A book on tree gnome cocktails	read	2	f	f	t
852	Beads of the dead	A curious looking neck ornament		35	f	t	t
853	cocktail glass	For sipping cocktails	drink	2	f	f	t
854	cocktail glass	For sipping cocktails	drink	2	f	f	t
855	lemon	It's very fresh	eat	2	f	f	t
856	lemon slices	It's very fresh	eat	2	f	f	t
857	orange	It's very fresh	eat	2	f	f	t
858	orange slices	It's very fresh	eat	2	f	f	t
859	Diced orange	Fresh chunks of orange	eat	2	f	f	t
860	Diced lemon	Fresh chunks of lemon	eat	2	f	f	t
861	Fresh Pineapple	It can be cut up with a knife	eat	1	f	f	t
862	Pineapple chunks	Fresh chunks of pineapple	eat	1	f	f	t
863	lime	It's very fresh	eat	2	f	f	t
864	lime chunks	Fresh chunks of lime	eat	1	f	f	t
865	lime slices	It's very fresh	eat	2	f	f	t
866	fruit blast	A cool refreshing fruit mix	drink	2	f	f	t
867	odd looking cocktail	A cool refreshing mix	drink	2	f	f	t
868	Whisky	A locally brewed Malt	drink	5	f	f	t
869	vodka	A strong spirit	drink	5	f	f	t
870	gin	A strong spirit	drink	5	f	f	t
871	cream	Fresh cream	eat	2	f	f	t
872	Drunk dragon	A warm creamy alcoholic beverage	drink	2	f	f	t
873	Equa leaves	Small sweet smelling leaves	eat	2	f	f	t
874	SGG	A short green guy..looks good	drink	2	f	f	t
875	Chocolate saturday	A warm creamy alcoholic beverage	drink	2	f	f	t
876	brandy	A strong spirit	drink	5	f	f	t
877	blurberry special	Looks good..smells strong	drink	2	f	f	t
878	wizard blizzard	Looks like a strange mix	drink	2	f	f	t
879	pineapple punch	A fresh healthy fruit mix	drink	2	f	f	t
880	gnomebatta dough	Dough formed into a base		2	f	f	t
881	gianne dough	It's made from a secret recipe	mould	2	f	f	t
882	gnomebowl dough	Dough formed into a bowl shape		2	f	f	t
883	gnomecrunchie dough	Dough formed into cookie shapes		2	f	f	t
884	gnomebatta	A baked dough base		2	f	f	t
885	gnomebowl	A baked dough bowl	eat	2	f	f	t
886	gnomebatta	It's burnt to a sinder		2	f	f	t
887	gnomecrunchie	They're burnt to a sinder		2	f	f	t
888	gnomebowl	It's burnt to a sinder		2	f	f	t
889	Uncut Red Topaz	A semi precious stone		40	f	f	t
890	Uncut Jade	A semi precious stone		30	f	f	t
891	Uncut Opal	A semi precious stone		20	f	f	t
892	Red Topaz	A semi precious stone		200	f	f	t
893	Jade	A semi precious stone		150	f	f	t
894	Opal	A semi precious stone		100	f	f	t
895	Swamp Toad	Slippery little blighters	remove legs	2	f	f	t
896	Toad legs	Gnome delicacy apparently	eat	2	f	f	t
897	King worm	Gnome delicacy apparently	eat	2	f	f	t
898	Gnome spice	Aluft Giannes secret reciepe		2	f	f	t
899	gianne cook book	Aluft Giannes favorite dishes	read	2	f	f	t
900	gnomecrunchie	yum ... smells good	eat	2	f	f	t
901	cheese and tomato batta	Smells really good	eat	2	f	f	t
902	toad batta	actually smells quite good	eat	2	f	f	t
903	gnome batta	smells like pants	eat	2	f	f	t
904	worm batta	actually smells quite good	eat	2	f	f	t
905	fruit batta	actually smells quite good	eat	2	f	f	t
906	Veg batta	well..it looks healthy	eat	2	f	f	t
907	Chocolate bomb	Looks great	eat	2	f	f	t
908	Vegball	Looks pretty healthy	eat	2	f	f	t
909	worm hole	actually smells quite good	eat	2	f	f	t
910	Tangled toads legs	actually smells quite good	eat	2	f	f	t
911	Choc crunchies	yum ... smells good	eat	2	f	f	t
912	Worm crunchies	actually smells quite good	eat	2	f	f	t
913	Toad crunchies	actually smells quite good	eat	2	f	f	t
914	Spice crunchies	yum ... smells good	eat	2	f	f	t
915	Crushed Gemstone	A gemstone that has been smashed		2	f	f	t
916	Blurberry badge	an official cocktail maker		2	f	f	t
917	Gianne badge	an official gianne chef		2	f	f	t
918	tree gnome translation	Translate the old gnome tounge	read	2	f	f	t
919	Bark sample	A sample from the grand tree		2	f	t	t
920	War ship	A model of a karamja warship	play with	2	f	f	t
921	gloughs journal	Glough's private notes	read	2	f	t	t
922	invoice	A note with foreman's timber order	read	2	f	t	t
923	Ugthanki Kebab	A strange smelling Kebab made from Ugthanki meat - it doesn't look too good	eat	20	f	f	t
924	special curry	It's a spicy hot curry	Eat	20	f	f	t
925	glough's key	Glough left this at anita's		1	f	t	t
926	glough's notes	Scribbled notes and diagrams	read	2	f	t	t
927	Pebble	The pebble has an inscription		2	f	t	t
928	Pebble	The pebble has an inscription		2	f	t	t
929	Pebble	The pebble has an inscription		2	f	t	t
930	Pebble	The pebble has an inscription		2	f	t	t
931	Daconia rock	A magicaly crafted stone		40	f	t	t
932	Sinister key	You get a sense of dread from this key		1	f	f	t
933	Herb	I need a closer look to identify this	Identify	1	f	f	t
934	Torstol	A useful herb		25	f	f	t
935	Unfinished potion	I need Jangerberries to finish this Torstol potion		25	f	f	t
936	Jangerberries	They don't look very ripe	eat	1	f	f	t
937	fruit blast	A cool refreshing fruit mix	drink	30	f	f	t
938	blurberry special	Looks good..smells strong	drink	30	f	f	t
939	wizard blizzard	Looks like a strange mix	drink	30	f	f	t
940	pineapple punch	A fresh healthy fruit mix	drink	30	f	f	t
941	SGG	A short green guy..looks good	drink	30	f	f	t
942	Chocolate saturday	A warm creamy alcoholic beverage	drink	30	f	f	t
943	Drunk dragon	A warm creamy alcoholic beverage	drink	30	f	f	t
944	cheese and tomato batta	Smells really good	eat	120	f	f	t
945	toad batta	actually smells quite good	eat	120	f	f	t
946	gnome batta	smells like pants	eat	120	f	f	t
947	worm batta	actually smells quite good	eat	120	f	f	t
948	fruit batta	actually smells quite good	eat	120	f	f	t
949	Veg batta	well..it looks healthy	eat	120	f	f	t
950	Chocolate bomb	Looks great	eat	160	f	f	t
951	Vegball	Looks pretty healthy	eat	150	f	f	t
952	worm hole	actually smells quite good	eat	150	f	f	t
953	Tangled toads legs	actually smells quite good	eat	160	f	f	t
954	Choc crunchies	yum ... smells good	eat	85	f	f	t
955	Worm crunchies	actually smells quite good	eat	85	f	f	t
956	Toad crunchies	actually smells quite good	eat	85	f	f	t
957	Spice crunchies	yum ... smells good	eat	85	f	f	t
958	Stone-Plaque	A stone plaque with carved letters in it	Read	5	f	t	t
959	Tattered Scroll	An ancient tattered scroll	Read	5	f	t	t
960	Crumpled Scroll	An ancient crumpled scroll	Read	5	f	t	t
961	Bervirius Tomb Notes	Notes taken from the tomb of Bervirius	Read	5	f	t	t
962	Zadimus Corpse	The remains of Zadimus	Bury	1	f	t	t
963	Potion of Zamorak	It looks scary	drink	25	f	f	t
964	Potion of Zamorak	It looks scary	drink	25	f	f	t
965	Potion of Zamorak	It looks scary	drink	25	f	f	t
966	Boots	They're soft and silky		200	f	f	t
967	Boots	They're soft and silky		200	f	f	t
968	Boots	They're soft and silky		200	f	f	t
969	Boots	They're soft and silky		200	f	f	t
970	Boots	They're soft and silky		200	f	f	t
971	Santa's hat	It's a santa claus' hat		160	f	f	f
972	Locating Crystal	A magical crystal sphere	Activate	100	f	t	t
973	Sword Pommel	An ivory sword pommel		100	f	t	t
974	Bone Shard	A slender piece of bone	Look	1	f	t	t
975	Steel Wire	Useful for crafting items		200	f	f	t
976	Bone Beads	Beads carved out of bone		1	f	t	t
977	Rashiliya Corpse	The remains of the Zombie Queen	Bury	1	f	t	t
978	ResetCrystal	Helps reset things in game	Activate	100	f	f	t
979	Bronze Wire	Useful for crafting items		20	f	f	t
980	Present	Click to use this on a friend	open	160	f	f	f
981	Gnome Ball	Lets play	shoot	10	f	t	t
982	Papyrus	Used for making notes		9	f	f	t
983	A lump of Charcoal	a lump of cooked coal good for making marks.		45	f	f	t
984	Arrow	linen wrapped around an arrow head		10	f	t	t
985	Lit Arrow	A flamming arrow		10	t	t	t
986	Rocks	A few Large rocks		10	f	t	t
987	Paramaya Rest Ticket	Allows you to rest in the luxurius Paramaya Inn		5	f	t	t
988	Ship Ticket	Allows you passage on the 'Lady of the Waves' ship.		5	f	t	t
989	Damp cloth	It smells as if it's been doused in alcohol		10	f	t	t
990	Desert Boots	Boots made specially for the desert		20	f	f	t
991	Orb of light	The orb gives you a safe peaceful feeling		10	f	t	t
992	Orb of light	The orb gives you a safe peaceful feeling		10	f	t	t
993	Orb of light	The orb gives you a safe peaceful feeling		10	f	t	t
994	Orb of light	The orb gives you a safe peaceful feeling		10	f	t	t
995	Railing	A broken metal rod		10	f	t	t
996	Randas's journal	An old journal with several pages missing	read	1	f	t	t
997	Unicorn horn	Poor unicorn went splat!		20	f	t	t
998	Coat of Arms	A symbol of truth and all that is good		10	f	t	t
999	Coat of Arms	A symbol of truth and all that is good		10	f	t	t
1000	Staff of Iban	It's a slightly magical stick		15	f	t	t
1001	Dwarf brew	It's a bucket of home made brew		2	f	t	t
1002	Ibans Ashes	A heap of ashes		2	f	t	t
1003	Cat	She's sleeping..i think!		2	f	t	t
1004	A Doll of Iban	A strange doll made from sticks and cloth	search	2	f	t	t
1005	Old Journal	I wonder who wrote this!	read	1	f	t	t
1006	Klank's gauntlets	Heavy hand protection		6	f	t	t
1007	Iban's shadow	A dark mystical liquid		2	f	t	t
1008	Iban's conscience	The remains of a dove that died long ago		2	f	t	t
1009	Amulet of Othainian	A strange looking amulet		0	f	t	t
1010	Amulet of Doomion	A strange looking amulet		0	f	t	t
1011	Amulet of Holthion	A strange looking amulet		0	f	t	t
1012	keep key	A small prison key		1	f	t	t
1013	Bronze Throwing Dart	A deadly throwing dart with a bronze tip.		2	t	f	t
1014	Prototype Throwing Dart	A proto type of a deadly throwing dart.		70	t	t	t
1015	Iron Throwing Dart	A deadly throwing dart with an iron tip.		5	t	f	t
1016	Full Water Skin	A skinful of water		30	f	f	t
1017	Lens mould	A peculiar mould in the shape of a disc		10	f	t	t
1018	Lens	A perfectly formed glass disc		10	f	t	t
1019	Desert Robe	Cool light robe to wear in the desert		40	f	f	t
1020	Desert Shirt	A light cool shirt to wear in the desert		40	f	f	t
1021	Metal Key	A large metalic key.		1	f	t	t
1022	Slaves Robe Bottom	A dirty desert skirt		40	f	f	t
1023	Slaves Robe Top	A dirty desert shirt		40	f	f	t
1024	Steel Throwing Dart	A deadly throwing dart with a steel tip.		20	t	f	t
1025	Astrology Book	A book on Astrology in runescape	Read	2	f	t	t
1026	Unholy Symbol mould	use this with silver in a furnace		200	f	t	t
1027	Unholy Symbol of Zamorak	this needs stringing		200	f	t	t
1028	Unblessed Unholy Symbol of Zamorak	this needs blessing		200	f	t	t
1029	Unholy Symbol of Zamorak	a symbol indicating allegiance to Zamorak		200	f	t	t
1030	Shantay Desert Pass	Allows you into the desert through the Shantay pass worth 5 gold.		5	t	t	t
1031	Staff of Iban	The staff is damaged	wield	15	f	t	t
1032	Dwarf cannon base	bang	set down	200000	f	f	t
1033	Dwarf cannon stand	bang		200000	f	f	t
1034	Dwarf cannon barrels	bang		200000	f	f	t
1035	Dwarf cannon furnace	bang		200000	f	f	t
1036	Fingernails	Ugh gross!		0	f	t	t
1037	Powering crystal1	An intricately cut gemstone		0	f	t	t
1038	Mining Barrel	A roughly constructed barrel for carrying rock.		100	f	t	t
1039	Ana in a Barrel	A roughly constructed barrel with an Ana in it!	Look	100	f	t	t
1040	Stolen gold	I wish I could spend it		300	f	t	t
1041	multi cannon ball	A heavy metal spiked ball		10	t	f	t
1042	Railing	A metal railing replacement		10	f	t	t
1043	Ogre tooth	big sharp and nasty		0	f	t	t
1044	Ogre relic	A grotesque symbol of the ogres		0	f	t	t
1045	Skavid map	A map of cave locations		0	f	t	t
1046	dwarf remains	The remains of a dwarf savaged by goblins		1	f	t	t
1047	Key	A key for a chest		1	f	t	t
1048	Ogre relic part	A piece of a statue		0	f	t	t
1049	Ogre relic part	A piece of a statue		0	f	t	t
1050	Ogre relic part	A piece of a statue		0	f	t	t
1051	Ground bat bones	The ground bones of a bat		20	f	t	t
1052	Unfinished potion	I need another ingredient to finish the shaman potion		3	f	t	t
1053	Ogre potion	A strange liquid		120	f	t	t
1054	Magic ogre potion	A strange liquid that bubbles with power		120	f	t	t
1055	Tool kit	These could be handy!		120	f	t	t
1056	Nulodion's notes	Construction notes for dwarf cannon ammo	read	1	f	t	t
1057	cannon ammo mould	Used to make cannon ammo		5	f	f	t
1058	Tenti Pineapple	The most delicious in the whole of Kharid		1	f	t	t
1059	Bedobin Copy Key	A copy of a key for the captains of the mining camps chest		20	f	t	t
1060	Technical Plans	Very technical looking plans for making a thrown weapon of some sort	Read	500	f	t	t
1061	Rock cake	Yum... I think!	eat	0	f	t	t
1062	Bronze dart tips	Dangerous looking dart tips - need feathers for flight		1	t	f	t
1063	Iron dart tips	Dangerous looking dart tips - need feathers for flight		3	t	f	t
1064	Steel dart tips	Dangerous looking dart tips - need feathers for flight		9	t	f	t
1065	Mithril dart tips	Dangerous looking dart tips - need feathers for flight		25	t	f	t
1066	Adamantite dart tips	Dangerous looking dart tips - need feathers for flight		65	t	f	t
1067	Rune dart tips	Dangerous looking dart tips - need feathers for flight		350	t	f	t
1068	Mithril Throwing Dart	A deadly throwing dart with a mithril tip.		50	t	f	t
1069	Adamantite Throwing Dart	A deadly throwing dart with an adamantite tip.		130	t	f	t
1070	Rune Throwing Dart	A deadly throwing dart with a runite tip.		700	t	f	t
1145	Trowel	A small device for digging		1	f	t	t
1071	Prototype dart tip	Dangerous looking dart tip - needs feathers for flight		1	t	t	t
1072	info document	read to access variable choices	read	2	f	t	t
1073	Instruction manual	An old note book	read	1	f	t	t
1074	Unfinished potion	I need another ingredient to finish this potion		3	f	t	t
1075	Iron throwing knife	A finely balanced knife		6	f	f	t
1076	Bronze throwing knife	A finely balanced knife		2	f	f	t
1077	Steel throwing knife	A finely balanced knife		21	f	f	t
1078	Mithril throwing knife	A finely balanced knife		54	f	f	t
1079	Adamantite throwing knife	A finely balanced knife		133	f	f	t
1080	Rune throwing knife	A finely balanced knife		333	f	f	t
1081	Black throwing knife	A finely balanced knife		37	f	f	t
1082	Water Skin mostly full	A half full skin of water		27	f	f	t
1083	Water Skin mostly empty	A half empty skin of water		24	f	f	t
1084	Water Skin mouthful left	A waterskin with a mouthful of water left		18	f	f	t
1085	Empty Water Skin	A completely empty waterskin		15	f	f	t
1086	nightshade	Deadly!	eat	30	f	t	t
1087	Shaman robe	This has been left by one of the dead ogre shaman	search	40	f	t	t
1088	Iron Spear	An iron tipped spear		13	f	f	t
1089	Steel Spear	A steel tipped spear		46	f	f	t
1090	Mithril Spear	A mithril tipped spear		119	f	f	t
1091	Adamantite Spear	An adamantite tipped spear		293	f	f	t
1092	Rune Spear	A rune tipped spear		1000	f	f	t
1093	Cat	it's fluffs	Stroke	2	f	t	t
1094	Seasoned Sardine	They don't smell any better		10	f	f	t
1095	Kittens	purrr		2	f	t	t
1096	Kitten	purrr	stroke	2	f	t	t
1097	Wrought iron key	This key clears unlocks a very sturdy gate of some sort.		1	f	t	t
1098	Cell Door Key	A roughly hewn key		1	f	t	t
1099	A free Shantay Disclaimer	Very important information.	Read	1	f	f	t
1100	Doogle leaves	Small sweet smelling leaves		2	f	f	t
1101	Raw Ugthanki Meat	I need to cook this first		2	f	f	t
1102	Tasty Ugthanki Kebab	A fresh Kebab made from Ugthanki meat	eat	20	f	f	t
1103	Cooked Ugthanki Meat	Freshly cooked Ugthanki meat	Eat	5	f	f	t
1104	Uncooked Pitta Bread	I need to cook this.		4	f	f	t
1105	Pitta Bread	Mmmm I need to add some other ingredients yet.		10	f	f	t
1106	Tomato Mixture	A mixture of tomatoes in a bowl		3	f	f	t
1107	Onion Mixture	A mixture of onions in a bowl		3	f	f	t
1108	Onion and Tomato Mixture	A mixture of onions and tomatoes in a bowl		3	f	f	t
1109	Onion and Tomato and Ugthanki Mix	A mixture of onions and tomatoes and Ugthanki meat in a bowl		3	f	f	t
1110	Burnt Pitta Bread	Urgh - it's all burnt		1	f	f	t
1111	Panning tray	used for panning gold	search	1	f	t	t
1112	Panning tray	this tray contains gold nuggets	take gold	1	f	t	t
1113	Panning tray	this tray contains mud	search	1	f	t	t
1114	Rock pick	a sharp pick for cracking rocks		1	f	t	t
1115	Specimen brush	stiff brush for cleaning specimens		1	f	t	t
1116	Specimen jar	a jar for holding soil samples		1	f	t	t
1117	Rock Sample	A rock sample		1	f	t	t
1118	gold Nuggets	Real gold pieces!		1	t	t	t
1119	cat	looks like a healthy one		1	f	f	t
1120	Scrumpled piece of paper	A piece of paper with barely legible writing - looks like a recipe!	Read	10	f	f	t
1121	Digsite info	IAN ONLY	read	63	f	t	t
1122	Poisoned Bronze Throwing Dart	A venomous throwing dart with a bronze tip.		2	t	f	t
1123	Poisoned Iron Throwing Dart	A venomous throwing dart with an iron tip.		5	t	f	t
1124	Poisoned Steel Throwing Dart	A venomous throwing dart with a steel tip.		20	t	f	t
1125	Poisoned Mithril Throwing Dart	A venomous throwing dart with a mithril tip.		50	t	f	t
1126	Poisoned Adamantite Throwing Dart	A venomous throwing dart with an adamantite tip.		130	t	f	t
1127	Poisoned Rune Throwing Dart	A deadly venomous dart with a runite tip.		700	t	f	t
1128	Poisoned Bronze throwing knife	A finely balanced knife with a coating of venom		2	f	f	t
1129	Poisoned Iron throwing knife	A finely balanced knife with a coating of venom		6	f	f	t
1130	Poisoned Steel throwing knife	A finely balanced knife with a coating of venom		21	f	f	t
1131	Poisoned Mithril throwing knife	A finely balanced knife with a coating of venom		54	f	f	t
1132	Poisoned Black throwing knife	A finely balanced knife with a coating of venom		37	f	f	t
1133	Poisoned Adamantite throwing knife	A finely balanced knife with a coating of venom		133	f	f	t
1134	Poisoned Rune throwing knife	A finely balanced knife with a coating of venom		333	f	f	t
1135	Poisoned Bronze Spear	A bronze tipped spear with added venom 		4	f	f	t
1136	Poisoned Iron Spear	An iron tipped spear with added venom		13	f	f	t
1137	Poisoned Steel Spear	A steel tipped spear with added venom		46	f	f	t
1138	Poisoned Mithril Spear	A mithril tipped spear with added venom		119	f	f	t
1139	Poisoned Adamantite Spear	An adamantite tipped spear with added venom		293	f	f	t
1140	Poisoned Rune Spear	A rune tipped spear with added venom		1000	f	f	t
1141	Book of experimental chemistry	A book on experiments with volatile chemicals	read	1	f	t	t
1142	Level 1 Certificate	A Certificate of education	read	1	f	t	t
1143	Level 2 Certificate	A Certificate of education	read	1	f	t	t
1144	Level 3 Certificate	A Certificate of education	read	1	f	t	t
1146	Stamped letter of recommendation	A stamped scroll with a recommendation on it		1	f	t	t
1147	Unstamped letter of recommendation	I hereby recommend this student to undertake the Varrock City earth sciences exams		5	f	t	t
1148	Rock Sample	A rock sample		1	f	t	t
1149	Rock Sample	A rock sample		1	f	t	t
1150	Cracked rock Sample	It's been cracked open		1	f	t	t
1151	Belt buckle	been here some time		1	f	t	t
1152	Powering crystal2	An intricately cut gemstone		0	f	t	t
1153	Powering crystal3	An intricately cut gemstone		0	f	t	t
1154	Powering crystal4	An intricately cut gemstone		0	f	t	t
1155	Old boot	that's been here some time		1	f	t	t
1156	Bunny ears	Get another from the clothes shop if you die		1	f	t	f
1157	Damaged armour	that's been here some time		1	f	t	t
1158	Damaged armour	that's been here some time		1	f	t	t
1159	Rusty sword	that's been here some time		1	f	t	t
1160	Ammonium Nitrate	An acrid chemical		20	f	t	t
1161	Nitroglycerin	A strong acidic formula		2	f	t	t
1162	Old tooth	a large single tooth		0	f	t	t
1163	Radimus Scrolls	Scrolls that Radimus gave you	Read Scrolls	5	f	t	t
1164	chest key	A small key for a chest		1	f	t	t
1165	broken arrow	that's been here some time		1	f	t	t
1166	buttons	they've been here some time		1	f	t	t
1167	broken staff	that's been here some time		1	f	t	t
1168	vase	An old vase		1	f	t	t
1169	ceramic remains	some ancient pottery		1	f	t	t
1170	Broken glass	smashed glass		0	f	t	t
1171	Unidentified powder	who knows what this is for?		20	f	t	t
1172	Machette	A purpose built tool for cutting through thick jungle.		40	f	f	t
1173	Scroll	A letter written by the expert	read	5	f	t	t
1174	stone tablet	some ancient script is engraved on here	read	1	f	t	t
1175	Talisman of Zaros	an ancient item		1	f	t	t
1176	Explosive compound	A dark mystical powder		2	f	t	t
1177	Bull Roarer	A sound producing instrument - it may attract attention	Swing	1	f	t	t
1178	Mixed chemicals	A pungent mix of 2 chemicals		2	f	t	t
1179	Ground charcoal	Powdered charcoal!		20	f	f	t
1180	Mixed chemicals	A pungent mix of 3 chemicals		2	f	t	t
1181	Spell scroll	A magical scroll	read	5	f	t	t
1182	Yommi tree seed	A magical seed that grows into a Yommi tree - these need to be germinated	Inspect	200	t	t	t
1183	Totem Pole	A well crafted totem pole		500	f	t	t
1184	Dwarf cannon base	bang	set down	200000	f	f	t
1185	Dwarf cannon stand	bang		200000	f	f	t
1186	Dwarf cannon barrels	bang		200000	f	f	t
1187	Dwarf cannon furnace	bang		150000	f	f	t
1188	Golden Bowl	A specially made bowl constructed out of pure gold		1000	f	t	t
1189	Golden Bowl with pure water	A golden bowl filled with pure water		1000	f	t	t
1190	Raw Manta ray	A rare catch!		500	f	f	t
1191	Manta ray	A rare catch!	eat	500	f	f	t
1192	Raw Sea turtle	A rare catch!		500	f	f	t
1193	Sea turtle	Tasty!	eat	500	f	f	t
1194	Annas Silver Necklace	A necklace coated with silver		1	f	t	t
1195	Bobs Silver Teacup	A tea cup coated with silver		1	f	t	t
1196	Carols Silver Bottle	A little bottle coated with silver		1	f	t	t
1197	Davids Silver Book	An ornamental book coated with silver		1	f	t	t
1198	Elizabeths Silver Needle	An ornamental needle coated with silver		1	f	t	t
1199	Franks Silver Pot	A small pot coated with silver		1	f	t	t
1200	Thread	A piece of red thread discovered at the scene of the crime		1	f	t	t
1201	Thread	A piece of green thread discovered at the scene of the crime		1	f	t	t
1202	Thread	A piece of blue thread discovered at the scene of the crime		1	f	t	t
1203	Flypaper	Sticky paper for catching flies		1	f	t	t
1204	Murder Scene Pot	The pot has a sickly smell of poison mixed with wine		1	f	t	t
1205	A Silver Dagger	Dagger Found at crime scene		1	f	t	t
1206	Murderers fingerprint	An impression of the murderers fingerprint		1	f	t	t
1207	Annas fingerprint	An impression of Annas fingerprint		1	f	t	t
1208	Bobs fingerprint	An impression of Bobs fingerprint		1	f	t	t
1209	Carols fingerprint	An impression of Carols fingerprint		1	f	t	t
1210	Davids fingerprint	An impression of Davids fingerprint		1	f	t	t
1211	Elizabeths fingerprint	An impression of Elizabeths fingerprint		1	f	t	t
1212	Franks fingerprint	An impression of Franks fingerprint		1	f	t	t
1213	Zamorak Cape	A cape from the almighty zamorak		100	f	t	t
1214	Saradomin Cape	A cape from the almighty saradomin		100	f	t	t
1215	Guthix Cape	A cape from the almighty guthix		100	f	t	t
1216	Staff of zamorak	It's a stick of the gods		80000	f	t	t
1217	Staff of guthix	It's a stick of the gods		80000	f	t	t
1218	Staff of Saradomin	It's a stick of the gods		80000	f	t	t
1219	A chunk of crystal	A reddish crystal fragment - it looks like it formed a shape at one time.		2000	f	t	t
1220	A lump of crystal	A reddish crystal fragment - it looks like it formed a shape at one time.		2000	f	t	t
1221	A hunk of crystal	A reddish crystal fragment - it looks like it formed a shape at one time.		2000	f	t	t
1222	A red crystal	A heart shaped red crystal 	Inspect	2000	f	t	t
1223	Unidentified fingerprint	An impression of the murderers fingerprint		1	f	t	t
1224	Annas Silver Necklace	A silver necklace coated with flour		1	f	t	t
1225	Bobs Silver Teacup	A silver tea cup coated with flour		1	f	t	t
1226	Carols Silver Bottle	A little silver bottle coated with flour		1	f	t	t
1227	Davids Silver Book	An ornamental silver book coated with flour		1	f	t	t
1228	Elizabeths Silver Needle	An ornamental silver needle coated with flour		1	f	t	t
1229	Franks Silver Pot	A small silver pot coated with flour		1	f	t	t
1230	A Silver Dagger	Dagger Found at crime scene coated with flour		1	f	t	t
1231	A glowing red crystal	A glowing heart shaped red crystal - great magic must be present in this item		2000	f	t	t
1232	Unidentified liquid	A strong acidic formula		2	f	t	t
1233	Radimus Scrolls	Mission briefing and the completed map of Karamja - Sir Radimus will be pleased...	Read Scrolls	5	f	t	t
1234	Robe	A worn robe		15	f	t	t
1235	Armour	An unusually red armour		40	f	t	t
1236	Dagger	Short but pointy		35	f	t	t
1237	eye patch	It makes me look very piratical		2	f	t	t
1238	Booking of Binding	An ancient tome on Demonology	read	1	f	t	t
1239	Holy Water Vial	A deadly potion against evil kin	Throw	3	f	t	t
1240	Enchanted Vial	This enchanted vial is empty - but is ready for magical liquids.		200	f	t	t
1241	Scribbled notes	It looks like a page ripped from a book	Read	20	f	t	t
1242	Scrawled notes	It looks like a page ripped from a book	Read	20	f	t	t
1243	Scatched notes	It looks like a page ripped from a book	Read	20	f	t	t
1244	Shamans Tome	An ancient tome on various subjects...	read	1	f	t	t
1245	Edible seaweed	slightly damp seaweed	eat	2	f	f	t
1246	Rough Sketch of a bowl	A roughly sketched picture of a bowl made from metal	Read	5	f	t	t
1247	Burnt Manta ray	oops!		500	f	f	t
1248	Burnt Sea turtle	oops!		500	f	f	t
1249	Cut reed plant	A narrow long tube - it might be useful for something		2	f	f	t
1250	Magical Fire Pass	A pass which allows you to cross the flaming walls into the Flaming Octagon		1	f	t	t
1251	Snakes Weed Solution	Snakes weed in water - part of a potion		1	f	t	t
1252	Ardrigal Solution	Ardrigal herb in water - part of a potion		1	f	t	t
1253	Gujuo Potion	A potion to help against fear of the supernatural	Drink	1	f	t	t
1254	Germinated Yommi tree seed	A magical seed that grows into a Yommi tree - these have been germinated.	Inspect	200	t	t	t
1255	Dark Dagger	An unusual looking dagger made of dark shiny obsidian		91	f	t	t
1256	Glowing Dark Dagger	An unusual looking dagger made of dark shiny obsidian - it has an unnatural glow .		91	f	t	t
1257	Holy Force Spell	A powerful incantation - it affects spirits of the underworld	Cast	1	f	t	t
1258	Iron Pickaxe	Used for mining		140	f	f	f
1259	Steel Pickaxe	Requires level 6 mining to use		500	f	f	f
1260	Mithril Pickaxe	Requires level 21 mining to use		1300	f	f	f
1261	Adamantite Pickaxe	Requires level 31 mining to use		3200	f	f	f
1262	Rune Pickaxe	Requires level 41 mining to use		32000	f	f	f
1263	Sleeping Bag	Not as comfy as a bed but better than nothing	sleep	30	f	f	f
1264	A blue wizards hat	An ancient wizards hat.		2	f	t	t
1265	Gilded Totem Pole	A well crafted totem pole - given to you as a gift from Gujuo	Inspect	20	f	t	t
1266	Blessed Golden Bowl	A specially made bowl constructed out of pure gold - it looks magical somehow		1000	f	t	t
1267	Blessed Golden Bowl with Pure Water	A golden bowl filled with pure water - it looks magical somehow		1000	f	t	t
1268	Raw Oomlie Meat	Raw meat from the Oomlie bird		10	f	f	t
1269	Cooked Oomlie meat Parcel	Deliciously cooked Oomlie meat in a palm leaf pouch.	eat	35	f	f	t
1270	Dragon Bone Certificate	Each certificate exchangable at Yanille for 5 Dragon Bones		10	t	f	t
1271	Limpwurt Root Certificate	Each certificate exchangable at Yanille for 5 Limpwort roots		10	t	f	t
1272	Prayer Potion Certificate	Each certificate exchangable at Yanille for 5 prayer potions		10	t	f	t
1273	Super Attack Potion Certificate	Exchangable at Yanille for 5		10	t	f	t
1274	Super Defense Potion Certificate	Exchangable at Yanille for 5		10	t	f	t
1275	Super Strength Potion Certificate	Exchangable at Yanille for 5		10	t	f	t
1276	Half Dragon Square Shield	The Right Half of an ancient and powerful looking Dragon Square shield.		500000	f	f	t
1277	Half Dragon Square Shield	Left Half of an ancient and powerful looking Dragon Square shield.		110000	f	f	t
1278	Dragon Square Shield	An ancient and powerful looking Dragon Square shield.		500000	f	f	t
1279	Palm tree leaf	A thick green plam leaf - natives use this to cook meat in		5	f	f	t
1280	Raw Oomlie Meat Parcel	Oomlie meat in a palm leaf pouch - just needs to be cooked.		16	f	f	t
1281	Burnt Oomlie Meat parcel	Oomlie meat in a palm leaf pouch - it's burnt.		1	f	f	t
1282	Bailing Bucket	It's a water tight bucket	bail with 	10	f	f	t
1283	Plank	Damaged remains of the ship		1	f	f	t
1284	Arcenia root	the root of an arcenia plant		7	f	t	t
1285	display tea	A nice cup of tea - for display only		10	f	f	t
1286	Blessed Golden Bowl with plain water	A golden bowl filled with plain water	Empty	1000	f	t	t
1287	Golden Bowl with plain water	A golden bowl filled with plain water	Empty	1000	f	t	t
1288	Cape of legends	Shows I am a member of the legends guild		450	f	t	t
1289	Scythe	Get another from the clothes shop if you die	null	1	f	f	f
\.


--
-- Data for Name: npc_drops; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.npc_drops (npcid, itemid, minamount, maxamount, probability) FROM stdin;
11	10	1	200	0.10558944
57	38	1	1	0.039596036
57	41	1	2	0.058074195
60	38	1	2	0.065993406
60	41	1	4	0.10558944
65	10	50	300	0.10558944
66	10	50	1000	0.07919207
67	619	1	1	0.02639736
67	38	1	4	0.10558944
127	10	200	700	0.10558944
184	402	1	1	0.001959604
184	406	1	1	0.001959604
184	81	1	1	0.001959604
184	401	1	1	0.001959604
184	407	1	1	0.001959604
184	404	1	1	0.002759604
184	400	1	1	0.002759604
184	93	1	1	0.002960405
184	112	1	1	0.0012859604
184	399	1	1	0.04461154
184	1288	1	1	0.047515254
196	643	1	50	0.105
196	645	1	50	0.07
196	647	1	50	0.035
196	10	1	1500	0.2
201	619	10	15	0.065993406
201	594	1	1	0.00105
201	593	1	1	0.0014478153
201	402	1	1	0.003959604
201	406	1	1	0.003959604
201	81	1	1	0.003959604
201	401	1	1	0.003959604
201	407	1	1	0.003959604
201	404	1	1	0.003959604
201	400	1	1	0.003959604
201	93	1	1	0.027515255
202	619	5	10	0.065993406
202	93	1	1	0.034316566
202	112	1	1	0.003959604
202	402	1	1	0.003959604
202	81	1	1	0.003959604
202	401	1	1	0.003959604
202	404	1	1	0.003959604
202	400	1	1	0.003959604
264	619	1	3	0.065993406
264	38	1	10	0.065993406
264	10	100	200	0.065993406
290	1215	1	1	0.00055
290	1214	1	1	0.00055
290	1213	1	1	0.00055
290	402	1	1	0.003959604
290	406	1	1	0.003959604
290	81	1	1	0.003959604
290	401	1	1	0.003959604
290	407	1	1	0.003959604
290	404	1	1	0.004759604
290	400	1	1	0.004759604
290	93	1	1	0.012159604
290	112	1	1	0.012159604
290	795	1	1	0.0002478153
291	1216	1	1	0.00055
291	1217	1	1	0.00055
291	1218	1	1	0.00055
291	795	1	1	0.00063
291	1278	1	1	0.00105
291	594	1	1	0.0020794717
291	593	1	1	0.0020589435
291	619	15	25	0.05279472
291	402	1	1	0.003959604
291	81	1	1	0.003959604
291	401	1	1	0.003959604
291	404	1	1	0.003959604
291	400	1	1	0.003959604
311	619	1	3	0.05279472
311	38	1	6	0.1319868
344	10	1	5000	0.003959604
344	619	1	4	0.06
344	33	1	75	0.2
344	31	1	100	0.3
361	31	1	200	0.1
\.


--
-- Data for Name: npc_locations; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.npc_locations (id, startx, minx, maxx, starty, miny, maxy) FROM stdin;
19	597	592	597	3354	3350	3354
19	595	592	597	3354	3350	3354
19	593	592	597	3354	3350	3354
19	592	592	597	3352	3350	3354
552	414	414	414	2996	2996	2996
551	417	417	419	2994	2994	2994
550	419	419	419	2991	2991	2991
582	690	690	694	2338	2338	2338
582	703	703	703	2346	2346	2349
582	711	711	715	2338	2338	2338
582	703	703	703	2327	2327	2331
540	717	714	717	1400	1400	1403
540	717	714	717	1402	1400	1403
540	717	714	717	1402	1400	1403
540	715	714	717	1400	1400	1403
586	715	714	717	1402	1400	1403
586	716	715	717	1411	1400	1411
586	715	715	717	1409	1400	1411
591	711	711	711	1391	1388	1393
536	714	713	717	1388	1387	1389
581	714	713	717	1387	1387	1389
591	716	713	717	1387	1387	1389
581	724	722	725	1387	1386	1389
593	725	722	725	1389	1386	1389
592	722	722	725	1387	1386	1389
592	721	717	721	1388	1388	1388
537	716	714	716	1380	1378	1380
586	715	715	715	1382	1382	1385
10	117	111	121	710	704	714
43	473	471	485	520	514	524
43	475	471	485	519	514	524
43	473	471	485	517	514	524
43	519	519	521	3379	3369	3379
43	520	519	521	3371	3369	3379
786	513	512	518	3390	3383	3391
786	518	512	518	3388	3383	3391
786	513	512	518	3389	3383	3391
786	515	512	518	3388	3383	3391
786	516	512	518	3385	3383	3391
786	513	512	518	3384	3383	3391
786	516	512	518	3384	3383	3391
786	512	512	518	3383	3383	3391
786	518	512	518	3391	3383	3391
43	509	508	510	3393	3391	3396
43	509	508	510	3395	3391	3396
787	505	504	511	3371	3369	3374
787	507	504	511	3371	3369	3374
787	509	504	511	3371	3369	3374
787	509	504	511	3373	3369	3374
787	504	504	511	3374	3369	3374
144	325	323	327	490	487	492
150	257	257	257	626	626	628
152	323	321	326	447	445	449
151	323	321	326	447	445	449
22	283	280	284	184	183	188
22	283	280	284	185	183	188
22	283	280	284	186	183	188
22	283	280	284	187	183	188
104	178	174	187	206	201	211
104	182	174	187	201	196	206
23	272	258	278	299	290	307
62	140	140	150	635	635	643
89	299	287	310	609	595	615
292	382	370	390	3350	3348	3353
292	386	375	389	3344	3340	3347
292	391	383	394	3339	3330	3344
292	391	388	395	3345	3337	3347
292	394	390	398	3350	3345	3354
292	397	380	400	3343	3338	3346
294	388	384	393	3326	3323	3330
294	391	377	398	3316	3310	3320
294	395	390	398	3329	3322	3334
82	126	124	129	516	513	517
95	100	98	106	511	510	515
95	104	98	106	513	510	515
58	115	113	116	515	512	516
95	150	147	153	498	498	506
95	149	147	153	504	498	506
224	174	172	176	3527	3521	3528
224	174	172	176	3523	3521	3528
224	171	170	176	3523	3521	3525
222	164	163	169	3523	3521	3525
223	167	163	169	3524	3521	3525
466	175	173	178	3531	3530	3535
467	175	173	178	3534	3530	3535
217	140	136	154	3534	3522	3536
217	138	136	154	3527	3522	3536
217	145	136	154	3532	3522	3536
2	148	136	154	3522	3522	3536
2	151	136	154	3533	3522	3536
392	130	128	133	3546	3542	3547
221	117	117	120	3537	3534	3537
221	117	117	120	3537	3534	3537
318	584	583	589	619	619	623
338	582	582	589	618	618	623
13	580	579	582	621	620	623
13	575	573	577	618	615	620
338	572	564	578	613	608	613
8	569	564	572	620	614	620
8	569	564	572	615	614	620
334	559	557	561	615	613	617
95	552	551	554	610	609	616
95	552	551	554	613	609	616
95	553	551	554	615	609	616
95	554	551	554	614	609	616
336	578	577	583	600	598	602
337	580	577	583	601	598	602
318	580	574	588	606	603	607
43	205	202	213	105	100	114
43	206	202	213	107	100	114
43	208	202	213	108	100	114
43	209	202	213	109	100	114
43	210	202	213	105	100	114
43	210	202	213	107	100	114
43	212	202	213	116	100	114
793	441	438	442	3376	3376	3377
792	452	451	456	3376	3376	3378
61	268	263	273	2964	2961	2974
61	270	265	275	2964	2961	2974
61	268	263	273	2967	2961	2974
61	269	264	274	2969	2961	2974
61	268	263	273	2972	2961	2974
61	270	265	275	2972	2961	2974
343	284	278	290	2968	2962	2974
343	283	277	289	2964	2958	2970
343	286	280	292	2965	2959	2971
343	283	277	289	2961	2955	2967
343	281	275	287	2955	2949	2961
343	282	276	288	2958	2952	2964
343	278	272	284	2957	2951	2963
343	284	278	290	2964	2958	2970
343	266	260	272	2949	2943	2955
343	272	266	278	2953	2947	2959
343	267	261	273	2955	2949	2961
344	270	265	273	2955	2947	2956
344	270	265	273	2948	2947	2956
344	265	265	273	2951	2947	2956
243	110	108	142	128	121	136
243	121	108	142	130	121	136
243	133	108	142	126	121	136
243	140	108	142	131	121	136
136	120	113	142	108	103	119
136	127	113	142	109	103	119
136	128	113	142	112	103	119
136	132	113	142	113	103	119
136	135	113	142	115	103	119
136	137	113	142	112	103	119
342	75	70	80	104	99	109
342	76	71	81	106	101	111
342	78	73	83	105	100	110
342	79	74	84	106	101	111
48	116	113	119	503	498	505
54	102	101	103	523	522	525
503	94	94	95	521	521	522
309	84	82	85	523	521	523
501	84	82	85	525	524	526
488	84	83	86	534	532	535
130	136	133	138	524	522	527
56	137	133	138	526	522	527
35	111	107	114	550	549	554
2	107	93	109	564	563	570
2	103	93	109	568	563	570
2	97	93	109	567	563	570
2	98	93	109	563	563	570
161	92	92	94	650	647	652
162	91	89	91	649	647	652
9	113	109	113	667	660	669
1	124	122	124	669	666	670
7	135	131	137	661	659	665
198	130	129	133	1601	1601	1606
2	155	136	156	619	618	633
2	151	136	156	621	618	633
2	153	136	156	629	618	633
2	145	136	156	630	618	633
2	144	136	156	626	618	633
2	144	136	156	621	618	633
2	146	136	156	618	618	633
2	140	136	156	620	618	633
2	139	136	156	623	618	633
2	140	136	156	626	618	633
2	138	136	156	633	618	633
2	145	136	156	629	618	633
2	147	136	156	623	618	633
2	152	136	156	618	618	633
2	153	136	156	621	618	633
77	159	157	161	618	617	620
63	185	180	190	608	604	621
123	197	194	198	641	637	642
121	195	194	198	640	637	642
118	200	199	203	640	639	642
122	206	203	212	626	624	631
95	218	216	223	635	634	638
95	220	216	223	637	634	638
226	227	225	228	632	632	635
225	229	229	232	631	630	633
227	228	226	229	628	627	630
299	429	427	430	484	481	485
778	604	601	606	744	742	745
11	603	601	606	745	742	745
95	589	585	590	753	750	758
95	586	585	590	751	750	758
341	583	582	585	563	561	565
95	583	577	585	573	572	576
95	582	577	585	576	572	576
95	580	577	585	574	572	576
95	578	577	585	572	572	576
368	586	582	588	526	524	527
371	600	599	600	517	516	518
95	282	280	286	566	564	573
95	285	280	286	570	564	573
95	284	280	286	573	564	573
95	331	328	334	551	549	557
268	87	87	93	694	689	700
268	90	87	93	690	689	700
268	92	87	93	699	689	700
106	319	317	322	531	530	536
105	320	317	322	535	530	536
191	274	272	277	565	563	567
191	268	265	270	3379	3379	3380
172	84	83	85	676	673	677
72	84	83	85	674	673	677
339	371	368	372	580	578	581
160	104	102	106	520	518	521
173	51	49	52	675	673	677
95	503	498	504	449	447	453
95	499	498	504	449	447	453
95	500	498	504	451	447	453
95	502	498	504	451	447	453
318	517	515	518	459	459	461
301	524	522	525	463	462	467
318	520	519	526	451	448	454
306	524	519	526	451	448	454
318	497	495	526	457	455	463
318	511	495	526	455	455	463
282	448	448	450	492	492	493
95	438	437	443	494	491	496
95	441	437	443	491	491	496
11	433	432	435	491	490	493
310	433	432	436	483	480	485
289	427	425	429	489	487	491
250	418	416	421	487	484	488
149	273	271	274	632	630	634
131	265	262	267	629	626	632
129	275	273	276	658	655	659
167	279	276	280	647	646	650
157	279	277	280	632	630	634
528	652	651	653	536	534	539
115	327	327	328	539	537	540
141	301	300	306	578	577	580
172	345	343	346	1554	1551	1557
69	140	138	142	504	503	506
59	137	135	139	516	515	517
75	235	233	239	508	507	512
231	348	341	349	606	599	612
67	364	361	367	605	601	613
67	363	361	367	603	601	613
67	366	361	367	605	601	613
2	342	336	353	584	581	593
2	347	336	353	583	581	593
2	347	336	353	586	581	593
136	268	258	276	3370	3361	3380
136	264	258	276	3373	3361	3380
70	269	258	276	3369	3361	3380
70	264	258	276	3369	3361	3380
70	265	260	272	3364	3360	3370
70	268	263	272	3362	3351	3365
70	268	263	272	3358	3351	3365
70	268	263	272	3353	3351	3361
70	271	259	271	3350	3340	3350
65	326	324	331	546	543	548
102	317	288	337	563	521	577
102	298	288	337	525	521	577
102	329	288	337	532	521	577
102	323	288	337	553	521	577
102	298	279	309	570	558	589
102	327	299	337	512	507	524
101	309	307	312	533	529	535
142	321	316	322	549	544	550
142	320	316	322	1491	1488	1494
114	310	300	336	548	537	577
114	323	300	336	566	537	577
132	311	306	315	568	567	571
110	304	301	314	1507	1503	1510
138	319	316	320	2456	2454	2459
228	379	377	381	501	500	503
230	369	366	370	506	504	508
773	293	289	296	3329	3327	3334
143	294	289	295	3339	3337	3342
84	77	74	78	676	674	678
87	55	54	58	683	679	684
88	55	54	58	679	679	684
85	55	54	57	686	685	688
103	54	53	56	694	693	696
90	88	85	88	684	682	685
165	347	344	349	714	712	715
168	360	358	364	715	712	716
169	363	358	364	714	712	716
331	556	554	557	603	600	605
330	551	549	553	598	597	602
328	554	553	558	594	592	595
325	546	543	546	599	597	602
329	546	542	547	591	588	593
514	600	599	601	1703	1703	1705
522	459	453	462	754	750	760
617	401	399	404	851	848	854
617	402	399	404	849	848	854
617	400	399	404	853	848	854
620	418	415	420	847	846	849
186	223	222	227	441	439	443
185	225	222	227	441	439	443
174	252	249	252	467	458	468
204	362	355	369	459	455	469
200	360	355	369	500	480	513
200	375	355	369	495	480	513
200	374	358	384	493	480	513
200	370	358	384	493	480	513
200	369	358	384	497	480	513
200	379	358	384	494	480	513
200	369	358	384	494	480	513
200	374	358	384	481	480	513
29	318	308	342	450	439	465
29	321	308	342	441	439	465
29	329	308	342	448	439	465
29	329	308	342	454	439	465
29	317	308	342	456	439	465
154	327	308	342	449	439	465
154	318	308	342	452	439	465
154	314	308	342	443	439	465
154	329	308	342	443	439	465
154	332	308	342	456	439	465
154	325	308	342	458	439	465
154	313	308	342	453	439	465
153	332	308	342	456	439	465
153	318	308	342	455	439	465
153	332	308	342	448	439	465
153	322	308	342	440	439	465
153	319	308	342	446	439	465
270	213	208	218	3256	3251	3261
270	216	211	221	3243	3238	3248
270	219	214	224	3247	3242	3252
270	218	213	223	3243	3238	3248
270	214	209	219	3244	3239	3249
270	212	207	217	3248	3243	3253
270	616	615	620	553	550	555
270	617	615	620	552	550	555
270	618	615	620	552	550	555
270	618	615	620	554	550	555
312	615	610	618	3400	3396	3402
312	613	610	618	3400	3396	3402
312	611	610	618	3399	3396	3402
312	617	610	618	3398	3396	3402
312	615	610	618	3402	3396	3402
312	612	610	618	3402	3396	3402
43	595	591	596	3592	3590	3597
43	593	591	596	3596	3590	3597
43	602	600	606	3552	3552	3557
43	605	600	606	3553	3552	3557
134	291	279	295	713	704	718
46	376	374	378	3336	3334	3345
46	377	374	378	3340	3334	3345
46	376	374	378	3343	3334	3345
206	374	374	374	3333	3333	3333
206	374	374	374	3331	3331	3331
194	281	276	286	3494	3489	3499
197	287	282	292	457	452	462
43	358	353	366	3330	3321	3335
43	355	353	366	3331	3321	3335
43	355	353	366	3333	3321	3335
43	357	353	366	3333	3321	3335
43	360	353	366	3331	3321	3335
193	260	258	265	641	636	642
179	345	340	350	2518	2513	2523
187	244	239	249	443	438	448
270	345	341	348	3316	3314	3321
270	342	341	348	3317	3314	3321
270	342	341	348	3319	3314	3321
270	345	341	348	3314	3314	3321
270	344	341	348	3318	3314	3321
189	342	336	357	3354	3335	3363
189	346	336	357	3344	3335	3363
189	355	336	357	3340	3335	3363
271	344	339	356	3379	3362	3388
271	342	339	356	3364	3362	3388
271	351	339	356	3370	3362	3388
189	356	341	362	3426	3391	3436
189	351	341	362	3433	3391	3436
189	348	341	362	3400	3391	3436
189	346	341	362	3433	3391	3436
189	346	341	362	3426	3391	3436
189	350	341	362	3426	3391	3436
189	349	341	362	3418	3391	3436
189	343	341	362	3419	3391	3436
189	353	341	362	3412	3391	3436
189	343	341	362	3413	3391	3436
189	353	341	362	3418	3391	3436
272	361	359	362	3428	3428	3430
29	361	359	362	3429	3428	3430
29	357	341	362	3427	3391	3436
29	345	341	362	3426	3391	3436
29	353	341	362	3420	3391	3436
266	348	341	355	3431	3422	3436
22	349	344	354	3352	3350	3363
22	349	344	354	3354	3350	3363
22	352	344	354	3353	3350	3363
280	440	439	441	502	502	507
113	268	267	271	3329	3327	3331
519	640	635	646	753	744	767
519	640	635	646	755	744	767
519	639	635	646	756	744	767
519	639	635	646	753	744	767
519	642	635	646	754	744	767
519	642	635	646	753	744	767
518	644	635	646	754	744	767
18	138	138	142	1406	1395	1407
0	451	435	466	462	454	481
0	448	435	466	460	454	481
0	519	513	529	493	487	504
0	521	513	529	493	487	504
4	620	568	627	493	486	529
4	620	568	627	497	486	529
4	619	568	627	503	486	529
4	615	568	627	506	486	529
4	616	568	627	511	486	529
4	602	568	627	525	486	529
4	591	568	627	528	486	529
4	581	568	627	528	486	529
4	576	568	627	527	486	529
4	573	568	627	521	486	529
4	573	568	627	524	486	529
29	582	577	592	3353	3345	3358
29	580	577	592	3353	3345	3358
29	583	577	592	3351	3345	3358
19	594	587	599	3349	3345	3356
19	592	587	599	3351	3345	3356
19	592	587	599	3353	3345	3356
19	598	587	599	3351	3345	3356
29	594	587	599	3355	3345	3356
29	597	587	599	3353	3345	3356
62	605	605	621	3323	3313	3332
62	606	605	621	3323	3313	3332
62	608	605	621	3324	3313	3332
62	608	605	621	3326	3313	3332
62	616	605	621	3319	3313	3332
62	616	605	621	3318	3313	3332
62	616	605	621	3317	3313	3332
62	616	605	621	3316	3313	3332
4	620	605	621	3315	3313	3332
4	619	605	621	3316	3313	3332
4	619	605	621	3317	3313	3332
4	619	605	621	3318	3313	3332
4	619	605	621	3319	3313	3332
4	618	605	621	3319	3313	3332
4	618	605	621	3318	3313	3332
4	618	605	621	3317	3313	3332
47	600	596	605	3339	3332	3344
47	599	596	605	3340	3332	3344
47	598	596	605	3340	3332	3344
47	596	596	605	3341	3332	3344
43	596	585	596	3322	3321	3338
43	596	585	596	3324	3321	3338
43	593	585	596	3326	3321	3338
43	592	585	596	3326	3321	3338
4	584	582	621	3344	3314	3354
4	583	582	621	3343	3314	3354
4	583	582	621	3339	3314	3354
4	588	582	621	3339	3314	3354
4	592	582	621	3337	3314	3354
4	595	582	621	3334	3314	3354
4	595	582	621	3334	3314	3354
4	596	582	621	3331	3314	3354
29	597	582	621	3328	3314	3354
29	598	582	621	3328	3314	3354
29	600	582	621	3323	3314	3354
29	605	582	621	3323	3314	3354
29	608	582	621	3326	3314	3354
29	610	582	621	3325	3314	3354
29	216	209	225	626	606	631
125	223	220	225	626	623	627
34	221	209	225	627	606	631
34	221	209	225	623	606	631
34	219	209	225	623	606	631
34	214	209	225	625	606	631
97	215	214	217	618	618	621
124	214	212	219	627	622	627
29	202	195	213	619	605	624
367	584	583	593	3475	3472	3477
367	584	583	593	3473	3472	3477
420	614	613	619	610	606	612
420	617	613	619	609	606	612
420	617	613	619	610	606	612
267	443	440	446	672	670	675
262	570	568	572	587	577	591
262	570	568	572	581	577	591
262	571	568	572	584	577	591
262	571	568	572	580	577	591
262	572	568	572	579	577	591
262	570	568	572	583	577	591
407	659	641	673	644	625	644
407	657	641	673	636	625	644
428	660	641	673	630	625	644
409	650	635	670	644	625	671
409	657	635	670	644	625	671
409	666	635	670	650	625	671
409	663	635	670	655	625	671
409	646	635	670	654	625	671
409	660	635	670	667	625	671
409	650	635	670	656	625	671
409	653	635	670	652	625	671
409	645	635	670	654	625	671
409	656	635	670	650	625	671
409	657	635	670	647	625	671
409	652	635	670	646	625	671
29	665	635	670	636	625	671
29	656	635	670	630	625	671
29	653	635	670	630	625	671
29	666	635	670	630	625	671
29	661	635	670	666	625	671
29	651	635	670	669	625	671
66	277	263	296	448	435	466
66	283	263	296	451	435	466
66	276	263	296	454	435	466
100	266	264	270	438	436	445
100	267	264	270	440	436	445
100	266	264	270	442	436	445
100	272	271	274	437	435	442
100	272	271	274	440	435	442
100	273	271	275	1382	1379	1385
100	272	271	275	1384	1379	1385
66	275	275	282	438	433	442
66	282	275	282	438	433	442
66	277	275	282	441	433	442
109	281	276	282	1378	1377	1382
107	280	276	282	1379	1377	1382
108	281	276	282	1381	1377	1382
66	272	271	282	2323	2321	2335
66	277	271	282	2328	2321	2335
34	280	271	282	2331	2321	2335
66	277	271	282	2332	2321	2335
66	275	271	282	2334	2321	2335
66	273	271	282	2331	2321	2335
770	275	262	292	493	483	507
94	264	259	270	3349	3336	3361
94	262	259	270	3344	3336	3361
94	268	259	270	3344	3336	3361
94	264	259	270	3340	3336	3361
94	279	267	283	3336	3325	3341
94	279	267	283	3330	3325	3341
94	276	267	283	3329	3325	3341
94	282	276	314	3349	3334	3358
94	284	276	314	3348	3334	3358
94	286	276	314	3347	3334	3358
94	288	276	314	3345	3334	3358
94	290	276	314	3344	3334	3358
94	289	276	314	3347	3334	3358
94	287	276	314	3349	3334	3358
94	293	276	314	3348	3334	3358
93	251	249	265	461	458	472
93	261	249	265	466	458	472
93	264	249	265	462	458	472
93	263	249	265	466	458	472
93	254	249	265	469	458	472
304	264	249	265	1403	1402	1412
63	152	135	176	583	562	592
70	263	248	276	409	395	424
70	248	248	276	395	395	424
70	263	248	276	409	395	424
70	267	248	276	409	395	424
70	263	248	276	414	395	424
70	269	248	276	406	395	424
70	270	248	276	417	395	424
46	268	264	290	376	370	390
46	274	264	290	379	370	390
46	280	264	290	379	370	390
46	284	264	290	380	370	390
46	282	264	290	382	370	390
46	282	264	290	376	370	390
46	216	212	236	392	388	420
46	219	212	236	394	388	420
46	218	212	236	395	388	420
46	227	212	236	399	388	420
46	221	212	236	409	388	420
46	230	212	236	410	388	420
46	224	212	236	415	388	420
46	232	212	236	419	388	420
139	687	682	692	634	629	639
139	689	684	694	639	634	644
139	690	685	695	642	637	647
139	700	695	705	649	644	654
140	686	681	691	642	637	647
140	701	696	706	640	635	645
140	697	692	702	646	641	651
8	691	686	696	658	653	663
664	687	686	698	658	649	659
664	686	686	698	654	649	659
664	690	686	698	659	649	659
664	689	686	698	655	649	659
664	692	686	698	656	649	659
665	692	690	695	650	650	652
666	689	689	695	652	650	652
666	691	689	695	651	650	652
706	662	638	671	734	733	757
706	667	638	671	741	733	757
706	663	638	671	737	733	757
706	663	638	671	733	733	757
706	658	638	671	734	733	757
521	537	528	547	752	738	767
521	541	528	547	749	738	767
521	542	528	547	746	738	767
521	540	528	547	751	738	767
521	543	528	547	745	738	767
521	533	528	547	761	738	767
521	534	528	547	750	738	767
521	544	528	547	756	738	767
521	534	528	547	752	738	767
521	534	528	547	761	738	767
521	539	528	547	764	738	767
521	539	528	547	759	738	767
521	537	528	547	755	738	767
521	529	528	547	755	738	767
425	577	574	590	3419	3406	3423
425	583	574	590	3420	3406	3423
425	588	574	590	3411	3406	3423
425	581	574	590	3413	3406	3423
425	583	574	590	3410	3406	3423
425	576	574	590	3419	3406	3423
323	610	607	615	605	599	607
323	608	607	615	601	599	607
562	704	697	707	526	523	529
562	700	697	707	526	523	529
592	743	673	751	523	433	530
592	743	673	751	495	433	530
592	733	673	751	518	433	530
592	702	673	751	488	433	530
592	698	673	751	516	433	530
592	750	673	751	481	433	530
592	689	673	751	513	433	530
592	745	673	751	517	433	530
592	689	673	751	513	433	530
592	677	673	751	486	433	530
592	731	673	751	494	433	530
592	750	673	751	480	433	530
592	714	673	751	514	433	530
592	702	673	751	497	433	530
592	679	673	751	512	433	530
592	728	673	751	519	433	530
592	734	673	751	512	433	530
592	725	673	751	522	433	530
592	736	673	751	487	433	530
592	746	673	751	510	433	530
593	684	673	751	486	433	530
593	695	673	751	529	433	530
593	736	673	751	517	433	530
593	710	673	751	503	433	530
593	723	673	751	489	433	530
593	711	673	751	529	433	530
593	746	673	751	525	433	530
593	674	673	751	484	433	530
593	743	673	751	528	433	530
593	701	673	751	496	433	530
593	702	673	751	519	433	530
593	706	673	751	526	433	530
593	684	673	751	491	433	530
593	745	673	751	502	433	530
593	716	673	751	484	433	530
593	702	673	751	509	433	530
593	748	673	751	520	433	530
593	720	673	751	509	433	530
593	697	673	751	529	433	530
593	733	673	751	516	433	530
585	720	673	751	514	433	530
585	746	673	751	494	433	530
585	704	673	751	499	433	530
585	681	673	751	498	433	530
585	710	673	751	480	433	530
585	707	673	751	508	433	530
585	748	673	751	519	433	530
585	701	673	751	487	433	530
585	689	673	751	488	433	530
585	675	673	751	500	433	530
585	731	673	751	508	433	530
585	704	673	751	528	433	530
585	682	673	751	525	433	530
585	728	673	751	497	433	530
585	685	673	751	515	433	530
585	735	673	751	480	433	530
585	749	673	751	506	433	530
585	698	673	751	499	433	530
585	727	673	751	487	433	530
585	674	673	751	511	433	530
562	705	700	706	465	460	468
562	701	700	706	467	460	468
551	699	699	706	454	454	460
551	701	699	706	455	454	460
188	229	203	241	331	312	345
188	224	203	241	328	312	345
188	209	207	222	364	359	374
188	213	207	222	365	359	374
188	289	281	298	332	327	339
188	291	281	298	332	327	339
199	249	245	271	345	340	368
199	257	245	271	346	340	368
199	263	245	271	346	340	368
199	260	245	271	359	340	368
199	254	245	271	360	340	368
199	265	245	271	355	340	368
53	67	58	72	326	314	330
53	65	58	72	320	314	330
53	70	58	72	318	314	330
53	84	80	88	314	311	324
53	66	63	85	306	303	310
61	80	75	85	329	324	334
74	304	300	334	298	285	310
74	302	300	334	298	285	310
23	313	300	334	309	285	310
74	316	300	334	303	285	310
47	323	300	334	305	285	310
47	320	300	334	293	285	310
68	172	164	188	324	303	330
68	170	164	188	322	303	330
68	168	164	188	324	303	330
41	172	164	188	319	303	330
41	168	164	188	316	303	330
41	166	164	188	319	303	330
41	167	164	188	323	303	330
68	168	164	177	309	303	315
68	168	164	177	311	303	315
41	167	164	177	312	303	315
41	170	164	177	313	303	315
41	165	164	177	307	303	315
41	168	164	177	307	303	315
3	525	520	529	620	616	622
3	523	520	529	621	616	622
3	523	520	529	617	616	622
3	526	520	529	619	616	622
321	563	562	567	567	564	572
321	565	562	567	569	564	572
321	549	544	552	591	558	603
321	548	544	552	595	558	603
356	396	390	402	3299	3295	3302
356	393	390	402	3300	3295	3302
356	391	390	402	3297	3295	3302
356	395	390	402	3296	3295	3302
356	399	390	402	3297	3295	3302
356	402	390	402	3299	3295	3302
356	400	390	402	3300	3295	3302
356	396	390	402	3302	3295	3302
356	395	390	402	3300	3295	3302
184	82	77	87	159	154	164
184	82	77	87	162	157	167
574	638	634	639	1681	1680	1685
574	636	634	639	1681	1680	1685
574	634	634	639	1683	1680	1685
574	637	634	639	1685	1680	1685
45	348	343	353	2515	2510	2520
45	343	341	351	2516	2514	2524
53	347	343	353	1575	1568	1578
53	347	343	353	1576	1568	1578
178	348	342	352	1573	1571	1581
53	346	340	350	1570	1565	1575
181	342	340	348	3462	3460	3465
182	346	345	347	3463	3462	3467
71	70	65	75	667	662	672
194	281	276	286	3473	3468	3478
47	343	338	348	632	627	636
47	341	335	345	634	630	640
180	343	338	348	3469	3464	3474
68	344	338	348	3468	3463	3473
72	316	314	319	1607	1603	1608
148	317	315	318	667	665	669
13	57	54	62	603	599	611
13	69	63	73	640	635	646
13	66	66	81	651	645	656
13	71	66	81	647	645	656
13	76	66	81	649	645	656
13	68	66	81	645	645	656
13	56	56	83	661	655	671
13	63	56	83	663	655	671
13	75	56	83	661	655	671
13	82	56	83	659	655	671
571	88	88	88	661	661	661
72	83	58	86	680	664	682
72	79	58	86	672	664	682
72	69	58	86	673	664	682
72	62	58	86	670	664	682
72	59	58	86	677	664	682
72	62	58	86	681	664	682
29	68	58	86	680	664	682
13	83	82	86	692	688	701
322	547	544	554	586	582	608
323	551	544	554	596	582	608
323	546	544	554	604	582	608
321	564	564	572	595	592	599
321	568	564	572	597	592	599
317	467	467	467	651	646	656
261	453	447	455	680	678	681
64	448	447	455	685	682	687
264	452	447	455	683	682	687
260	453	447	455	686	682	687
332	452	447	455	687	682	688
262	458	454	458	676	672	680
262	457	454	458	680	672	680
262	455	454	458	677	672	680
257	464	459	466	682	681	685
258	469	459	469	676	672	680
258	467	459	469	678	672	680
258	465	459	469	676	672	680
259	463	459	469	677	672	680
271	494	483	508	659	644	671
271	498	483	508	653	644	671
271	502	483	508	659	644	671
271	499	483	508	667	644	671
271	506	483	508	665	644	671
594	520	511	523	661	658	670
594	518	511	523	662	658	670
594	516	511	523	664	658	670
594	515	511	523	666	658	670
421	480	479	512	705	687	717
421	487	479	512	709	687	717
421	491	479	512	703	687	717
421	497	479	512	708	687	717
421	502	479	512	703	687	717
421	507	479	512	700	687	717
421	498	479	512	713	687	717
421	492	479	512	711	687	717
255	437	436	550	695	691	697
256	439	436	550	692	691	697
64	442	441	447	692	689	698
29	446	441	447	697	689	698
278	444	440	449	706	696	709
279	451	449	454	705	702	706
264	453	445	459	698	691	709
264	450	445	459	699	691	709
264	449	445	459	702	691	709
264	453	445	459	703	691	709
264	451	445	459	701	691	709
264	449	445	459	704	691	709
70	464	459	467	717	713	717
421	466	448	466	742	738	749
521	466	449	477	732	727	740
521	457	449	477	735	727	740
67	466	451	473	763	750	771
517	436	434	439	751	748	754
523	420	416	430	764	755	767
421	416	416	430	761	755	767
521	417	416	430	763	755	767
521	419	416	430	760	755	767
521	422	416	430	759	755	767
521	426	416	430	759	755	767
521	422	416	430	766	755	767
521	420	416	430	767	755	767
521	424	416	430	766	755	767
556	425	425	425	764	764	764
521	426	425	429	771	769	773
521	428	425	429	772	769	773
521	396	393	413	764	760	771
557	401	398	405	763	763	766
523	374	370	381	774	768	787
521	374	370	381	782	768	787
521	377	370	381	771	768	787
521	385	379	385	784	777	795
421	385	379	385	788	777	795
70	379	379	385	786	777	795
523	440	427	450	833	812	840
521	436	427	450	831	812	840
521	437	427	450	825	812	840
521	433	427	450	820	812	840
523	432	427	450	816	812	840
523	440	427	450	816	812	840
521	444	427	450	821	812	840
521	442	427	450	828	812	840
521	450	427	450	822	812	840
542	382	380	390	847	845	856
542	384	380	390	854	845	856
542	386	380	390	854	845	856
542	388	380	390	854	845	856
542	386	380	390	852	845	856
542	387	380	390	849	845	856
542	388	380	390	849	845	856
542	390	380	390	845	845	856
542	388	380	390	850	845	856
542	385	380	390	851	845	856
542	387	380	390	856	845	856
542	390	380	390	855	845	856
765	397	394	402	865	863	866
765	424	422	428	863	861	864
311	521	520	524	3436	3431	3438
311	524	520	524	3435	3431	3438
312	511	509	513	3418	3415	3422
312	511	508	513	3410	3408	3415
312	510	508	513	3444	3440	3447
294	506	500	508	3428	3424	3430
36	208	204	211	561	556	569
40	210	208	211	546	545	546
38	213	208	215	2441	2438	2442
91	209	208	215	2439	2438	2442
62	294	287	298	666	655	672
62	297	287	298	662	655	672
62	295	287	298	671	655	672
47	299	288	300	676	672	679
47	296	288	300	675	672	679
164	339	337	341	706	704	708
70	424	390	434	686	665	710
70	419	390	434	701	665	710
70	412	390	434	699	665	710
70	403	390	434	697	665	710
70	406	390	434	681	665	710
114	408	390	434	687	665	710
114	414	390	434	687	665	710
114	422	390	434	689	665	710
114	419	390	434	700	665	710
114	425	390	434	707	665	710
70	421	390	434	707	665	710
70	430	390	434	700	665	710
67	414	408	430	606	603	620
67	422	408	430	608	603	620
67	429	408	430	618	603	620
41	654	649	655	3535	3531	3536
67	652	649	655	3536	3531	3536
67	652	649	655	3532	3531	3536
407	669	667	673	595	590	599
407	671	667	673	594	590	599
407	667	667	673	594	590	599
410	671	667	673	592	590	599
189	279	277	280	665	665	667
21	280	277	281	662	662	664
192	282	282	284	662	662	664
137	283	282	284	659	658	661
64	282	282	284	655	655	657
714	161	158	166	513	511	519
783	166	158	166	513	511	519
31	168	165	174	500	494	502
133	178	176	181	484	480	487
33	143	141	145	519	518	521
26	128	125	130	505	504	506
28	133	133	137	528	528	529
64	146	141	151	534	532	535
27	150	145	151	534	533	535
64	147	145	151	531	524	532
64	148	145	151	528	524	532
111	149	148	152	557	554	560
183	150	148	152	1506	1504	1507
112	149	148	152	1499	1498	1501
14	136	133	137	511	508	515
42	126	124	128	474	472	477
20	128	126	134	456	455	459
39	100	97	106	493	484	496
32	100	97	106	477	473	482
49	98	97	99	1428	1428	1431
98	83	78	85	444	442	452
44	83	78	85	453	451	454
314	84	83	85	1387	1386	1388
189	276	270	282	187	177	196
189	265	250	272	186	177	196
189	269	250	272	186	177	196
189	259	250	272	181	177	196
189	258	250	272	180	177	196
189	255	250	272	185	177	196
61	217	213	221	184	180	188
61	222	218	226	185	181	189
61	228	224	232	182	178	186
61	225	221	229	171	167	175
22	231	230	245	174	174	182
22	238	230	245	178	174	182
184	253	250	265	3012	3002	3018
184	257	250	265	3005	3002	3018
184	260	250	265	3015	3002	3018
99	578	578	579	3460	3460	3461
34	580	580	581	3460	3460	3461
34	580	580	581	3461	3460	3461
34	579	578	583	3458	3458	3459
74	582	582	583	3460	3460	3461
62	580	580	581	3457	3457	3457
29	578	578	579	3457	3456	3457
19	579	578	579	3457	3456	3457
19	579	578	583	3458	3458	3459
4	588	576	597	3478	3464	3481
4	585	576	597	3479	3464	3481
4	586	576	597	3478	3464	3481
4	577	576	597	3479	3464	3481
4	580	576	597	3477	3464	3481
367	588	583	594	3474	3472	3477
367	589	583	594	3477	3472	3477
367	587	583	594	3475	3472	3477
367	584	583	594	3476	3472	3477
367	585	583	594	3474	3472	3477
367	583	583	594	3476	3472	3477
367	583	583	594	3475	3472	3477
41	666	657	675	578	575	587
41	666	657	675	581	575	587
41	662	657	675	583	575	587
41	661	657	675	580	575	587
41	668	657	675	578	575	587
41	664	657	675	576	575	587
29	666	657	675	578	575	587
29	670	657	675	577	575	587
29	664	657	675	578	575	587
29	689	675	701	564	564	570
29	690	675	701	564	564	570
29	691	675	701	564	564	570
29	692	675	701	564	564	570
29	693	675	701	564	564	570
29	694	675	701	564	564	570
29	695	675	701	564	564	570
29	689	675	701	564	564	570
29	690	675	701	564	564	570
29	691	675	701	564	564	570
29	692	675	701	564	564	570
29	693	675	701	564	564	570
29	694	675	701	564	564	570
29	695	675	701	564	564	570
29	689	675	701	564	564	570
29	690	675	701	564	564	570
29	691	675	701	564	564	570
29	692	675	701	564	564	570
29	693	675	701	564	564	570
29	694	675	701	564	564	570
29	695	675	701	564	564	570
29	689	675	701	564	564	570
29	690	675	701	564	564	570
29	691	675	701	564	564	570
29	692	675	701	564	564	570
29	693	675	701	564	564	570
29	694	675	701	564	564	570
29	695	675	701	564	564	570
311	522	520	524	3434	3433	3436
521	540	537	544	748	744	750
521	541	537	544	748	744	750
521	540	537	544	746	744	750
521	542	537	544	746	744	750
521	540	537	547	752	750	757
521	541	537	547	754	750	757
521	540	537	547	755	750	757
521	542	537	547	754	750	757
521	542	537	547	755	750	757
521	537	529	540	760	756	765
521	536	529	540	760	756	765
521	534	529	540	761	756	765
521	535	529	540	758	756	765
521	537	529	540	759	756	765
99	168	146	179	258	248	265
99	166	146	179	259	248	265
99	166	146	179	257	248	265
99	163	146	179	256	248	265
99	163	146	179	260	248	265
99	172	146	179	258	248	265
99	174	146	179	255	248	265
99	166	146	179	252	248	265
99	155	146	179	260	248	265
158	331	323	335	158	149	162
158	329	323	335	155	149	162
158	330	323	335	154	149	162
158	334	323	335	153	149	162
158	333	323	335	157	149	162
158	334	323	335	154	149	162
135	335	323	335	148	142	149
135	328	323	335	148	142	149
135	331	323	335	146	142	149
158	333	326	334	125	123	141
158	328	326	334	126	123	141
158	334	326	334	127	123	141
158	329	326	334	136	123	141
158	330	326	334	137	123	141
158	332	326	334	135	123	141
135	334	326	334	134	123	141
135	331	326	334	129	123	141
263	325	317	325	134	122	141
263	324	317	325	130	122	141
263	320	317	325	126	122	141
263	324	317	325	138	122	141
263	325	317	325	127	122	141
263	325	317	325	132	122	141
17	214	211	219	1640	1633	1642
81	214	211	215	2580	2580	2582
60	212	211	215	2581	2580	2582
722	20	3	30	526	505	542
722	20	3	30	526	505	542
722	20	3	30	526	505	542
722	20	3	30	526	505	542
722	20	3	30	526	505	542
722	20	3	30	526	505	542
722	20	3	30	526	505	542
724	20	3	30	526	505	542
724	20	3	30	526	505	542
724	20	3	30	526	505	542
724	20	3	30	526	505	542
724	20	3	30	526	505	542
725	20	3	30	526	505	542
725	20	3	30	526	505	542
725	20	3	30	526	505	542
725	20	3	30	526	505	542
725	20	3	30	526	505	542
727	20	3	30	526	505	542
727	20	3	30	526	505	542
727	20	3	30	526	505	542
727	20	3	30	526	505	542
727	20	3	30	526	505	542
726	6	3	14	541	539	546
726	6	3	14	541	539	546
726	6	3	14	541	539	546
726	6	3	14	541	539	546
723	20	17	23	568	563	573
723	20	17	23	568	563	573
723	20	17	23	568	563	573
723	20	17	23	568	563	573
723	20	17	23	568	563	573
728	26	24	30	571	570	573
728	26	24	30	571	570	573
728	26	24	30	571	570	573
40	9	5	17	3393	3386	3399
40	12	5	17	3392	3386	3399
40	13	5	17	3395	3386	3399
40	11	5	17	3396	3386	3399
40	9	5	17	3398	3386	3399
40	6	5	17	3395	3386	3399
40	8	5	17	3391	3386	3399
40	10	5	17	3390	3386	3399
40	16	5	17	3395	3386	3399
40	15	5	17	3398	3386	3399
40	12	5	17	3398	3386	3399
626	713	713	713	582	582	582
627	702	702	702	3420	3420	3420
628	723	723	723	3461	3461	3461
629	763	763	763	3441	3441	3441
631	715	712	718	3413	3410	3415
631	716	712	718	3411	3410	3415
631	716	712	718	3413	3410	3415
631	713	712	718	3412	3410	3415
631	715	712	718	3411	3410	3415
631	716	712	718	3415	3410	3415
631	714	712	717	3441	3438	3442
631	714	712	717	3439	3438	3442
631	716	712	717	3440	3438	3442
631	716	712	717	3441	3438	3442
631	715	712	717	3442	3438	3442
631	713	712	717	3438	3438	3442
631	715	712	717	3440	3438	3442
631	723	723	733	585	579	595
631	731	723	733	589	579	595
631	728	723	733	589	579	595
631	727	723	733	589	579	595
631	727	723	733	591	579	595
631	724	723	733	592	579	595
631	726	723	733	594	579	595
631	732	723	733	581	579	595
631	733	723	733	584	579	595
34	727	723	733	591	579	595
34	724	723	733	592	579	595
34	726	723	733	594	579	595
34	732	723	733	581	579	595
34	733	723	733	584	579	595
641	727	723	733	581	579	595
4	711	703	715	3431	3419	3438
4	715	703	715	3430	3419	3438
4	708	703	715	3430	3419	3438
4	708	703	715	3434	3419	3438
4	750	742	754	3461	3460	3464
4	745	742	754	3462	3460	3464
4	745	742	754	3463	3460	3464
4	744	742	754	3464	3460	3464
4	743	742	754	3461	3460	3464
43	708	703	715	3434	3419	3438
43	689	684	693	3437	3436	3440
43	691	684	693	3440	3436	3440
43	677	676	679	3448	3436	3450
43	676	676	679	3447	3436	3450
43	693	690	693	3448	3444	3451
43	763	756	765	3476	3472	3479
43	759	756	765	3477	3472	3479
43	760	756	765	3474	3472	3479
43	763	756	765	3474	3472	3479
43	765	756	765	3474	3472	3479
43	748	740	754	3489	3484	3490
43	725	722	730	3477	3474	3485
43	726	722	730	3476	3474	3485
43	726	722	730	3481	3474	3485
43	724	722	730	3482	3474	3485
43	724	722	730	3481	3474	3485
43	722	722	730	3481	3474	3485
43	727	722	730	3479	3474	3485
43	727	722	730	3481	3474	3485
43	741	737	744	586	577	587
43	738	737	744	584	577	587
43	741	737	744	579	577	587
43	737	737	744	578	577	587
43	758	754	765	623	607	625
43	762	754	765	620	607	625
43	758	754	765	618	607	625
43	762	754	765	614	607	625
41	707	700	710	3448	3442	3452
41	705	700	710	3446	3442	3452
41	704	700	710	3448	3442	3452
41	703	700	710	3447	3442	3452
41	706	700	710	3444	3442	3452
41	741	738	751	3427	3425	3434
41	742	738	751	3429	3425	3434
41	750	738	751	3428	3425	3434
542	750	738	751	3426	3425	3434
19	741	733	743	3440	3440	3442
19	736	733	743	3440	3440	3442
19	738	733	743	3440	3440	3442
23	734	733	743	3442	3440	3442
312	749	745	752	3441	3438	3443
312	751	745	752	3440	3438	3443
312	746	745	752	3440	3438	3443
630	735	730	746	3464	3460	3464
630	731	730	746	3462	3460	3464
630	730	730	746	3462	3460	3464
630	736	730	746	3464	3460	3464
630	745	740	754	3489	3484	3490
630	746	740	754	3489	3484	3490
630	747	740	754	3489	3484	3490
630	748	740	754	3489	3484	3490
630	745	740	754	3489	3484	3490
630	746	740	754	3489	3484	3490
630	747	740	754	3489	3484	3490
630	748	740	754	3489	3484	3490
630	745	740	754	3489	3484	3490
630	746	740	754	3489	3484	3490
630	747	740	754	3489	3484	3490
630	748	740	754	3489	3484	3490
634	745	743	746	3459	3457	3459
635	744	743	746	3459	3457	3459
638	745	743	746	3466	3465	3467
639	743	743	746	3466	3465	3467
636	749	748	751	3459	3457	3459
637	750	748	751	3459	3457	3459
640	750	748	751	3465	3465	3467
634	749	748	751	3465	3465	3467
45	762	752	765	3497	3494	3501
46	760	752	765	3497	3494	3501
46	758	752	765	3496	3494	3501
46	764	752	765	3498	3494	3501
46	763	752	765	3496	3494	3501
542	742	738	751	3429	3425	3434
542	744	738	751	3430	3425	3434
542	744	738	751	3429	3425	3434
542	748	738	751	3431	3425	3434
542	749	738	751	3430	3425	3434
542	748	738	751	3429	3425	3434
542	745	738	751	3428	3425	3434
542	744	738	751	3427	3425	3434
52	751	738	751	3429	3425	3434
52	739	738	751	3429	3425	3434
632	725	722	732	3416	3410	3421
633	725	722	732	3414	3410	3421
633	725	722	732	3414	3410	3421
658	740	737	741	624	622	625
642	759	755	764	661	657	667
648	761	755	764	661	657	667
657	763	755	764	661	657	667
655	751	724	760	633	600	640
655	753	724	760	628	600	640
655	750	724	760	626	600	640
655	748	724	760	627	600	640
655	748	724	760	629	600	640
655	746	724	760	632	600	640
655	749	724	760	634	600	640
655	751	724	760	635	600	640
655	754	724	760	635	600	640
655	741	724	760	629	600	640
655	739	724	760	630	600	640
655	740	724	760	632	600	640
655	736	724	760	631	600	640
655	734	724	760	628	600	640
655	734	724	760	627	600	640
655	736	724	760	624	600	640
655	737	724	760	621	600	640
655	735	724	760	618	600	640
655	740	724	760	616	600	640
655	742	724	760	615	600	640
655	743	724	760	613	600	640
655	746	724	760	611	600	640
655	742	724	760	606	600	640
655	739	724	760	605	600	640
655	738	724	760	607	600	640
655	737	724	760	609	600	640
655	736	724	760	611	600	640
655	735	724	760	613	600	640
655	736	724	760	615	600	640
655	738	724	760	622	600	640
67	224	213	237	253	238	267
67	226	213	237	252	238	267
67	229	213	237	253	238	267
67	228	213	237	255	238	267
67	226	213	237	257	238	267
67	223	213	237	256	238	267
67	225	213	237	250	238	267
67	224	213	237	248	238	267
67	226	213	237	248	238	267
67	228	213	237	248	238	267
67	232	213	237	253	238	267
67	229	213	237	256	238	267
232	270	266	280	329	319	332
234	270	255	281	302	289	322
234	272	255	281	301	289	322
234	274	255	281	301	289	322
232	272	255	281	308	289	322
232	271	255	281	318	289	322
232	269	255	281	319	289	322
235	269	267	273	290	289	292
457	493	493	495	615	615	618
456	516	514	517	612	611	614
460	483	483	484	614	614	620
461	484	483	486	1560	1558	1564
462	493	493	493	1560	1560	1560
463	494	494	494	1558	1558	1558
464	495	495	495	1564	1564	1564
463	490	490	490	1564	1564	1564
462	488	488	488	1561	1561	1561
464	490	490	490	1558	1558	1558
462	494	494	494	614	614	614
462	492	492	492	620	620	620
463	490	490	490	620	620	620
40	20	16	23	3380	3378	3386
40	21	16	23	3379	3378	3386
40	23	16	23	3380	3378	3386
40	22	16	23	3382	3378	3386
40	20	16	23	3383	3378	3386
40	20	16	23	3384	3378	3386
40	21	16	23	3384	3378	3386
40	18	16	23	3382	3378	3386
40	19	16	23	3379	3378	3386
40	21	16	23	3378	3378	3386
40	9	5	11	3374	3373	3381
40	8	5	11	3375	3373	3381
40	6	5	11	3376	3373	3381
40	7	5	11	3377	3373	3381
40	9	5	11	3377	3373	3381
40	8	5	11	3378	3373	3381
40	11	5	11	3378	3373	3381
40	9	5	11	3380	3373	3381
40	7	5	11	3380	3373	3381
40	6	5	11	3378	3373	3381
40	5	5	11	3376	3373	3381
67	654	647	659	3539	3533	3542
67	647	647	659	3538	3533	3542
41	650	647	659	3540	3533	3542
67	662	661	666	3533	3530	3536
41	664	661	666	3535	3530	3536
352	620	618	622	3497	3494	3500
351	619	618	622	3499	3494	3500
351	616	616	622	3496	3494	3500
195	558	550	559	3293	3281	3294
195	556	550	559	3287	3281	3294
195	555	550	559	3285	3281	3294
195	553	533	557	3284	3281	3292
195	552	533	557	3283	3281	3292
195	551	533	557	3284	3281	3292
195	547	533	557	3285	3281	3292
195	546	533	557	3284	3281	3292
22	534	531	535	3291	3286	3296
22	533	531	535	3294	3289	3299
22	533	531	535	3297	3292	3302
70	546	544	547	3333	3330	3337
70	545	544	547	3332	3330	3337
70	546	544	547	3333	3330	3337
70	544	544	547	3331	3330	3337
45	125	105	140	262	255	273
45	126	105	140	266	255	273
45	118	105	140	265	255	273
45	123	105	123	280	273	286
45	120	105	123	278	273	286
45	123	105	123	275	273	286
45	120	105	121	280	273	286
276	460	460	463	2407	2406	2407
521	469	466	473	734	733	737
521	471	466	473	735	733	737
290	389	382	391	3370	3360	3374
290	389	382	391	3364	3360	3374
290	384	382	391	3370	3360	3374
290	386	382	391	3364	3360	3374
202	370	365	374	3352	3347	3357
202	370	365	374	3355	3350	3363
22	213	211	213	2585	2583	2585
40	630	629	634	3295	3293	3307
40	631	629	634	3297	3293	3307
40	631	629	634	3300	3293	3307
40	633	629	634	3296	3293	3307
41	630	629	634	3301	3293	3307
41	633	629	634	3294	3293	3307
41	634	629	634	3301	3293	3307
104	632	629	634	3297	3293	3307
41	637	634	648	3306	3301	3308
41	639	634	648	3305	3301	3308
40	637	634	648	3304	3301	3308
40	639	634	648	3306	3301	3308
40	640	634	648	3307	3301	3308
40	641	634	648	3302	3301	3308
43	664	663	670	3290	3290	3296
46	668	663	670	3296	3290	3297
19	665	663	670	3297	3290	3297
344	662	656	663	3286	3282	3288
344	660	656	663	3286	3282	3288
344	657	656	663	3285	3282	3288
344	661	656	663	3283	3282	3288
344	665	663	670	3293	3290	3294
344	667	663	670	3292	3290	3294
195	670	667	670	3288	3281	3289
195	668	667	670	3285	3281	3289
195	670	667	670	3285	3281	3289
195	654	648	662	3293	3289	3302
195	651	648	662	3292	3289	3302
195	653	648	662	3291	3289	3302
195	650	648	662	3292	3289	3302
195	650	648	662	3293	3289	3302
104	632	631	641	508	505	516
104	635	631	641	511	505	516
104	631	631	641	491	485	496
29	663	660	670	3277	3269	3279
29	668	660	670	3277	3269	3279
29	663	660	670	3278	3269	3279
29	661	660	670	3272	3269	3279
29	667	660	670	3276	3269	3279
343	658	652	664	3297	3291	3303
343	659	653	665	3296	3290	3302
343	659	653	665	3298	3292	3304
343	661	655	667	3297	3291	3303
343	658	652	664	3298	3292	3304
498	655	648	662	3297	3289	3302
498	653	648	662	3296	3289	3302
498	652	648	662	3294	3289	3302
135	421	416	426	448	443	453
135	425	420	430	446	441	451
99	256	251	260	152	152	160
99	260	251	260	154	152	160
99	259	251	260	159	152	160
99	253	251	260	156	152	160
409	644	642	646	666	664	668
409	644	642	646	664	662	666
409	645	643	647	667	665	669
409	667	665	669	662	660	664
409	668	666	670	663	661	665
409	669	667	671	667	665	669
409	669	667	671	669	667	671
409	669	667	671	663	661	665
408	656	656	656	662	662	662
34	11	8	14	3395	3392	3398
40	10	8	13	3394	3392	3397
40	11	9	13	3397	3395	3398
40	11	8	14	3395	3392	3398
190	270	260	300	234	234	265
190	276	260	300	240	234	265
190	296	260	300	250	234	265
190	282	260	300	240	234	265
146	330	329	332	660	658	664
156	331	328	332	668	666	670
166	270	268	270	654	646	654
171	268	268	270	646	646	654
163	325	323	333	713	713	713
163	331	323	333	713	713	713
317	467	467	467	649	646	655
316	537	534	542	615	615	616
145	330	329	332	662	658	664
214	418	409	418	3377	3373	3382
214	411	409	418	3377	3373	3382
214	410	409	418	3379	3373	3382
214	415	409	418	3381	3373	3382
214	413	409	418	3382	3373	3382
184	413	410	414	3390	3387	3395
184	412	410	414	3392	3387	3395
184	413	410	414	3394	3387	3395
297	435	434	437	565	561	567
395	406	402	413	562	559	565
3	412	405	412	545	545	550
3	410	405	412	546	545	550
3	407	405	412	545	545	550
3	405	405	412	546	545	550
3	406	405	412	548	545	550
3	405	405	412	550	545	550
3	407	405	412	550	545	550
0	432	399	435	535	531	540
0	425	399	435	540	531	540
0	415	399	435	540	531	540
8	432	431	436	540	532	548
213	427	426	428	548	546	549
213	426	426	428	546	546	549
239	409	408	411	483	480	485
249	402	402	403	466	463	469
249	403	402	403	466	463	469
239	395	393	398	479	475	483
239	399	396	403	481	476	483
239	395	391	399	489	485	493
621	395	394	396	1766	1766	1770
624	415	411	415	833	833	835
320	616	608	621	531	528	541
320	619	608	621	533	528	541
320	617	608	621	535	528	541
320	614	608	621	539	528	541
432	621	621	621	529	529	529
432	622	622	622	529	529	529
432	623	623	623	529	529	529
430	578	578	578	566	566	566
430	578	578	578	565	565	565
430	578	578	578	564	564	564
431	581	581	581	564	564	564
541	417	414	418	163	161	165
735	517	515	518	545	543	548
736	511	510	516	551	551	552
736	514	510	516	551	551	552
788	512	510	515	1481	1479	1483
95	515	510	517	2422	2421	2423
95	512	510	517	2423	2421	2423
95	511	510	517	2422	2421	2423
779	510	508	513	2426	2424	2427
513	597	596	597	756	755	758
269	372	370	377	438	435	440
253	372	368	375	443	441	445
24	111	105	115	3369	3362	3372
16	139	133	143	475	469	479
15	108	102	112	674	668	678
117	218	212	222	2580	2574	2584
128	252	246	256	632	626	636
119	70	64	74	697	691	701
120	76	70	80	683	677	687
4	354	349	359	508	503	513
4	361	356	366	506	501	511
205	379	374	384	487	482	492
274	472	467	477	450	445	455
275	474	469	479	453	448	458
393	462	457	467	1387	1382	1392
63	492	487	497	487	482	492
63	493	488	498	493	488	498
301	524	519	529	1406	1401	1411
333	545	542	548	576	576	580
327	549	544	554	585	580	590
326	565	560	570	594	589	599
435	578	573	583	591	586	596
486	580	575	585	586	581	591
335	560	555	565	579	574	584
455	521	516	526	619	614	624
512	613	608	618	602	597	607
437	620	615	625	581	576	586
419	616	611	621	617	612	622
422	611	606	616	620	615	625
418	619	614	624	620	615	625
423	616	611	621	1565	1560	1570
418	613	608	618	1562	1557	1567
419	613	608	618	1562	1557	1567
366	580	575	585	637	632	642
350	580	575	585	665	660	670
372	575	570	580	679	674	684
733	538	533	543	703	698	708
733	538	533	543	703	698	708
381	589	584	594	715	710	720
382	591	586	596	718	713	723
385	591	586	596	718	713	723
385	589	584	594	708	703	713
385	588	583	593	701	696	706
385	588	583	593	701	696	706
381	590	585	595	697	692	702
381	598	593	603	685	680	690
381	599	594	604	686	681	691
385	603	598	608	688	683	693
385	602	597	607	687	682	692
385	600	595	605	679	674	684
385	610	605	615	684	679	689
385	607	602	612	693	688	698
381	605	600	610	695	690	700
385	618	613	623	688	683	693
385	618	613	623	688	683	693
385	618	613	623	688	683	693
385	618	613	623	688	683	693
385	618	613	623	688	683	693
385	618	613	623	688	683	693
381	619	614	624	679	674	684
385	621	616	626	674	669	679
575	639	634	644	741	736	746
575	639	634	644	741	736	746
575	639	634	644	741	736	746
72	614	609	619	763	758	768
72	614	609	619	763	758	768
529	629	624	634	766	761	771
680	662	657	667	738	733	743
434	595	590	600	539	534	544
345	583	578	588	624	619	629
754	485	480	490	403	398	408
748	482	477	487	385	380	390
742	481	476	486	389	384	394
746	501	496	506	388	383	393
756	495	490	500	385	380	390
751	496	491	501	389	384	394
752	485	480	490	385	380	390
755	495	490	500	1329	1324	1334
745	495	490	500	1334	1329	1339
753	488	483	493	1329	1324	1334
741	487	482	492	1336	1331	1341
743	494	489	499	390	385	395
747	496	491	501	399	394	404
747	482	477	487	402	397	407
750	496	491	501	416	411	421
305	494	489	499	547	542	552
305	494	489	499	547	542	552
6	550	545	555	561	556	566
6	550	545	555	558	553	563
6	549	544	554	556	551	561
6	548	543	553	551	546	556
6	547	542	552	548	543	553
6	540	535	545	566	561	571
2	552	547	557	552	547	557
2	551	546	556	549	544	554
2	551	546	556	551	546	556
2	551	546	556	548	543	553
2	551	546	556	551	546	556
94	612	607	617	478	473	483
94	612	607	617	478	473	483
771	605	600	610	467	462	472
43	612	607	617	460	455	465
43	613	608	618	462	457	467
43	615	610	620	461	456	466
472	657	652	662	491	486	496
481	654	649	659	500	495	505
207	178	173	183	667	662	672
208	178	173	183	667	662	672
209	177	172	182	669	664	674
210	180	175	185	669	664	674
300	511	506	516	1452	1447	1457
302	380	376	386	3353	3349	3359
355	427	422	432	456	451	461
715	126	121	131	504	499	509
170	268	263	273	648	643	653
3	560	555	565	495	490	500
3	562	557	567	493	488	498
3	560	555	565	491	486	496
3	558	553	563	492	487	497
349	563	558	568	493	488	498
346	568	563	573	489	484	494
353	570	565	575	500	495	505
354	570	565	575	503	498	508
347	566	561	571	491	486	496
355	384	379	389	464	459	469
345	560	555	565	485	480	490
291	407	397	412	3336	3331	3341
291	412	407	422	3337	3332	3342
294	394	389	399	3314	3309	3319
95	102	97	107	515	510	520
95	280	275	285	566	561	571
95	280	275	285	568	563	573
95	333	328	334	551	549	557
81	280	275	285	609	604	614
219	106	101	111	3545	3540	3550
220	113	108	118	3546	3541	3551
218	107	102	112	3541	3536	3546
221	113	108	118	3537	3532	3542
221	117	112	122	3541	3536	3546
747	488	483	493	411	406	416
747	495	490	500	418	413	423
355	386	381	391	471	466	476
262	547	542	552	455	450	460
262	548	543	553	470	465	475
262	566	561	571	438	433	443
262	569	564	574	472	467	477
70	593	590	598	619	616	623
243	584	577	591	610	603	617
431	581	581	581	565	565	565
431	581	581	581	566	566	566
433	603	603	603	585	585	585
433	603	603	603	586	586	586
70	595	590	598	620	616	623
357	617	612	622	638	633	643
458	512	507	517	638	633	643
459	510	505	515	636	631	641
374	621	616	626	701	696	706
379	619	614	624	701	696	706
380	619	614	624	707	702	712
385	616	611	621	713	708	718
387	609	604	614	715	710	720
377	605	600	610	716	711	721
373	604	599	609	717	712	722
377	613	613	613	709	709	709
378	613	613	613	707	707	707
396	625	625	625	675	675	675
4	627	622	632	679	674	684
61	629	624	634	679	674	684
4	631	626	636	677	672	682
47	625	620	630	689	684	694
402	634	629	639	681	676	686
4	637	632	642	681	676	686
4	640	635	645	682	677	687
4	645	640	650	683	678	688
47	646	641	651	679	674	684
47	627	622	632	695	690	700
29	625	620	630	700	695	705
4	627	622	632	710	705	715
4	627	622	632	705	700	710
4	631	626	636	714	709	719
402	629	624	634	716	711	721
29	625	620	630	718	713	723
4	631	626	636	718	713	723
4	635	630	640	718	713	723
29	651	646	656	717	712	722
4	633	628	638	714	709	719
4	641	636	646	714	709	719
402	648	643	653	712	707	717
4	653	648	658	708	703	713
4	660	655	665	708	703	713
4	662	657	667	708	703	713
61	662	657	667	714	709	719
402	659	654	664	703	698	708
4	664	659	669	703	698	708
4	652	647	657	708	703	713
4	668	663	673	718	713	723
47	640	635	645	708	703	713
4	631	626	636	693	688	698
4	639	634	644	710	705	715
4	663	658	668	687	682	692
4	658	653	663	683	678	688
4	660	655	665	681	676	686
4	662	657	667	681	676	686
402	662	657	667	679	674	684
4	663	658	668	675	670	680
4	666	661	671	675	670	680
402	653	648	658	674	669	679
29	651	646	656	673	668	678
4	642	637	647	675	670	680
4	633	628	638	675	670	680
8	650	645	655	683	678	688
399	640	635	645	699	694	704
399	645	640	650	697	692	702
399	635	630	640	696	691	701
397	640	640	640	698	698	698
399	645	640	650	688	683	693
400	656	656	656	695	695	695
406	668	668	668	638	638	638
404	653	653	653	627	627	627
405	649	644	654	629	624	634
652	715	710	720	681	676	686
652	714	709	719	679	674	684
652	715	710	720	682	677	687
654	714	709	719	683	678	688
367	712	707	717	3513	3508	3518
660	717	712	722	3512	3507	3517
367	715	710	720	3506	3501	3511
660	713	708	718	3506	3501	3511
660	705	700	710	3506	3501	3511
660	703	698	708	3511	3506	3516
660	709	704	714	3516	3511	3521
660	706	701	711	3523	3518	3528
4	698	693	703	3512	3507	3517
660	699	694	704	3507	3502	3512
660	686	681	691	3507	3502	3512
660	682	677	687	3515	3510	3520
4	675	670	680	3506	3501	3511
660	677	672	682	3511	3506	3516
4	678	673	683	3522	3517	3527
660	681	676	686	3527	3522	3532
660	681	676	686	3536	3531	3541
4	695	690	700	3537	3532	3542
660	688	683	693	3538	3533	3543
367	694	689	699	3549	3544	3554
660	675	670	680	3540	3535	3545
660	704	699	709	3544	3539	3549
660	705	700	710	3550	3545	3555
367	712	707	717	3548	3543	3553
660	714	709	719	3529	3524	3534
4	685	680	690	3513	3508	3518
651	690	690	690	3512	3512	3512
660	716	711	721	3521	3516	3526
660	716	711	721	3533	3528	3538
662	713	708	718	698	693	703
450	615	610	620	581	576	586
437	625	620	630	3417	3412	3422
29	623	618	628	3416	3411	3421
29	630	625	635	3411	3406	3416
29	634	629	639	3411	3406	3416
19	631	626	636	3421	3416	3426
443	634	629	639	589	584	594
445	633	628	638	582	577	587
446	644	639	649	566	561	571
448	646	641	651	566	561	571
447	646	641	651	567	562	572
449	645	640	650	1511	1506	1516
445	636	631	641	603	598	608
445	631	626	636	604	599	609
445	628	623	633	603	598	608
29	630	625	635	606	601	611
441	633	628	638	604	599	609
440	634	629	639	599	594	604
29	651	646	656	584	579	589
34	635	630	640	607	602	612
34	635	630	640	609	604	614
34	639	634	644	607	602	612
465	637	632	642	3449	3444	3454
454	646	641	651	586	581	591
593	691	686	696	1412	1407	1417
592	691	686	696	1412	1407	1417
586	691	686	696	1412	1407	1417
580	692	687	697	1404	1399	1409
580	690	685	695	1403	1398	1408
593	694	689	699	1402	1397	1407
534	691	686	696	1404	1399	1409
593	685	680	690	1402	1397	1407
592	683	678	688	1401	1396	1406
586	683	678	688	1401	1396	1406
535	685	680	690	1402	1397	1407
586	693	688	698	1396	1391	1401
592	696	691	701	1394	1389	1399
592	690	685	695	1388	1383	1393
592	690	685	695	1390	1385	1395
592	690	685	695	1380	1375	1385
587	691	686	696	1379	1374	1384
586	685	680	690	1389	1384	1394
586	682	677	687	1390	1385	1395
586	683	678	688	1388	1383	1393
532	717	712	722	1409	1404	1414
625	747	742	752	448	443	453
601	746	746	748	450	446	453
596	728	728	728	450	450	450
611	736	735	737	445	444	446
605	741	740	742	446	445	447
606	742	741	743	449	448	450
609	741	740	741	443	443	444
610	741	740	741	456	455	456
611	736	735	737	455	454	456
588	489	484	494	555	550	560
763	490	485	495	452	447	457
394	411	406	416	560	555	565
68	169	164	174	310	305	315
41	168	163	173	309	304	314
41	168	163	173	310	305	315
47	323	318	328	301	296	306
13	62	57	67	683	678	688
13	63	58	68	740	735	745
13	76	71	81	751	746	756
13	58	53	63	770	765	775
13	63	58	68	786	781	791
13	65	60	70	801	796	806
13	114	109	119	804	799	809
13	145	140	150	803	798	808
13	152	147	157	774	769	779
13	171	166	176	788	783	793
13	179	174	184	804	799	809
13	131	126	136	785	780	790
13	128	123	133	762	757	767
13	116	111	121	745	740	750
13	93	88	98	734	729	739
653	72	67	77	751	746	756
653	70	65	75	765	760	770
653	96	91	101	761	756	766
653	111	106	116	754	749	759
653	118	113	123	768	763	773
653	118	113	123	785	780	790
653	97	92	102	795	790	800
653	128	123	133	800	795	805
653	157	152	162	792	787	797
653	165	160	170	809	804	814
653	59	54	64	751	746	756
780	90	85	95	519	514	524
70	67	62	72	752	747	757
70	54	49	59	769	764	774
70	86	81	91	778	773	783
70	95	90	100	790	785	795
70	109	104	114	788	783	793
70	131	126	136	799	794	804
70	131	126	136	769	764	774
70	125	120	130	753	748	758
70	109	104	114	733	728	738
70	108	103	113	738	733	743
70	87	82	92	731	726	736
95	149	144	154	500	495	505
270	219	214	224	3244	3239	3249
270	217	212	222	3244	3239	3249
593	680	675	685	1413	1408	1418
76	487	482	492	549	544	554
76	490	485	495	547	542	552
76	491	486	496	544	539	549
76	491	486	496	550	545	555
76	491	486	496	554	549	559
3	490	485	495	548	543	553
3	492	487	497	545	540	550
243	300	295	305	134	129	139
243	303	298	308	135	130	140
243	303	298	308	133	128	138
243	300	295	305	134	129	139
243	296	291	301	130	125	135
243	296	291	301	132	127	137
243	296	291	301	129	124	134
243	296	291	301	130	125	135
45	294	289	299	2940	2935	2945
45	294	289	299	2944	2939	2949
45	300	295	305	2940	2935	2945
45	291	286	296	2940	2935	2945
34	289	284	294	2936	2931	2941
34	291	286	296	2934	2929	2939
34	302	297	307	2944	2939	2949
45	308	303	313	114	109	119
273	462	457	467	432	427	437
287	462	461	462	2336	2336	2336
475	666	661	671	3526	3521	3531
470	656	651	661	448	443	453
471	664	659	669	464	459	469
263	402	397	407	3276	3271	3281
263	409	404	414	3276	3271	3281
263	411	406	416	3282	3277	3287
263	405	400	410	3284	3279	3289
263	398	393	403	3284	3279	3289
263	405	400	410	3279	3274	3284
263	399	394	404	3280	3275	3285
254	395	390	400	3272	3267	3277
360	615	610	620	587	582	592
340	619	614	624	587	582	592
484	624	619	629	616	611	621
390	540	535	545	703	698	708
546	534	529	539	755	750	760
675	630	625	635	792	787	797
675	631	626	636	791	786	796
312	622	617	627	784	779	789
312	617	612	622	783	778	788
312	612	607	617	785	780	790
312	612	607	617	791	786	796
8	608	603	613	774	769	779
94	603	598	608	780	775	785
94	620	615	625	776	771	781
312	621	616	626	804	799	809
312	618	613	623	810	805	815
312	624	619	629	817	812	822
312	620	615	625	819	814	824
531	621	616	626	827	822	832
312	628	623	633	823	818	828
312	630	625	635	833	828	838
312	621	616	626	837	832	842
312	616	611	621	844	839	849
312	625	620	630	853	848	858
312	627	622	632	853	848	858
243	619	614	624	851	846	856
243	611	606	616	853	848	858
61	619	614	624	854	849	859
19	610	605	615	847	842	852
243	595	590	600	856	851	861
243	585	580	590	853	848	858
61	593	588	598	855	850	860
61	588	583	593	853	848	858
312	587	582	592	832	827	837
312	591	586	596	835	830	840
531	590	585	595	832	827	837
672	636	631	641	2626	2621	2631
672	639	634	644	2627	2622	2632
672	637	632	642	2625	2620	2630
312	647	642	652	839	834	844
312	657	652	662	829	824	834
243	662	657	667	853	848	858
312	663	658	668	846	841	851
312	668	663	673	843	838	848
61	649	644	654	852	847	857
95	439	434	444	491	486	496
95	589	584	594	756	751	761
95	502	497	507	447	442	452
70	252	247	257	399	394	404
268	91	88	93	691	689	696
62	291	286	296	666	661	671
523	419	414	424	3575	3570	3580
523	423	418	428	3580	3575	3585
523	414	409	419	3577	3572	3582
539	382	377	387	852	847	857
622	399	394	404	846	841	851
619	417	412	422	853	848	858
616	395	390	400	1780	1775	1785
721	54	49	59	753	748	758
721	70	65	75	764	759	769
721	87	82	92	765	760	770
721	112	107	117	763	758	768
721	135	130	140	774	769	779
721	157	152	162	785	780	790
721	158	153	163	801	796	806
721	180	175	185	808	803	813
721	125	120	130	801	796	806
721	102	97	107	800	795	805
303	487	485	495	543	538	548
23	607	602	612	3460	3455	3465
23	607	602	612	3466	3461	3471
34	605	600	610	3463	3458	3468
34	607	602	612	3468	3463	3473
19	607	602	612	3462	3457	3467
34	601	596	606	3465	3460	3470
40	598	593	603	3460	3455	3465
40	578	573	583	3464	3459	3469
40	583	578	588	3458	3453	3463
312	607	602	612	3481	3476	3486
312	613	608	618	3484	3479	3489
312	603	598	608	3477	3472	3482
153	581	576	586	3463	3458	3468
523	422	417	427	3587	3582	3592
403	585	580	590	460	455	465
29	585	580	590	460	455	465
29	585	580	590	461	456	466
29	582	577	587	449	444	454
401	413	411	413	11	9	12
416	404	399	409	33	28	38
416	409	404	414	34	29	39
416	398	393	403	27	22	32
416	397	392	402	35	30	40
416	401	396	406	36	31	41
416	401	396	406	30	25	35
416	416	411	421	27	22	32
414	392	387	397	29	24	34
413	414	409	419	35	30	40
413	415	410	420	36	31	41
413	419	414	424	33	28	38
413	423	418	428	34	29	39
499	224	219	229	742	737	747
474	231	226	236	735	730	740
474	229	224	234	734	729	739
478	215	210	220	730	725	735
478	214	209	219	726	721	731
478	216	211	221	726	721	731
480	209	204	214	731	726	736
480	208	203	213	728	723	733
480	207	202	212	732	727	737
479	201	196	206	725	720	730
479	199	194	204	723	718	728
479	198	193	203	728	723	733
482	198	193	203	738	733	743
482	203	198	208	742	737	747
482	205	200	210	738	733	743
485	198	193	203	748	743	753
485	198	193	203	752	747	757
485	200	195	205	751	746	756
489	206	201	211	750	745	755
489	209	204	214	753	748	758
489	209	204	214	748	743	753
493	206	201	211	758	753	763
493	210	205	215	760	755	765
493	211	206	216	756	751	761
494	218	213	223	758	753	763
494	220	215	225	758	753	763
3	219	214	224	755	750	760
3	218	213	223	755	750	760
3	219	214	224	756	751	761
774	224	219	229	758	753	763
774	224	219	229	761	756	766
774	222	217	227	761	756	766
496	229	224	234	761	756	766
496	231	226	236	760	755	765
497	226	221	231	753	748	758
413	515	510	520	33	28	38
413	517	512	522	37	32	42
413	517	512	522	39	34	44
413	520	515	525	33	28	38
415	516	511	521	35	30	40
417	505	500	510	26	21	31
417	506	501	511	21	16	26
417	515	510	520	24	19	29
417	490	485	495	31	26	36
417	494	489	499	27	22	32
263	542	537	547	3349	3344	3354
263	547	542	552	3350	3345	3355
263	543	538	548	3351	3346	3356
319	559	554	564	553	548	558
319	559	554	564	550	545	555
491	621	616	626	589	584	594
483	605	600	610	573	568	578
424	619	614	624	3477	3472	3482
321	532	527	537	595	590	600
318	607	602	612	584	579	589
318	606	601	611	585	580	590
318	606	601	611	580	575	585
590	739	734	744	3334	3329	3339
592	743	738	748	1406	1401	1411
593	742	737	747	1404	1399	1409
412	419	414	424	34	29	39
509	94	89	99	528	523	533
520	94	89	99	526	521	531
508	82	77	87	534	529	539
315	232	226	236	3248	3240	3250
190	341	336	346	3365	3360	3370
190	342	337	347	3367	3362	3372
190	342	337	347	3369	3364	3374
190	339	334	344	3365	3360	3370
8	214	209	219	370	365	375
313	311	306	316	3348	3343	3353
307	67	62	72	601	596	606
563	708	703	713	534	529	539
500	647	642	652	603	598	608
495	632	627	637	1514	1509	1519
361	542	538	548	3299	3295	3305
365	552	547	557	3292	3287	3297
487	629	624	634	613	608	618
547	714	709	719	1424	1419	1429
565	739	734	744	1382	1377	1387
593	744	739	749	1382	1377	1387
592	746	741	751	1381	1376	1386
591	749	744	754	1382	1377	1387
592	751	746	756	1381	1376	1386
538	64	59	69	738	733	743
669	95	90	100	806	801	811
671	83	78	88	807	802	812
671	85	80	90	809	804	814
669	99	94	104	807	802	812
669	98	93	103	805	800	810
711	88	83	93	807	802	812
702	85	80	90	1747	1739	1749
700	172	167	177	807	802	812
701	172	167	177	807	802	812
701	175	170	180	801	796	806
701	176	171	181	797	792	802
701	171	166	176	795	790	800
701	169	164	174	794	789	799
701	175	170	180	799	794	804
701	169	164	174	803	798	808
718	83	78	88	3625	3620	3630
718	87	82	92	3620	3615	3625
716	89	84	94	3623	3618	3628
716	89	84	94	3638	3633	3643
718	92	87	97	3627	3622	3632
670	83	78	88	3641	3636	3646
716	72	67	77	3641	3636	3646
718	73	68	78	3638	3633	3643
718	67	62	72	3642	3637	3647
718	53	48	58	3629	3624	3634
716	52	47	57	3625	3620	3630
718	56	51	61	3615	3610	3620
718	59	54	64	3612	3607	3617
716	61	56	66	3612	3607	3617
718	62	57	67	3623	3618	3628
718	68	63	73	3625	3620	3630
716	64	59	69	3618	3613	3623
718	67	62	72	3615	3610	3620
716	71	66	76	3606	3601	3611
668	70	65	75	3605	3600	3610
668	68	63	73	3605	3600	3610
554	68	63	73	3605	3600	3610
671	71	66	76	3605	3600	3610
671	65	60	70	3606	3601	3611
671	76	71	81	3618	3613	3623
671	56	51	61	3606	3601	3611
671	51	46	56	3608	3603	3613
542	352	347	357	3652	3647	3657
40	354	349	359	3654	3649	3659
542	354	349	359	3667	3662	3672
40	345	340	350	3712	3707	3717
542	345	340	350	3719	3714	3724
542	352	347	357	3720	3715	3725
40	354	349	359	3715	3710	3720
542	360	355	365	3717	3712	3722
40	466	461	471	3666	3661	3671
542	469	464	474	3661	3656	3666
542	347	342	352	3612	3607	3617
542	349	344	354	3614	3609	3619
542	347	342	352	3622	3617	3627
40	355	350	360	3624	3619	3629
542	356	351	361	3630	3625	3635
542	363	358	368	3627	3622	3632
542	371	366	376	3637	3632	3642
542	365	360	370	3609	3604	3614
542	347	342	352	3634	3629	3639
542	354	349	359	3674	3669	3679
542	353	348	358	3680	3675	3685
542	355	350	360	3683	3678	3688
40	351	346	356	3671	3666	3676
542	358	353	363	3672	3667	3677
576	696	695	696	492	492	494
578	694	694	695	1451	1451	1451
577	696	695	696	502	502	504
577	689	688	690	504	502	504
291	272	267	277	2999	2989	3004
291	273	268	278	2998	2988	3003
578	690	690	691	1450	1449	1451
579	687	686	688	502	502	504
579	685	685	687	500	499	501
579	681	681	682	500	500	501
576	682	681	682	492	492	493
8	311	306	316	461	456	466
8	311	306	316	467	462	472
8	303	298	308	596	591	601
177	347	333	353	635	617	637
47	343	331	351	635	623	643
47	341	336	356	628	619	639
67	359	349	369	614	604	624
95	214	212	220	449	448	453
95	217	212	220	452	448	453
95	213	212	220	452	448	453
4	206	196	216	497	487	507
4	204	194	214	500	490	510
53	322	312	332	282	272	292
53	316	306	326	278	268	288
53	316	306	326	285	275	295
53	310	300	320	284	274	294
53	317	307	327	276	266	286
53	319	309	329	277	267	287
3	450	440	460	763	753	773
3	451	441	461	768	758	778
3	438	428	448	764	754	774
3	435	425	445	771	761	781
421	453	443	463	781	771	791
4	442	432	452	800	790	810
70	460	450	470	802	792	812
67	450	440	460	802	792	812
523	473	463	483	819	809	829
523	459	449	469	815	805	825
521	463	453	473	816	806	826
521	450	440	460	820	810	830
521	461	451	471	856	846	866
521	471	461	481	860	850	870
523	476	466	486	865	855	875
523	462	452	472	861	851	871
523	462	452	472	861	851	871
775	457	447	467	875	865	885
775	475	465	485	875	865	885
765	459	449	469	863	853	873
523	360	350	370	845	835	855
523	352	342	362	843	833	853
421	358	348	368	843	833	853
421	347	337	357	827	817	837
421	338	328	348	823	813	833
421	348	338	358	836	826	846
421	359	349	369	833	823	843
521	353	343	363	828	818	838
521	337	327	347	832	822	842
521	342	332	352	814	804	824
521	336	326	346	817	807	827
523	336	326	346	817	807	827
523	337	327	347	800	790	810
523	350	340	360	801	791	811
521	344	334	354	793	783	803
521	343	333	353	793	783	803
421	338	328	348	774	764	784
523	366	356	376	761	751	771
523	369	359	379	766	756	776
421	366	356	376	815	805	825
421	374	364	384	817	807	827
521	370	360	380	796	786	806
521	381	371	391	820	810	830
521	358	348	368	816	806	826
53	370	360	380	3319	3309	3329
53	361	351	371	3320	3310	3330
53	350	340	360	3321	3311	3331
293	406	396	416	3350	3340	3360
293	411	401	421	3349	3339	3359
293	410	400	420	3345	3335	3355
293	406	396	416	3346	3336	3356
11	498	488	508	411	401	421
750	486	476	496	417	407	427
750	493	483	503	414	404	424
750	486	476	496	423	413	433
750	492	482	502	411	401	421
750	498	488	508	418	408	428
3	499	489	509	406	396	416
3	496	491	501	403	398	408
3	493	488	498	405	400	410
195	538	533	543	3324	3319	3329
195	530	525	535	3327	3322	3332
195	531	526	536	3331	3326	3336
195	536	531	541	3330	3325	3335
195	533	528	538	3323	3318	3328
46	535	530	540	537	532	542
46	546	541	551	538	533	543
188	529	524	534	566	561	571
188	515	510	520	567	562	572
188	517	512	522	580	575	585
188	530	525	535	579	574	584
188	521	516	526	574	569	579
188	527	522	532	573	568	578
188	516	511	521	563	558	568
19	561	556	566	3419	3414	3424
19	558	553	563	3424	3419	3429
19	557	552	562	3435	3430	3440
19	563	558	568	3435	3430	3440
41	546	541	551	3435	3430	3440
41	549	544	554	3433	3428	3438
4	581	576	586	3484	3479	3489
4	584	579	589	3481	3476	3486
4	577	572	582	3483	3478	3488
4	577	572	582	3478	3473	3483
4	609	604	614	505	500	510
4	617	612	622	503	498	508
4	613	608	618	513	508	518
4	605	600	610	499	494	504
4	623	618	628	502	497	507
4	618	613	623	497	492	502
4	618	613	623	514	509	519
4	610	605	615	518	513	523
4	618	613	623	506	501	511
4	609	604	614	624	619	629
318	609	604	614	1523	1518	1528
318	609	604	614	1525	1520	1530
322	607	602	612	606	601	611
206	608	608	608	2490	2490	2490
206	608	608	608	2492	2492	2492
206	613	613	613	2490	2490	2490
206	613	613	613	2492	2492	2492
0	375	370	380	457	452	462
312	617	612	622	3391	3386	3396
8	765	760	770	516	511	521
8	748	743	753	522	517	527
4	767	762	772	511	506	516
4	750	745	755	490	485	495
4	760	755	765	490	485	495
4	748	743	753	482	477	487
4	762	757	767	487	482	492
2	755	750	760	446	441	451
243	742	737	747	566	561	571
243	757	752	762	559	554	564
243	759	754	764	567	562	572
243	728	723	733	565	560	570
243	725	720	730	566	561	571
4	742	737	747	566	561	571
4	733	728	738	571	566	576
4	727	722	732	569	564	574
4	737	732	742	564	559	569
4	742	737	747	569	564	574
4	748	743	753	569	564	574
243	713	708	718	569	564	574
243	716	711	721	567	562	572
243	714	709	719	563	558	568
243	712	707	717	570	565	575
4	721	716	726	567	562	572
4	722	717	727	570	565	575
4	731	726	736	564	559	569
312	616	611	621	3484	3479	3489
407	639	634	644	645	640	650
407	637	632	642	642	637	647
407	644	639	649	648	643	653
407	640	635	645	647	642	652
389	612	607	617	453	448	458
407	648	638	658	633	623	643
407	648	638	658	636	626	646
407	653	643	663	641	631	651
407	658	648	668	643	633	653
407	666	656	676	640	630	650
407	660	650	670	638	628	648
407	658	648	668	629	619	639
407	644	634	654	633	623	643
407	651	641	661	640	630	650
407	655	645	665	644	634	654
407	647	637	657	640	630	650
407	645	635	655	628	618	638
407	647	637	657	625	615	635
407	641	631	651	630	620	640
399	635	625	645	1639	1629	1649
399	641	631	651	1636	1626	1646
555	601	596	606	3577	3572	3582
555	599	594	604	3575	3570	3580
555	596	591	601	3574	3569	3579
555	600	595	605	3570	3565	3575
555	601	596	606	3567	3562	3572
400	648	638	658	1635	1625	1645
398	635	625	645	1633	1623	1643
195	609	599	619	3569	3559	3579
195	609	599	619	3566	3556	3576
195	613	603	623	3569	3559	3579
195	613	603	623	3566	3556	3576
195	613	603	623	3567	3557	3577
67	658	648	668	3534	3524	3544
292	602	592	612	3535	3525	3545
292	604	594	614	3536	3526	3546
292	603	593	613	3510	3500	3520
292	604	594	614	3521	3511	3531
292	605	595	615	3513	3503	3523
292	606	596	616	3523	3513	3533
292	601	591	611	3536	3526	3546
270	580	575	585	3584	3579	3589
270	582	577	587	3583	3578	3588
270	582	577	587	3588	3583	3593
270	585	580	590	3587	3582	3592
270	586	581	591	3584	3579	3589
706	665	655	675	759	749	769
706	658	648	668	761	751	771
706	661	651	671	759	749	769
519	626	616	636	1687	1677	1697
519	633	623	643	1687	1677	1697
519	642	632	652	1689	1679	1699
519	645	635	655	1697	1687	1707
519	643	633	653	1707	1697	1717
67	662	652	672	3535	3525	3545
681	663	653	673	759	749	769
407	658	648	668	644	634	654
407	655	645	665	647	637	657
407	654	644	664	633	623	643
62	695	685	705	675	665	685
62	691	681	701	686	676	696
62	698	688	708	681	671	691
62	702	692	712	683	673	693
62	701	691	711	673	663	683
62	703	693	713	685	675	695
62	706	696	716	687	677	697
29	697	687	707	695	685	705
29	701	691	711	694	684	704
29	685	675	695	698	688	708
8	685	675	695	709	699	719
0	633	623	643	855	845	865
0	649	639	659	855	845	865
0	650	640	660	854	844	864
0	643	633	653	853	843	863
0	643	633	653	847	837	857
312	663	653	673	844	834	854
312	666	656	676	842	832	852
312	662	652	672	830	820	840
4	589	579	599	830	820	840
4	576	566	586	819	809	829
4	588	578	598	836	826	846
4	583	573	593	834	824	844
4	578	568	588	827	817	837
312	608	598	618	794	784	804
312	601	591	611	805	795	815
312	606	596	616	803	793	813
312	603	593	613	798	788	808
682	606	596	616	803	793	813
683	607	598	618	794	783	803
8	603	593	613	779	769	789
0	612	602	622	778	768	788
8	593	583	603	774	764	784
321	546	536	556	598	588	608
321	548	538	558	588	578	598
295	157	147	167	105	95	115
295	159	149	169	103	93	113
295	159	149	169	105	95	115
295	160	150	170	104	94	114
295	161	151	171	105	95	115
295	161	151	171	106	96	116
295	162	152	172	105	95	115
137	263	258	268	108	103	113
137	264	259	269	102	97	107
137	264	259	269	103	98	108
137	264	259	269	108	103	113
137	264	259	269	109	104	114
137	266	261	271	108	103	113
137	266	261	271	109	104	114
137	267	262	272	104	99	109
137	267	262	272	109	104	114
137	268	263	273	101	96	106
137	268	263	273	104	99	109
45	66	61	71	287	282	292
45	68	63	73	283	278	288
45	80	75	85	285	280	290
45	83	78	88	279	274	284
45	72	67	77	286	281	291
45	71	66	76	281	276	286
45	101	96	106	277	272	282
45	98	93	103	273	268	278
45	96	91	101	281	276	286
45	121	116	126	261	256	266
45	115	110	120	260	255	265
45	119	114	124	268	263	273
45	136	131	141	264	259	269
45	136	131	141	271	266	276
45	131	126	136	268	263	273
189	58	53	63	263	258	268
189	68	63	73	262	257	267
22	60	55	65	172	167	177
22	64	59	69	183	178	188
22	76	71	81	176	171	181
47	109	104	114	415	410	420
57	316	311	321	412	407	417
57	311	306	316	408	403	413
188	238	233	243	312	307	317
53	231	226	236	288	283	293
53	239	234	244	297	292	302
292	280	275	285	3018	3013	3023
292	281	276	286	3017	3012	3022
292	281	276	286	3018	3013	3023
292	282	277	287	3014	3009	3019
190	110	105	115	206	201	211
190	107	102	112	221	216	226
190	113	108	118	220	215	225
190	275	270	280	2972	2967	2977
190	277	272	282	2971	2966	2976
190	279	274	284	2968	2963	2973
190	276	271	281	2969	2964	2974
34	216	211	221	628	623	633
34	223	218	228	629	624	634
34	222	217	227	634	629	639
45	181	176	186	202	197	207
45	179	174	184	210	205	215
74	171	166	176	161	156	166
292	282	277	287	3015	3010	3020
292	281	276	286	3020	3015	3025
102	307	302	312	570	565	575
3	310	305	315	565	560	570
116	288	286	296	543	541	551
50	217	212	222	3520	3515	3525
95	217	216	223	638	635	638
60	315	310	320	416	411	421
212	266	261	267	659	658	660
212	262	261	267	659	658	660
229	99	99	99	3537	3537	3537
296	115	110	120	311	306	316
296	118	113	123	310	305	315
296	120	115	125	312	307	317
296	107	102	112	307	302	312
296	109	104	114	312	307	317
296	119	114	124	306	301	311
298	174	169	179	3545	3540	3550
298	174	169	179	3547	3542	3552
298	176	171	181	3543	3538	3548
298	176	171	181	3548	3543	3553
298	178	173	183	3546	3541	3551
298	178	173	183	3545	3540	3550
298	173	168	178	3546	3541	3551
540	713	708	718	1453	1448	1458
540	715	710	720	1453	1448	1458
789	230	225	235	130	125	135
789	226	221	231	130	125	135
790	229	224	234	133	128	138
790	227	222	232	127	122	132
791	226	221	231	131	126	136
791	229	224	234	128	123	133
712	446	441	451	3371	3366	3376
784	470	465	475	3387	3382	3392
784	473	468	478	3387	3382	3392
201	153	143	163	190	180	200
201	135	125	145	201	191	211
201	134	124	144	197	187	207
251	170	165	175	395	390	400
251	168	163	173	392	387	397
546	532	527	537	754	749	759
559	407	402	412	760	755	765
559	402	397	407	758	753	763
559	401	396	406	756	751	761
559	398	393	403	759	754	764
559	402	397	407	754	749	759
452	651	646	656	585	580	590
323	609	604	614	1548	1543	1553
323	609	604	614	1550	1545	1555
323	611	606	616	1549	1544	1554
323	612	607	617	1550	1545	1555
436	586	581	591	606	601	611
546	532	527	537	754	749	759
559	407	402	412	760	755	765
559	402	397	407	758	753	763
559	401	396	406	756	751	761
559	398	393	403	759	754	764
559	402	397	407	754	749	759
559	401	396	406	749	744	754
559	402	397	407	745	740	750
559	402	397	407	741	736	746
559	406	401	411	739	734	744
559	409	404	414	747	742	752
558	410	405	415	751	746	756
558	408	403	413	743	738	748
560	396	391	401	740	735	745
561	409	404	414	753	747	757
564	708	708	708	509	509	509
569	390	385	395	754	749	759
572	402	402	402	461	461	461
570	58	58	58	503	503	503
545	704	699	709	3283	3278	3288
429	616	611	621	3452	3447	3457
427	582	576	586	3419	3413	3423
573	584	579	589	3575	3570	3580
270	581	576	586	3555	3550	3560
270	585	580	590	3555	3550	3560
270	583	578	588	3552	3547	3557
451	622	617	627	590	585	595
444	613	608	618	580	575	585
469	629	624	634	595	590	600
29	630	625	635	593	588	598
29	635	630	640	586	581	591
29	625	620	630	594	589	599
492	631	626	636	573	568	578
29	631	626	636	578	573	583
502	628	623	633	570	565	575
502	634	629	639	569	564	574
502	628	623	633	566	561	571
502	630	625	635	565	560	570
524	662	657	667	552	547	557
527	661	656	666	543	538	548
526	657	652	662	533	528	538
233	279	274	284	296	291	301
706	663	658	668	736	731	741
312	663	658	668	747	742	752
34	539	534	544	3373	3368	3378
34	540	535	545	3370	3365	3375
34	541	536	546	3373	3368	3378
263	552	547	557	3355	3350	3360
263	562	557	567	3352	3347	3357
263	562	557	567	3356	3351	3361
364	163	155	165	468	463	473
720	61	56	66	731	726	736
720	64	59	69	730	725	735
549	63	58	68	731	726	736
719	63	63	63	732	732	732
668	83	78	88	3630	3625	3635
671	82	77	87	3623	3618	3628
671	83	78	88	3617	3612	3622
671	85	80	90	3618	3613	3623
671	88	83	93	3625	3620	3630
703	170	165	175	794	789	799
692	70	70	70	3625	3625	3625
692	70	70	70	3627	3627	3627
671	57	52	62	3622	3617	3627
671	70	65	75	3610	3605	3615
671	73	68	78	3607	3602	3612
671	60	55	65	3613	3608	3618
690	72	72	72	3637	3637	3637
668	95	90	100	806	801	811
668	93	88	98	807	802	812
668	93	88	98	808	803	813
668	94	89	99	806	801	811
668	95	90	100	810	805	815
668	93	88	98	808	803	813
668	83	78	88	806	801	811
668	80	75	85	804	799	809
671	85	80	90	811	806	816
670	83	78	88	3639	3634	3644
716	91	86	96	3638	3633	3643
670	90	85	95	806	801	811
670	90	85	95	809	804	814
670	90	90	90	806	806	806
670	90	90	90	809	809	809
342	193	188	198	109	104	114
342	192	187	197	126	121	131
251	170	165	175	393	388	398
251	174	169	179	396	391	401
68	171	166	176	317	312	322
68	179	174	184	314	309	319
68	179	174	184	317	312	322
68	185	180	190	315	310	320
68	178	173	183	307	302	312
41	176	171	181	317	312	322
41	178	173	183	315	310	320
41	178	173	183	307	302	312
41	168	163	173	315	310	320
200	363	358	368	465	460	470
200	359	354	364	463	458	468
200	366	361	371	462	457	467
200	365	360	370	460	455	465
200	361	356	366	460	455	465
99	157	152	162	256	251	261
99	149	144	154	256	251	261
99	151	146	156	251	246	256
99	155	150	160	259	254	264
99	147	142	152	261	256	266
99	171	166	176	259	254	264
99	176	171	181	256	251	261
99	173	168	178	264	259	269
99	170	165	175	254	249	259
99	166	161	171	261	256	266
99	174	169	179	249	244	254
99	178	173	183	249	244	254
99	164	159	169	252	247	257
99	155	150	160	250	245	255
99	148	143	153	250	245	255
99	148	143	153	261	256	266
99	148	143	153	256	251	261
99	156	151	161	261	256	266
99	174	169	179	262	257	267
99	159	154	164	257	252	262
99	162	157	167	259	254	264
99	161	156	166	253	248	258
99	156	151	161	253	248	258
99	157	152	162	256	251	261
320	616	611	621	536	531	541
320	619	614	624	538	533	543
319	576	571	581	525	520	530
4	580	575	585	526	521	531
4	577	572	582	525	520	530
4	577	572	582	522	517	527
4	579	574	584	527	522	532
47	113	108	118	416	411	421
47	109	104	114	412	407	417
47	103	98	108	414	409	419
47	113	108	118	415	410	420
232	269	264	274	299	294	304
232	270	265	275	297	292	302
232	263	258	268	299	294	304
232	270	265	275	319	314	324
47	257	252	262	303	298	308
47	256	251	261	305	300	310
203	357	359	369	3352	3347	3357
203	358	359	369	3353	3348	3358
203	362	359	369	3349	3344	3354
203	363	359	369	3353	3348	3358
203	363	359	369	3357	3352	3362
203	369	364	374	3372	3367	3377
203	372	367	377	3372	3367	3377
367	594	589	599	3476	3471	3481
358	559	552	562	649	645	655
63	266	261	271	600	595	605
57	311	306	316	409	404	414
324	548	543	553	599	594	604
62	115	110	120	632	627	637
62	119	114	124	632	627	637
62	115	110	120	629	624	634
62	110	105	115	630	625	635
62	116	111	121	637	632	642
62	112	107	117	640	635	645
62	104	99	109	644	639	649
62	99	94	104	654	649	659
62	105	100	110	651	646	656
30	133	128	138	507	502	512
240	357	352	362	489	484	494
409	650	645	655	659	654	664
409	658	653	663	660	655	665
409	652	647	657	660	655	665
409	649	644	654	664	659	669
409	660	655	665	660	655	665
409	663	658	668	664	659	669
409	648	643	653	661	656	666
409	651	646	656	668	663	673
409	657	652	662	668	663	673
409	655	650	660	666	661	671
409	649	644	654	664	659	669
320	608	603	613	599	594	604
320	608	603	613	607	602	612
320	605	600	610	598	593	603
320	605	600	610	608	603	613
433	603	603	603	587	587	587
244	351	345	355	490	485	495
264	448	443	453	700	695	705
264	448	443	453	707	702	712
311	522	517	527	3433	3428	3438
311	521	516	526	3437	3432	3442
47	103	98	108	414	409	419
525	667	662	672	534	529	539
525	664	659	669	534	529	539
525	664	659	669	531	526	536
525	667	662	672	531	526	536
525	666	661	671	533	528	538
239	397	392	402	471	466	476
239	395	390	400	467	462	472
655	759	754	764	3449	3444	3454
655	759	754	764	3450	3445	3455
655	761	756	766	3449	3444	3454
658	776	771	781	3456	3451	3461
658	779	774	784	3455	3450	3460
658	786	781	791	3469	3464	3474
658	786	781	791	3469	3464	3474
658	785	780	790	3469	3464	3474
658	789	784	794	3468	3463	3473
658	788	783	793	3470	3465	3475
658	788	783	793	3468	3463	3473
658	792	787	797	3468	3463	3473
658	792	787	797	3474	3469	3479
658	793	788	798	3477	3472	3482
658	796	791	801	3476	3471	3481
658	792	787	797	3477	3472	3482
658	804	799	809	3475	3470	3480
658	803	798	808	3479	3474	3484
658	805	800	810	3464	3459	3469
655	807	802	812	3467	3462	3472
658	802	797	807	3461	3456	3466
658	795	790	800	3461	3456	3466
658	792	787	797	3460	3455	3465
658	789	784	794	3465	3460	3470
658	795	790	800	3463	3458	3468
655	795	790	800	3457	3452	3462
658	776	771	781	3482	3477	3487
658	778	773	783	3482	3477	3487
655	770	765	775	3492	3487	3497
655	772	767	777	3492	3487	3497
655	777	772	782	3494	3489	3499
655	784	779	789	3497	3492	3502
655	782	777	787	3550	3545	3555
644	768	763	773	3549	3544	3554
644	768	763	773	3546	3541	3551
643	773	768	778	3537	3532	3542
655	774	769	779	3527	3522	3532
655	769	764	774	3527	3522	3532
655	769	764	774	3526	3521	3531
655	785	780	790	3529	3524	3534
655	783	778	788	3529	3524	3534
646	807	802	812	3542	3538	3548
647	807	802	812	3529	3525	3535
655	795	790	800	3519	3514	3524
655	776	771	781	3510	3505	3515
655	770	765	775	3510	3505	3515
655	770	765	775	3510	3505	3515
655	795	790	800	3509	3504	3514
655	814	809	819	3512	3507	3517
655	814	809	819	3510	3505	3515
655	815	810	820	3514	3509	3519
655	809	804	814	3496	3491	3501
644	815	810	820	3500	3495	3505
644	815	810	820	3503	3498	3508
655	795	790	800	3490	3485	3495
655	781	776	786	3435	3430	3440
655	782	777	787	3444	3439	3449
655	803	798	808	3446	3441	3451
655	803	798	808	3444	3439	3449
644	815	810	820	3457	3452	3462
644	815	810	820	3455	3450	3460
644	801	801	801	3435	3435	3435
644	792	792	792	3435	3435	3435
644	788	788	788	3431	3431	3431
644	791	791	791	3426	3426	3426
644	791	791	791	3420	3420	3420
644	796	796	796	3414	3414	3414
644	791	791	791	3410	3410	3410
644	798	798	798	3409	3409	3409
644	802	802	802	3414	3414	3414
644	802	802	802	3420	3420	3420
644	802	802	802	3426	3426	3426
644	797	797	797	3431	3431	3431
644	797	797	797	3424	3424	3424
644	808	808	808	3429	3429	3429
644	812	812	812	3426	3426	3426
644	808	808	808	3434	3434	3434
644	808	808	808	3422	3422	3422
644	808	808	808	3414	3414	3414
644	810	810	810	3408	3408	3408
644	813	813	813	3408	3408	3408
655	783	783	783	3419	3419	3419
34	755	750	760	640	635	645
34	764	759	769	643	638	648
34	761	756	766	647	642	652
34	753	748	758	645	640	650
34	735	730	740	643	638	648
34	729	724	734	638	633	643
34	723	718	728	647	642	652
34	729	724	734	659	654	664
43	728	723	733	656	651	661
43	731	726	736	658	653	663
43	730	725	735	653	648	658
43	736	731	741	657	652	662
43	733	728	738	664	659	669
43	734	729	739	656	651	661
43	735	730	740	651	646	656
43	740	735	745	664	659	669
43	743	738	748	669	664	674
34	746	741	751	662	657	667
34	751	746	756	660	655	665
34	755	750	760	657	652	662
34	756	751	761	653	648	658
34	764	759	769	652	647	657
43	729	724	734	663	658	668
43	737	732	742	652	647	657
43	743	738	748	658	653	663
34	725	720	730	640	635	645
34	724	719	729	616	611	621
655	724	719	729	613	608	618
655	736	731	741	595	590	600
34	731	726	736	598	593	603
34	736	731	741	588	583	593
34	739	734	744	589	584	594
34	742	737	747	588	583	593
34	745	740	750	597	592	602
34	747	742	752	598	593	603
34	743	738	748	594	589	599
650	740	735	745	584	579	589
649	804	802	806	3469	3467	3471
645	795	792	796	3541	3540	3544
659	676	674	678	3493	3491	3495
631	710	708	712	3421	3419	3423
631	715	713	717	3422	3420	3424
631	712	710	714	3435	3433	3437
631	708	706	710	3438	3436	3440
630	764	762	766	3469	3467	3471
630	765	763	767	3471	3469	3473
630	765	763	767	3468	3466	3470
630	764	762	766	3468	3466	3470
630	764	762	766	3467	3465	3469
630	764	762	766	3467	3465	3469
630	764	762	766	3470	3468	3472
630	764	762	766	3466	3464	3468
630	764	762	766	3469	3467	3471
43	752	750	754	3477	3475	3479
43	746	744	748	3477	3475	3479
43	749	747	751	3480	3478	3482
630	750	748	752	3479	3477	3481
630	749	747	751	3480	3478	3482
630	749	747	751	3480	3478	3482
630	750	748	752	3480	3478	3482
630	747	745	749	3474	3472	3476
630	744	742	746	3474	3472	3476
630	744	742	746	3476	3474	3478
630	744	742	746	3476	3474	3478
43	747	745	749	3474	3472	3476
43	742	740	744	3481	3479	3483
43	740	738	742	3479	3477	3481
630	733	731	735	3483	3481	3485
630	737	735	739	3483	3481	3485
43	740	738	742	3475	3473	3477
43	736	734	738	3474	3472	3476
43	736	734	738	3479	3477	3481
43	734	732	736	3475	3473	3477
43	730	728	732	3475	3473	3477
630	729	727	731	3478	3476	3480
46	745	743	747	3496	3494	3498
46	743	741	745	3496	3494	3498
46	748	746	750	3498	3496	3500
45	746	744	748	3497	3495	3499
45	741	739	743	3497	3495	3499
0	744	742	746	3500	3498	3502
667	665	663	667	567	565	569
12	120	118	122	521	519	523
515	606	604	608	3585	3583	3587
358	558	556	560	1592	1590	1594
369	604	602	606	501	499	503
370	603	601	605	502	500	504
781	123	122	124	506	505	507
283	353	352	354	523	522	524
661	684	681	686	589	586	592
729	672	669	677	571	568	572
731	687	684	690	567	565	569
732	689	684	694	567	563	571
734	319	319	320	742	741	743
734	271	271	272	742	741	743
734	252	251	252	729	728	730
391	553	551	555	708	706	710
575	77	75	79	1643	1641	1645
196	418	414	422	3479	3475	3485
189	257	250	272	187	177	196
189	254	250	272	184	177	196
189	252	250	272	190	177	196
136	241	235	247	195	188	200
136	236	234	250	189	180	196
136	239	239	251	199	192	204
136	235	227	237	199	192	204
11	288	285	291	1467	1464	1470
95	330	328	332	555	553	557
95	334	328	334	555	551	557
95	331	328	334	554	549	557
67	362	358	366	611	607	615
67	362	360	364	611	609	613
67	365	361	369	604	600	608
95	151	150	152	501	500	502
95	283	282	284	568	567	569
95	219	218	220	636	635	637
95	101	100	102	512	511	513
95	90	89	91	694	693	695
95	331	330	332	553	552	554
95	369	368	370	715	714	716
673	653	645	661	3622	3614	3630
673	642	634	650	3622	3614	3630
673	660	652	668	3619	3611	3627
673	633	625	641	3615	3607	3623
673	655	647	663	3612	3604	3620
673	648	640	656	3607	3599	3615
499	224	222	226	741	739	743
499	222	220	224	738	736	740
494	219	217	221	761	759	763
496	232	230	234	760	758	762
497	227	225	229	754	752	756
497	230	228	232	757	755	759
199	257	249	265	345	337	353
199	256	247	265	360	351	369
199	261	252	270	353	344	362
251	170	167	173	401	398	404
251	169	164	174	397	392	402
47	109	107	111	416	414	418
47	114	112	116	416	414	418
47	110	108	112	417	415	419
47	112	110	114	419	417	421
57	314	312	316	414	412	416
57	310	308	312	409	407	411
57	312	310	314	409	407	411
60	312	310	314	414	412	416
95	73	71	75	1643	1641	1645
83	135	133	137	640	638	642
95	217	215	219	450	448	452
3	64	62	66	1639	1637	1641
6	64	62	66	1643	1641	1645
47	65	63	67	1640	1638	1642
11	64	62	66	1641	1639	1643
67	64	62	66	1642	1640	1644
61	64	62	66	1639	1637	1641
86	65	63	67	1643	1641	1645
89	486	474	498	431	419	443
6	496	490	502	406	400	412
383	615	607	625	1661	1651	1669
243	667	665	669	851	849	853
243	659	657	661	852	850	854
243	663	661	665	846	844	848
312	631	629	633	784	782	786
312	629	627	631	783	781	785
691	634	632	636	797	795	799
691	635	633	637	793	791	795
691	636	634	638	793	791	795
312	635	633	637	807	805	809
686	644	642	646	791	789	793
704	653	651	655	790	788	792
704	663	661	665	791	789	793
704	663	661	665	783	781	785
704	655	653	657	783	781	785
704	656	654	658	789	787	791
689	649	647	651	790	788	792
689	650	648	652	785	783	787
687	652	650	654	791	789	793
531	660	658	662	782	780	784
684	662	662	662	789	789	789
676	664	662	666	774	772	776
676	665	663	667	772	770	774
312	653	651	655	776	774	778
312	655	653	657	776	774	778
677	669	666	672	811	808	814
677	666	663	669	811	808	814
688	662	658	666	803	799	807
312	656	653	659	802	799	805
312	658	654	662	803	799	807
312	650	645	655	807	802	812
697	646	643	649	806	803	809
697	648	645	651	806	803	809
679	647	645	649	3596	3594	3598
679	650	648	652	3594	3592	3596
43	627	625	629	3591	3589	3593
43	630	628	632	3594	3592	3596
709	628	626	630	3592	3590	3594
709	631	629	633	3593	3591	3595
696	650	648	652	3556	3554	3558
696	650	648	652	3555	3553	3557
696	652	650	654	3558	3556	3560
707	637	635	639	3567	3565	3569
707	640	638	642	3564	3562	3566
708	626	624	628	3565	3563	3567
708	628	626	630	3568	3566	3570
678	627	625	629	3581	3579	3583
678	630	628	632	3574	3572	3576
678	626	624	628	3576	3574	3578
34	649	640	658	3628	3619	3637
531	639	631	647	3623	3615	3631
531	636	631	641	3619	3614	3624
202	635	630	640	3616	3611	3621
40	636	630	642	3611	3605	3617
531	640	636	644	3609	3605	3613
202	650	645	655	3608	3603	3613
34	644	638	650	3607	3601	3613
34	651	646	656	3607	3602	3612
40	656	649	663	3613	3606	3620
202	659	654	664	3617	3612	3622
531	657	650	664	3618	3611	3625
531	658	651	665	3623	3616	3630
43	665	657	673	3628	3620	3636
43	666	659	673	3630	3623	3637
290	664	656	672	3632	3624	3640
43	665	658	673	3629	3622	3636
290	666	657	675	3633	3624	3642
34	650	643	657	3628	3621	3635
672	494	492	496	3522	3520	3524
672	492	490	494	3521	3519	3523
672	492	490	494	3524	3522	3526
616	395	393	397	837	835	839
777	467	466	468	877	876	878
777	466	465	467	878	877	879
777	466	465	467	876	875	877
775	467	466	468	877	876	878
775	466	465	467	877	876	878
777	470	469	471	870	869	871
51	127	126	128	516	515	517
766	454	453	455	3708	3707	3709
308	81	80	83	666	666	668
738	25	22	30	3345	3339	3347
738	25	22	30	3346	3339	3347
37	105	102	108	1479	1476	1480
155	332	331	333	567	566	568
768	441	438	444	3719	3718	3726
768	439	438	444	3721	3719	3728
768	444	438	444	3725	3720	3730
768	444	438	444	3727	3721	3732
768	439	438	444	3728	3722	3734
768	444	438	444	3731	3723	3736
768	444	438	444	3733	3724	3738
768	451	446	454	3738	3724	3740
663	428	427	429	3711	3710	3712
663	425	424	426	3710	3709	3711
663	418	417	419	3709	3708	3710
663	418	417	419	3712	3711	3713
761	408	407	409	3712	3711	3713
761	404	403	405	3709	3708	3710
761	405	404	406	3704	3703	3705
761	405	404	406	3699	3698	3700
762	400	399	401	3719	3718	3720
762	396	395	397	3719	3718	3720
762	392	391	393	3721	3720	3722
762	393	392	394	3724	3723	3725
22	391	390	392	3732	3731	3733
22	393	392	394	3736	3735	3737
22	394	393	395	3741	3740	3742
22	400	399	401	3737	3736	3738
23	419	418	420	3735	3734	3736
23	421	420	422	3741	3740	3742
23	421	420	422	3725	3724	3726
23	423	422	424	3728	3727	3729
175	310	309	311	3346	3345	3347
785	514	508	517	535	533	539
114	200	179	207	457	441	467
734	300	299	300	729	728	730
607	742	741	743	452	451	453
608	741	740	742	455	454	456
603	740	739	741	446	445	447
600	735	734	736	448	447	449
595	735	734	736	451	450	452
604	737	736	738	451	450	452
602	734	733	735	454	453	455
597	731	730	732	446	445	447
599	731	730	732	454	453	455
598	733	732	734	450	449	451
611	736	735	737	445	444	446
611	735	734	736	444	443	445
611	737	736	738	443	442	444
611	738	737	739	456	455	457
611	736	735	737	456	455	457
63	121	111	131	611	601	621
34	140	134	146	2535	2529	2541
34	137	130	144	3488	3481	3495
34	140	133	147	3491	3484	3498
11	133	123	143	1602	1592	1612
29	140	135	145	667	662	672
29	138	133	143	667	662	672
5	137	127	147	651	641	661
11	118	108	128	658	648	668
11	126	116	136	656	646	666
11	129	119	139	658	648	668
114	127	112	142	660	645	675
11	129	119	139	649	639	659
11	132	122	142	639	629	649
62	129	119	139	629	619	639
62	124	114	134	664	654	674
62	96	86	106	664	654	674
62	98	88	108	662	652	672
23	120	110	130	636	626	646
23	111	101	121	640	630	650
62	103	93	113	636	626	646
62	120	110	130	625	615	635
62	100	90	110	633	623	643
62	103	93	113	630	620	640
62	105	95	115	646	636	656
62	113	103	123	631	621	641
62	116	106	126	632	622	642
62	114	104	124	630	620	640
62	116	106	126	629	619	639
11	97	87	107	579	569	589
4	101	91	111	572	562	582
717	63	53	73	725	715	735
717	63	53	73	727	717	737
717	61	51	71	726	716	736
79	51	41	61	717	707	727
34	53	43	63	717	707	727
34	52	42	62	715	705	725
34	49	39	59	715	705	725
70	87	77	97	715	705	725
86	76	66	86	694	684	704
86	66	56	76	691	681	701
86	68	58	78	694	684	704
86	69	59	79	693	683	703
86	78	68	88	689	679	699
86	78	68	88	691	681	701
86	66	56	76	689	679	699
86	66	56	76	686	676	696
86	69	59	79	694	684	704
6	50	40	60	693	683	703
29	57	47	67	687	677	697
29	57	47	67	686	676	696
70	69	59	79	592	582	602
70	73	63	83	589	579	599
70	70	60	80	585	575	595
70	69	59	79	600	590	610
114	69	59	79	608	598	618
8	73	58	88	561	546	576
0	83	73	93	557	547	567
0	86	76	96	549	539	559
57	117	112	122	555	550	560
57	116	111	121	550	545	555
57	112	107	117	557	552	562
57	105	100	110	549	544	554
57	109	104	114	546	541	551
57	105	100	110	553	548	558
60	108	103	113	557	552	562
60	113	108	118	546	541	551
65	131	126	136	536	531	541
65	134	129	139	537	532	542
3	145	140	150	561	556	566
3	145	140	150	560	555	565
3	146	141	151	562	557	567
3	147	142	152	560	555	565
3	145	140	150	562	557	567
34	149	144	154	553	548	558
21	162	152	172	548	538	558
25	121	116	126	525	520	530
63	115	105	125	595	585	605
65	149	144	154	507	502	512
65	155	150	160	508	503	513
65	152	147	157	507	502	512
65	145	140	150	507	502	512
21	105	100	110	528	523	533
47	95	90	100	535	530	540
65	87	82	92	507	502	512
65	91	86	96	509	504	514
65	89	84	94	508	503	513
65	87	82	92	514	509	519
66	123	118	128	437	432	442
64	81	76	86	449	444	454
63	80	75	85	444	439	449
11	82	77	87	444	439	449
66	85	75	95	447	437	457
114	97	82	112	438	423	453
114	114	99	129	447	432	462
65	109	104	114	464	459	469
65	113	108	118	463	458	468
64	111	106	116	3380	3375	3385
64	111	106	116	3373	3368	3378
64	112	107	117	3371	3366	3376
64	110	105	115	3371	3366	3376
65	128	123	133	483	478	488
65	133	128	138	483	478	488
65	136	131	141	484	479	489
65	141	136	146	456	451	461
159	136	131	141	461	456	466
159	139	134	144	462	457	467
65	138	133	143	1397	1392	1402
65	140	135	145	1395	1390	1400
65	142	137	147	1400	1395	1405
65	140	135	145	1406	1401	1411
114	154	144	164	522	512	532
114	110	100	120	530	520	540
114	130	115	145	451	436	466
114	142	127	157	477	462	492
202	375	370	380	3272	3267	3277
6	99	94	104	626	621	631
6	102	97	107	627	622	632
6	97	92	102	628	623	633
6	94	89	99	620	615	625
6	102	97	107	614	609	619
6	97	92	102	616	611	621
6	108	103	113	610	605	615
6	101	96	106	607	602	612
6	110	105	115	606	601	611
6	98	93	103	607	602	612
3	120	115	125	603	598	608
3	118	113	123	603	598	608
3	118	113	123	604	599	609
3	120	115	125	605	600	610
3	120	115	125	604	599	609
3	120	115	125	603	598	608
63	110	100	120	565	555	575
8	117	112	122	450	445	455
29	114	109	119	3306	3301	3311
29	112	107	117	3306	3301	3311
29	109	104	114	3306	3301	3311
19	94	89	99	3301	3296	3306
19	119	114	124	3298	3293	3303
52	101	91	111	3289	3279	3299
52	107	97	117	3291	3281	3301
70	104	99	109	3280	3275	3285
70	101	96	106	3280	3275	3285
47	90	85	95	3289	3284	3294
47	91	86	96	3290	3285	3295
40	95	90	100	3276	3271	3281
40	94	89	99	3276	3271	3281
40	93	88	98	3275	3270	3280
40	102	97	107	3276	3271	3281
40	103	98	108	3277	3272	3282
53	118	108	128	3274	3264	3284
53	116	106	126	3275	3265	3285
53	117	107	127	3273	3263	3283
41	121	116	126	3285	3280	3290
41	123	118	128	3288	3283	3293
41	128	123	133	3291	3286	3296
41	120	115	125	3293	3288	3298
45	138	133	143	3285	3280	3290
45	137	132	142	3286	3281	3291
34	137	132	142	3292	3287	3297
34	137	132	142	3289	3284	3294
34	135	130	140	3294	3289	3299
74	137	132	142	3295	3290	3300
99	161	156	166	3297	3292	3302
99	146	141	151	3293	3288	3298
99	152	147	157	3297	3292	3302
99	157	152	162	3299	3294	3304
99	154	149	159	3292	3287	3297
99	148	143	153	3295	3290	3300
104	164	159	169	3290	3285	3295
104	162	157	167	3306	3301	3311
8	165	155	175	439	429	449
29	159	154	164	468	463	473
29	161	156	166	466	461	471
29	160	150	170	463	453	473
78	233	223	243	499	489	509
11	117	107	127	1468	1458	1478
11	111	106	116	1477	1472	1482
11	108	103	113	527	522	532
11	97	92	102	526	521	531
11	110	105	115	523	518	528
11	112	107	117	522	517	527
11	124	119	129	522	517	527
64	113	103	123	531	521	541
76	233	223	243	502	492	512
76	235	225	245	501	491	511
76	235	225	245	499	489	509
76	233	223	243	506	496	516
76	230	220	240	509	499	519
76	230	220	240	522	512	532
76	236	226	246	519	509	529
0	227	217	237	482	472	492
114	235	225	245	447	437	457
114	230	220	240	479	469	489
114	225	210	240	451	436	466
11	212	207	217	441	436	446
11	216	211	221	441	436	446
11	213	208	218	443	438	448
11	215	210	220	444	439	449
65	206	201	211	436	431	441
65	206	201	211	435	430	440
65	203	198	208	433	428	438
65	202	197	207	438	433	443
23	215	207	223	3295	3287	3303
47	210	205	215	3285	3280	3290
47	207	202	212	3286	3281	3291
47	222	217	227	3282	3277	3287
40	217	212	222	3284	3279	3289
40	217	212	222	3281	3276	3286
40	216	211	221	3274	3269	3279
40	213	208	218	3273	3268	3278
46	195	190	200	3271	3266	3276
46	197	192	202	3273	3268	3278
41	180	172	188	3291	3283	3299
41	183	175	191	3286	3278	3294
41	179	171	187	3278	3270	3286
41	183	175	191	3278	3270	3286
45	189	174	204	3300	3285	3315
45	188	173	203	3297	3282	3312
45	191	176	206	3295	3280	3310
67	209	199	219	3312	3302	3322
67	207	197	217	3313	3303	3323
67	200	190	210	3298	3288	3308
67	205	195	215	3299	3289	3309
67	206	196	216	3294	3284	3304
68	201	191	211	3307	3297	3317
68	199	189	209	3305	3295	3315
61	208	203	213	3320	3315	3325
61	210	205	215	3322	3317	3327
61	206	201	211	3326	3321	3331
61	207	202	212	3329	3324	3334
61	201	196	206	3333	3328	3338
61	208	203	213	3328	3323	3333
251	198	193	203	3257	3252	3262
251	200	195	205	3256	3251	3261
251	200	195	205	3253	3248	3258
251	198	193	203	3253	3248	3258
251	197	192	202	3251	3246	3256
584	207	202	212	3210	3205	3215
584	208	203	213	3219	3214	3224
584	206	201	211	3218	3213	3223
584	206	201	211	3212	3207	3217
584	207	202	212	3208	3203	3213
584	205	200	210	3208	3203	3213
584	207	202	212	3207	3202	3212
99	206	201	211	3239	3234	3244
99	204	199	209	3234	3229	3239
99	202	197	207	3238	3233	3243
99	204	199	209	3237	3232	3242
99	207	202	212	3236	3231	3241
99	205	200	210	3238	3233	3243
195	213	208	218	3233	3228	3238
195	215	210	220	3239	3234	3244
195	219	214	224	3236	3231	3241
195	215	210	220	3235	3230	3240
195	218	213	223	3234	3229	3239
195	218	213	223	3232	3227	3237
195	211	206	216	3233	3228	3238
195	216	211	221	3233	3228	3238
292	230	220	240	3240	3230	3250
292	231	221	241	3239	3229	3249
292	231	221	241	3240	3230	3250
290	227	222	232	3231	3226	3236
290	229	224	234	3232	3227	3237
290	233	228	238	3232	3227	3237
476	218	217	219	747	746	748
476	219	218	220	743	742	744
476	217	216	218	746	745	747
474	229	228	230	736	735	737
473	230	229	231	734	733	735
473	228	227	229	731	730	732
774	225	224	226	761	760	762
477	567	557	577	3318	3308	3328
176	250	245	255	1410	1405	1415
93	263	253	273	1404	1394	1414
93	261	251	271	1412	1402	1422
93	261	251	271	1409	1399	1419
93	261	251	271	1403	1393	1413
93	251	241	261	1403	1393	1413
93	251	241	261	1411	1401	1421
699	270	265	275	485	480	490
699	275	270	280	489	484	494
699	272	267	277	494	489	499
699	268	263	273	486	481	491
94	280	275	285	492	487	497
94	269	259	279	496	486	506
94	267	257	277	501	491	511
81	222	207	237	3517	3502	3532
81	221	216	226	692	687	697
81	217	212	222	696	691	701
81	212	207	217	692	687	697
81	213	208	218	697	692	702
81	212	207	217	1640	1635	1645
81	212	207	217	1635	1630	1640
89	198	188	208	604	594	614
139	331	326	336	435	430	440
139	333	328	338	437	432	442
140	330	325	335	436	431	441
140	331	326	336	433	428	438
140	332	327	337	434	429	439
140	331	326	336	437	432	442
265	360	355	365	3426	3421	3431
237	271	261	281	300	290	310
238	270	260	280	296	286	306
236	268	263	273	309	304	314
567	583	578	588	3554	3549	3559
516	610	605	615	3583	3578	3588
516	607	602	612	3583	3578	3588
516	605	600	610	3584	3579	3589
516	607	602	612	3586	3581	3591
252	404	394	414	532	522	542
555	580	575	585	3587	3582	3592
94	361	356	366	3424	3419	3429
94	360	355	365	3422	3417	3427
295	328	318	338	3362	3352	3372
295	329	319	339	3369	3359	3379
66	220	210	230	652	642	662
29	220	210	230	621	611	631
29	219	209	229	622	612	632
29	215	210	220	612	607	617
114	237	222	252	631	616	646
62	213	203	223	617	607	627
3	158	153	163	616	611	621
3	157	152	162	615	610	620
3	159	154	164	615	610	620
3	161	156	166	616	611	621
3	160	155	165	616	611	621
6	162	152	172	581	571	591
6	161	151	171	583	573	593
6	174	164	184	565	555	575
6	170	160	180	588	578	598
6	150	140	160	580	570	590
6	154	144	164	588	578	598
6	146	136	156	583	573	593
55	134	131	137	643	640	646
8	183	173	193	653	643	663
8	169	159	179	653	643	663
0	185	175	195	661	651	671
23	173	168	178	637	632	642
4	177	167	187	654	644	664
4	191	181	201	669	659	679
4	162	152	172	672	662	682
4	187	177	197	605	595	615
4	170	160	180	616	606	626
23	195	190	200	567	562	572
62	208	203	213	589	584	594
23	233	228	238	554	549	559
11	276	271	281	639	634	644
64	277	272	282	642	637	647
64	276	271	281	642	637	647
11	214	209	219	1563	1558	1568
29	218	213	223	555	550	560
79	218	213	223	544	539	549
34	214	209	219	544	539	549
127	199	194	204	635	630	640
127	196	191	201	632	627	637
127	193	188	198	632	627	637
127	192	187	197	634	629	639
57	229	224	234	642	637	647
57	230	225	235	641	636	646
19	132	127	137	687	682	692
19	137	132	142	686	681	691
19	136	131	141	694	689	699
19	131	126	136	691	686	696
348	561	556	566	475	470	480
348	558	553	563	475	470	480
348	565	560	570	474	469	479
21	283	278	288	3544	3539	3549
21	282	277	287	3543	3538	3548
137	276	271	281	3523	3518	3528
137	282	277	287	3520	3515	3525
137	275	270	280	3527	3522	3532
137	282	277	287	3527	3522	3532
67	294	284	304	3523	3513	3533
67	292	282	302	3514	3504	3524
67	292	282	302	3523	3513	3533
67	290	280	300	3513	3503	3523
158	306	301	311	3521	3516	3526
158	302	297	307	3516	3511	3521
158	306	301	311	3517	3512	3522
158	301	296	306	3519	3514	3524
135	311	306	316	3524	3519	3529
135	310	305	315	3521	3516	3526
29	324	314	334	670	660	680
29	324	314	334	668	658	678
29	327	317	337	670	660	680
29	327	317	337	668	658	678
57	361	356	366	2459	2454	2464
60	362	357	367	2458	2453	2463
57	360	355	365	2460	2455	2465
57	360	355	365	2457	2452	2462
57	362	357	367	1517	1512	1522
60	362	357	367	1514	1509	1519
60	360	355	365	572	567	577
57	362	357	367	571	566	576
60	362	357	367	571	566	576
60	361	356	366	570	565	575
6	354	349	359	604	599	609
6	356	351	361	605	600	610
6	351	346	356	604	599	609
6	352	347	357	607	602	612
6	352	347	357	611	606	616
6	350	345	355	612	607	617
6	349	344	354	617	612	622
6	344	339	349	613	608	618
6	345	340	350	616	611	621
65	290	287	293	579	576	582
65	291	288	294	580	577	583
65	314	311	317	521	518	524
65	316	313	319	521	518	524
11	308	303	313	523	518	528
11	310	305	315	522	517	527
11	309	304	314	544	539	549
11	314	309	319	548	543	553
11	304	299	309	540	535	545
65	312	309	315	539	536	542
65	315	312	318	538	535	541
65	315	312	318	541	538	544
3	272	267	277	605	600	610
3	272	267	277	601	596	606
3	272	267	277	603	598	608
3	269	264	274	606	601	611
3	271	266	276	607	602	612
80	217	207	227	1497	1487	1507
80	222	212	232	1495	1485	1505
80	223	213	233	1498	1488	1508
80	204	194	214	1496	1486	1506
80	210	200	220	1497	1487	1507
34	210	200	220	1487	1477	1497
34	204	194	214	1491	1481	1501
46	419	414	424	3525	3520	3530
46	414	409	419	3522	3517	3527
46	420	415	425	3530	3525	3535
46	426	421	431	3531	3526	3536
22	408	403	413	3535	3530	3540
22	410	405	415	3539	3534	3544
22	414	409	419	3536	3531	3541
99	408	403	413	3522	3517	3527
99	407	402	412	3525	3520	3530
99	405	400	410	3520	3515	3525
99	406	401	411	3521	3516	3526
22	404	399	409	3503	3498	3508
22	406	401	411	3499	3494	3504
22	404	399	409	3496	3491	3501
22	407	402	412	3492	3487	3497
195	413	408	418	3484	3479	3489
195	410	405	415	3485	3480	3490
195	411	406	416	3482	3477	3487
195	410	405	415	3479	3474	3484
195	413	408	418	3478	3473	3483
195	411	406	416	3476	3471	3481
45	416	411	421	3468	3463	3473
45	419	414	424	3470	3465	3475
45	418	413	423	3465	3460	3470
45	420	415	425	3467	3462	3472
45	416	411	421	3463	3458	3468
22	415	407	423	616	608	624
104	423	418	428	639	634	644
104	421	416	426	637	632	642
136	410	405	415	615	610	620
136	406	401	411	606	601	611
136	406	401	411	616	611	621
45	408	398	418	627	617	637
45	407	397	417	635	625	645
248	413	408	418	466	461	471
248	416	411	421	464	459	469
248	412	407	417	457	452	462
248	401	396	406	459	454	464
248	411	406	416	440	435	445
248	406	401	411	437	432	442
158	412	407	417	432	427	437
158	421	416	426	432	427	437
158	410	405	415	3268	3263	3273
158	418	413	423	3268	3263	3273
158	420	415	425	3284	3279	3289
158	406	401	411	3271	3266	3276
158	415	410	420	3274	3269	3279
158	417	416	422	3281	3276	3286
158	410	405	415	3265	3260	3270
158	418	413	423	3266	3261	3271
158	424	419	429	3275	3270	3280
21	446	436	456	675	665	685
21	437	427	447	703	693	713
6	612	607	617	733	728	738
6	609	604	614	737	732	742
6	611	606	616	738	733	743
530	615	612	618	750	747	753
0	434	424	444	554	544	564
277	462	457	467	521	516	526
277	464	459	469	519	514	524
277	460	455	465	519	514	524
277	464	459	469	1465	1460	1470
277	458	453	463	1460	1455	1465
277	460	455	465	1464	1459	1469
277	462	457	467	1467	1462	1472
504	344	341	347	665	662	668
506	344	341	347	657	654	660
505	340	337	343	657	654	660
507	346	343	349	658	655	661
511	92	89	95	528	525	531
510	93	90	96	529	526	532
292	260	250	270	3007	2997	3017
292	262	252	272	3003	2993	3013
292	265	255	275	3002	2992	3012
292	269	259	279	3010	3000	3020
292	271	261	281	3005	2995	3015
292	273	263	283	3010	3000	3020
362	542	537	547	3278	3273	3283
362	543	538	548	3275	3270	3280
362	548	543	553	3278	3273	3283
363	542	537	547	3275	3270	3280
363	547	542	552	3275	3270	3280
363	545	540	550	3279	3274	3284
363	545	540	550	3280	3275	3285
\.


--
-- Data for Name: npcs; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.npcs (id, name, description, command, hits, attack, strength, defense, hostility) FROM stdin;
1	Bob	An axe seller		2	2	2	3	0
2	Sheep	A very wooly sheep		0	0	0	3	0
21	mugger	He jumps out and attacks people		8	15	10	8	5
41	zombie	The living dead		23	23	28	24	5
43	Giant bat	An angry flying rodent		32	32	32	32	5
7	cook	The head cook of Lumbridge castle		20	20	20	3	0
8	Bear	Eek! A bear!		26	25	23	25	1
9	Priest	A priest of Saradomin		0	0	0	3	0
10	Urhney	He looks a little grumpy		10	10	10	3	0
46	skeleton	It rattles as it walks		28	27	24	24	5
13	Camel	Oh its a camel		0	0	0	3	0
14	Gypsy	An old gypsy lady		0	0	0	3	0
15	Ghost	Ooh spooky		15	15	15	5	0
16	Sir Prysin	One of the king's knights		20	30	60	50	0
17	Traiborn the wizard	An old wizard		10	20	15	3	0
18	Captain Rovin	The head of the palace guard		30	40	70	65	0
19	Rat	Overgrown vermin		10	10	10	5	1
20	Reldo	I think he's the librarian		10	20	15	3	0
23	Giant Spider	I think this spider has been genetically modified		10	10	10	5	1
24	Man	A shifty looking man		30	30	30	30	0
25	Jonny the beard	I can see why he's called the beard		5	10	20	8	1
26	Baraek	A fur trader		30	30	30	30	0
27	Katrine	She doesn't look to friendly		30	35	25	10	0
28	Tramp	A scruffy looking chap		7	9	8	5	0
47	Rat	overgrown vermin		15	16	12	10	5
30	Romeo	He looks mildly confused		40	20	60	60	0
31	Juliet	She looks a little stressed		2	2	4	3	0
32	Father Lawrence	A kindly looking priest		0	0	0	3	0
53	Ghost	Ooh spooky		23	23	30	25	5
57	Darkwizard	He works evil magic		12	15	15	12	5
36	Veronica	She doesn't look too happy		1	1	1	5	0
60	Darkwizard	He works evil magic		27	27	24	24	5
38	Professor Oddenstein	A mad scientist if I ever saw one		3	3	3	7	0
39	Curator	He looks like he's daydreaming		2	3	2	3	0
40	skeleton	It rattles as it walks		23	24	20	17	1
42	king	King Roald the VIII		15	15	60	30	0
44	Bartender	A friendly barman		0	0	0	3	0
48	Horvik the Armourer	He looks strong		6	15	22	22	0
49	Bear	A  bear		0	0	0	3	0
61	Giant	A very large foe		40	37	36	35	5
51	Shopkeeper	Maybe he'd like to buy some of my junk		0	0	0	3	0
52	zombie	The living dead		19	18	20	22	1
55	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
56	shopkeeper	I can buy swords off him		0	0	0	3	0
58	lowe	The owner of the archery store		0	0	0	3	0
59	Thessalia	A young shop assistant		0	0	0	3	0
67	Hobgoblin	A large ugly green creature		34	32	34	29	5
68	zombie	The living dead		35	32	31	30	5
70	Scorpion	An extremely vicious scorpion		35	35	35	25	5
69	Zaff	He trades in staffs		0	0	0	3	0
71	silk trader	He sells silk		0	0	0	3	0
4	Goblin	An ugly green creature		13	16	14	12	3
73	Guide	He gives hints to new adventurers		0	0	0	7	0
75	Peksa	A helmet salesman		11	11	8	7	0
89	Highwayman	He holds up passers by		13	14	15	13	5
77	Fred the farmer	An old farmer		11	11	8	7	0
104	Moss Giant	his beard seems to have a life of its own		65	62	61	60	5
135	Ice Giant	He's got icicles in his beard		70	67	66	70	5
136	King Scorpion	Wow scorpions shouldn't grow that big		39	40	38	30	5
137	Pirate	A vicious pirate		30	35	25	20	5
82	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
83	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
84	Zeke	He sells Scimitars		0	0	0	3	0
85	Louie Legs	He might want to sell something		0	0	0	3	0
87	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
88	Shop Assistant	Maybe she'd like to buy some of my junk		0	0	0	3	0
90	Kebab Seller	A seller of strange food		0	0	0	3	0
91	Chicken	Yep definitely a chicken		4	3	4	3	0
92	Ernest	A former chicken		3	3	3	3	0
94	Dwarf	A short angry guy		20	20	17	16	1
95	Banker	He can look after my money		11	11	8	7	0
97	Morgan	A frigtened villager		11	11	8	7	0
98	Dr Harlow	His nose is very red		11	11	8	7	0
101	Cassie	She sells shields		30	35	25	10	0
103	Ranael	A shopkeeper of some sort		30	35	25	10	0
105	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
106	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
107	Witch	She's got warts		30	35	25	10	0
110	Sir Amik Varze	The leader of the white knights		58	55	60	52	0
112	Valaine	She runs the champion's store		30	35	25	10	0
113	Drogo	He runs a mining store		20	20	17	16	0
115	Flynn	The mace salesman		6	15	22	22	0
116	Wyson the gardener	An old gardener		8	10	8	7	0
117	Wizard Mizgog	An old wizard		10	20	15	3	0
118	Prince Ali	A young prince		20	20	20	20	0
119	Hassan	the Chancellor to the emir		20	20	20	20	0
120	Osman	He looks a little shifty		20	20	20	20	0
121	Joe	Lady Keli's head guard		40	40	40	40	0
122	Leela	She comes from Al Kharid		20	20	20	20	0
124	Ned	An old sailor		20	20	20	20	0
125	Aggie	A witch		30	35	25	10	0
126	Prince Ali	That is an effective disguise		10	10	10	10	0
128	Redbeard Frank	A pirate		30	35	25	10	0
129	Wydin	A grocer		0	0	0	3	0
130	shop assistant	I can buy swords off him		0	0	0	3	0
131	Brian	An axe seller		0	0	0	3	0
132	squire	A young squire		0	0	0	3	0
134	Thurgo	A short angry guy		20	20	17	16	0
138	Sir Vyvin	One of the white knights of Falador		58	55	60	52	0
139	Monk of Zamorak	An evil cleric		28	28	32	30	1
140	Monk of Zamorak	An evil cleric		18	18	22	20	1
141	Wayne	An armourer		6	15	22	22	0
142	Barmaid	a pretty barmaid		30	35	25	10	0
79	Witch	She's got warts		30	35	25	10	3
144	Doric	A dwarven smith		20	20	17	16	0
146	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
147	Guide	She gives hints to new adventurers		0	0	0	7	0
148	Hetty	A witch		30	35	25	10	0
149	Betty	A witch		30	35	25	10	0
150	Bartender	I could get a beer off him		0	0	0	3	0
151	General wartface	An ugly green creature		13	16	14	12	0
152	General Bentnoze	An ugly green creature		13	16	14	12	0
158	Ice warrior	A strange inhuman warrior		59	57	56	59	5
177	Rat	Overgrown vermin		15	16	12	10	5
155	Herquin	A gem merchant		0	0	0	3	0
156	Rommik	The owner of the crafting shop		0	0	0	3	0
178	Ghost	Ooh spooky		23	23	30	25	5
160	Thrander	A smith of some sort		6	15	22	22	0
161	Border Guard	a guard from Al Kharid		18	20	17	19	0
162	Border Guard	a guard from Al Kharid		18	20	17	19	0
163	Customs Officer	She is here to stop smugglers		14	23	12	15	0
164	Luthas	The owner of the banana plantation		14	23	12	15	0
165	Zambo	He will sell me exotic rum		14	23	12	15	0
166	Captain Tobias	An old sailor		20	20	20	20	0
168	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
169	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
170	Seaman Lorris	A young sailor		20	20	20	20	0
171	Seaman Thresnor	A young sailor		20	20	20	20	0
172	Tanner	He makes leather		40	20	60	60	0
173	Dommik	The owner of the crafting shop		0	0	0	3	0
174	Abbot Langley	A Peaceful monk		12	12	13	15	0
175	Thordur	He runs a a tourist attraction		20	20	17	16	0
176	Brother Jered	human		12	12	13	15	0
183	Scavvo	He has lopsided eyes		10	10	10	10	0
185	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
187	Oziach	A strange little man		0	0	0	3	0
191	dwarf	A dwarf who looks after the mining guild		20	20	17	16	0
180	zombie	the living dead		35	32	31	30	5
193	Klarense	A young sailor		20	20	20	20	0
194	Ned	An old sailor		20	20	20	20	0
197	Oracle	A mystic of unknown race		59	57	56	59	0
198	Duke of Lumbridge	Duke Horacio of Lumbridge		15	15	60	30	0
200	Druid	A worshipper of Guthix		28	28	32	30	1
204	Kaqemeex	A wise druid		28	28	32	30	0
205	Sanfew	An old druid		28	28	32	30	0
206	Suit of armour	A dusty old suit of armour		28	30	30	29	0
207	Adventurer	A cleric		12	12	13	15	0
208	Adventurer	A wizard		10	20	15	3	0
209	Adventurer	A Warrior		58	55	60	52	0
210	Adventurer	An archer		30	35	25	10	0
212	Monk of entrana	A Peaceful monk		12	12	13	15	0
213	Monk of entrana	A Peaceful monk		12	12	13	15	0
153	Goblin	An ugly green creature		13	16	14	12	3
181	Lesser Demon	Lesser but still very big		80	78	79	79	5
216	tree spirit	Ooh spooky		105	100	90	85	1
217	cow	It's a dairy cow		9	9	8	8	0
219	Fairy Lunderwin	A fairy merchant		2	2	2	3	0
220	Jakut	An unusual looking merchant		2	2	2	3	0
221	Doorman	He guards the entrance to the faerie market		58	55	60	52	0
222	Fairy Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
223	Fairy Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
224	Fairy banker	He can look after my money		11	11	8	7	0
225	Giles	He runs an ore exchange store		30	30	30	30	0
226	Miles	He runs a bar exchange store		30	30	30	30	0
227	Niles	He runs a fish exchange store		30	30	30	30	0
228	Gaius	he sells very big swords		6	15	22	22	0
230	Jatix	A hard working druid		28	28	32	30	0
231	Master Crafter	The man in charge of the crafter's guild		0	0	0	3	0
233	Noterazzo	A bandit shopkeeper		26	32	33	27	0
234	Bandit	A wilderness outlaw		26	32	33	27	1
235	Fat Tony	A Gourmet Pizza chef		20	20	20	3	0
236	Donny the lad	A bandit leader		36	42	43	37	1
237	Black Heather	A bandit leader		36	42	43	37	1
238	Speedy Keith	A bandit leader		36	42	43	37	1
240	Boy	He doesn't seem very happy		36	42	43	37	0
241	Rat	He seems to live here		2	3	4	2	0
244	shapeshifter	I've not seen anyone like this before		20	28	29	21	1
245	shapeshifter	I think this spider has been genetically modified		30	38	39	31	0
246	shapeshifter	Eek! A bear!		40	48	49	41	0
247	shapeshifter	A sinister looking wolf		50	58	59	51	0
250	Harry	I wonder what he's got for sale		0	0	0	3	0
232	Bandit	He's ready for a fight		26	32	33	27	5
253	Achetties	One of Asgarnia's greatest heros		48	45	50	42	0
255	Grubor	A rough looking thief		18	15	16	12	0
256	Trobert	A well dressed thief		13	14	15	13	0
257	Garv	A diligent guard		31	31	30	22	0
258	guard	A vicious pirate		30	35	25	20	0
239	White wolf sentry	A vicious mountain wolf		31	30	32	34	5
261	Charlie the cook	Head cook of the Shrimp and parrot		20	20	20	3	0
264	Pirate	A vicious pirate		33	38	28	23	1
267	Seth	He runs a fish exchange store		30	30	30	30	0
268	Banker	He can look after my money		11	11	8	7	0
269	Helemos	A retired hero		48	45	50	42	0
272	Velrak the explorer	he looks cold and hungry		3	3	3	3	0
273	Sir Lancelot	A knight of the round table		58	55	60	52	0
274	Sir Gawain	A knight of the round table		58	55	60	52	0
275	King Arthur	A wise old king		58	55	60	52	0
243	Grey wolf	A sinister looking wolf		65	60	62	69	5
278	Davon	An amulet trader		30	35	25	20	0
279	Bartender	I could get some grog off him		0	0	0	3	0
280	Arhein	A merchant		0	0	0	3	0
281	Morgan le faye	An evil sorceress		30	35	25	10	0
282	Candlemaker	He makes and sells candles		6	15	22	22	0
284	lady	She has a hint of magic about her		0	0	0	3	0
252	Firebird	Probably not a chicken		7	6	7	5	3
286	Beggar	A scruffy looking chap		7	9	8	5	0
287	Merlin	An old wizard		10	20	15	3	0
288	Thrantax	A freshly summoned demon		90	90	90	90	0
289	Hickton	The owner of the archery store		0	0	0	3	0
290	Black Demon	A big scary jet black demon		158	155	157	157	5
297	Frincos	A Peaceful monk		12	12	13	15	0
298	Otherworldly being	Is he invisible or just a set of floating clothes?		66	66	66	66	1
299	Owen	He runs a fish exchange store		30	30	30	30	0
300	Thormac the sorceror	A powerful sorcerrer		27	27	24	24	0
301	Seer	An old wizard		18	18	15	14	0
302	Kharid Scorpion	a smaller less dangerous scorpion		22	21	24	17	0
303	Kharid Scorpion	a smaller less dangerous scorpion		22	21	24	17	0
304	Kharid Scorpion	a smaller less dangerous scorpion		22	21	24	17	0
305	Barbarian guard	Not very civilised		18	18	15	14	0
307	man	A well dressed nobleman		11	11	8	7	0
308	gem trader	He sells gems		0	0	0	3	0
309	Dimintheis	A well dressed nobleman		11	11	8	7	0
310	chef	A busy looking chef		20	20	20	3	0
313	Boot the Dwarf	A short angry guy		20	20	17	16	0
314	Wizard	A young wizard		18	18	15	14	0
316	Captain Barnaby	An old sailor		20	20	20	20	0
291	Black Dragon	A fierce dragon with black scales!		210	210	190	190	5
293	Monk of Zamorak	An evil cleric		48	48	52	40	5
294	Hellhound	Hello nice doggy		114	115	112	116	5
295	Animated axe	a magic axe with a mind of it's own		45	50	60	30	5
311	Hobgoblin	An ugly green creature		48	49	47	49	5
312	Ogre	A large dim looking humanoid		70	72	33	60	5
315	Chronozon	Chronozon the blood demon		182	183	60	60	5
325	Baker	He sells hot baked bread		20	20	20	3	0
326	silk merchant	He buys silk		0	0	0	3	0
328	silver merchant	He deals in silver		0	0	0	3	0
329	spice merchant	He sells exotic spices		20	20	20	3	0
330	gem merchant	He sells gems		0	0	0	3	0
331	Zenesha	A shopkeeper of some sort		30	35	25	10	0
332	Kangai Mau	A tribesman		0	0	0	3	0
333	Wizard Cromperty	An old wizard		10	20	15	3	0
334	RPDT employee	A delivery man		12	12	12	13	0
335	Horacio	An old gardener		8	10	8	7	0
336	Aemad	He helps run the adventurers store		6	15	22	22	0
337	Kortan	He helps run the adventurers store		6	15	22	22	0
339	Make over mage	He can change how I look		0	0	0	3	0
340	Bartender	I could get a beer off him		0	0	0	3	0
341	chuck	A wood merchant		0	0	0	3	0
343	Shadow spider	Is it a spider or is it a shadow		52	54	51	55	5
345	Grandpa Jack	A wistful old man		20	20	20	20	0
346	Sinister stranger	not your average fisherman		35	40	65	35	0
347	Bonzo	Fishing competition organiser		30	30	30	30	0
344	Fire Giant	A big guy with red glowing skin		105	110	112	111	5
349	Morris	Fishing competition organiser		30	30	30	30	0
353	Big Dave	A well built fisherman		18	15	16	12	0
296	Black Unicorn	It's a sort of unicorn		33	31	33	29	3
355	Mountain Dwarf	A short angry guy		20	20	17	16	0
357	Brother Cedric	A Peaceful monk		12	12	13	15	0
359	zombie	The living dead		23	23	28	24	5
360	Lucien	He walks with a slight limp		23	24	22	17	0
362	guardian of Armadyl	A worshipper of Armadyl		58	58	52	50	0
363	guardian of Armadyl	A worshipper of Armadyl		58	58	52	50	0
361	The Fire warrior of lesarkus	A strange red humanoid		72	72	50	59	5
365	winelda	A witch		30	35	25	10	0
366	Brother Kojo	A Peaceful monk		12	12	13	15	0
368	Master fisher	The man in charge of the fishing guild		18	15	16	12	0
369	Orven	He runs a fish exchange store		30	30	30	30	0
370	Padik	He runs a fish exchange store		30	30	30	30	0
371	Shopkeeper	He smells of fish		0	0	0	3	0
372	Lady servil	She look's wealthy		1	1	1	5	0
373	Guard	It's one of General Khazard's guard's		31	31	30	22	0
374	Guard	It's one of General Khazard's guard's		31	31	30	22	0
376	Guard	It's one of General Khazard's guard's		31	31	30	22	0
377	Jeremy Servil	A young squire		0	0	0	3	0
378	Justin Servil	Jeremy servil's father		0	0	0	3	0
379	fightslave joe	He look's mistreated and weak		0	0	0	3	0
380	fightslave kelvin	He look's mistreated and weak		0	0	0	3	0
381	local	A scruffy looking chap		7	9	8	5	0
382	Khazard Bartender	A tough looking barman		0	0	0	3	0
384	Khazard Ogre	Khazard's strongest ogre warrior		70	72	33	60	5
385	Guard	It's one of General Khazard's guard's		31	31	30	22	0
387	hengrad	He look's mistreated and weak		0	0	0	3	0
389	Stankers	A cheerful looking fellow		0	0	0	3	0
390	Docky	An old sailor		20	20	20	20	0
391	Shopkeeper	Maybe he'd like to buy some of my junk		0	0	0	3	0
392	Fairy queen	A very little queen		2	2	2	3	0
393	Merlin	An old wizard		10	20	15	3	0
394	Crone	A strange old lady		30	35	25	10	0
395	High priest of entrana	A Peaceful monk		12	12	13	15	0
396	elkoy	It's a tree gnome		3	3	3	3	0
397	remsai	It's a tree gnome		3	3	3	3	0
398	bolkoy	It's a tree gnome		3	3	3	3	0
400	bolren	It's a gnome he look's important		3	3	3	3	0
388	Bouncer	Hello nice doggy		130	130	112	116	5
421	Tribesman	A primative warrior		40	38	39	39	5
403	brother Galahad	A Peaceful monk		12	12	13	15	0
404	tracker 1	It's a tree gnome		3	3	3	3	0
405	tracker 2	It's a tree gnome		3	3	3	3	0
406	tracker 3	It's a tree gnome		3	3	3	3	0
408	commander montai	It's a tree gnome		3	3	3	3	0
411	Sir Percival	He's covered in pieces of straw		58	55	60	52	0
412	Fisher king	an old king		15	15	60	30	0
413	maiden	She has a far away look in her eyes		2	2	4	3	0
414	Fisherman	an old fisherman		15	15	60	30	0
415	King Percival	The new fisher king		58	55	60	52	0
418	ceril	It's Sir ceril carnillean a local noblemen		11	11	8	7	0
422	henryeta	It's a wealthy looking woman		2	2	4	3	0
423	philipe	It's a young well dressed boy		0	0	0	3	0
358	Necromancer	A crazy evil necromancer		28	28	42	40	3
425	cult member	An suspicous looking man in black 		20	20	20	20	1
428	Khazard commander	It's one of General Khazard's commander's		45	50	50	22	5
427	alomone	A musculer looking man in black 		56	48	46	20	0
429	claus	the carnillean family cook		20	20	20	3	0
430	1st plague sheep	The sheep has the plague		0	0	0	3	0
431	2nd plague sheep	The sheep has the plague		0	0	0	3	0
432	3rd plague sheep	The sheep has the plague		0	0	0	3	0
434	Farmer brumty	He looks after livestock in this area		18	15	16	12	0
435	Doctor orbon	A local doctor		20	20	20	3	0
436	Councillor Halgrive	A town counceller		20	20	20	20	0
437	Edmond	A local civilian		20	20	20	20	0
477	King Black Dragon	The biggest meanest dragon around		250	250	240	240	5
443	Jethick	A cynical old man		20	18	12	10	0
444	Mourner	A mourner or plague healer		2	2	2	3	0
446	Ted Rehnison	The head of the Rehnison family		11	11	8	7	0
447	Martha Rehnison	A fairly poor looking woman		14	11	10	13	0
448	Billy Rehnison	The Rehnisons eldest son		40	20	60	60	0
449	Milli Rehnison	She doesn't seem very happy		36	42	43	37	0
450	Alrena	She look's concerned		1	1	1	5	0
451	Mourner	A mourner or plague healer		2	2	2	3	0
452	Clerk	A bueracratic administrator		2	2	4	3	0
453	Carla	She look's upset		1	1	1	5	0
455	Caroline	A well dressed middle aged lady		1	1	1	5	0
456	Holgart	An old sailor		20	20	20	20	0
457	Holgart	An old sailor		20	20	20	20	0
458	Holgart	An old sailor		20	20	20	20	0
459	kent	caroline's husband		40	20	60	60	0
460	bailey	the fishing platform cook		20	20	20	3	0
461	kennith	A young scared looking boy		0	0	0	3	0
465	Elena	She doesn't look too happy		1	1	1	5	0
467	Watto	He doesn't seem to mind his lack of legs		30	30	30	30	0
468	Recruiter	A member of the Ardougne royal army		30	40	70	65	0
469	Head mourner	In charge of people with silly outfits		2	2	2	3	0
470	Almera	A woman of the wilderness		1	1	1	5	0
471	hudon	A young boisterous looking lad		0	0	0	3	0
472	hadley	A happy looking fellow		15	15	60	30	0
474	Combat instructor	He will tell me how to fight		30	40	70	65	0
475	golrie	It's a tree gnome		3	3	3	3	0
478	cooking instructor	Talk to him to learn about runescape food		20	20	20	3	0
479	fishing instructor	He smells of fish		18	15	16	12	0
480	financial advisor	He knows about money		0	0	0	3	0
481	gerald	An old fisherman		18	15	16	12	0
482	mining instructor	A short angry guy		20	20	17	16	0
483	Elena	She looks concerned		1	1	1	5	0
484	Omart	A nervous looking fellow		15	15	60	30	0
485	Bank assistant	She can look after my stuff		11	11	8	7	0
486	Jerico	He looks friendly enough		18	15	16	12	0
487	Kilron	He looks shifty		18	15	16	12	0
489	Quest advisor	I wonder what advise he has to impart		30	40	70	65	0
490	chemist	human		3	3	3	7	0
491	Mourner	A mourner or plague healer		2	2	2	3	0
492	Mourner	A mourner or plague healer		2	2	2	3	0
426	Lord hazeel	He could do with some sun		78	75	80	170	3
494	Magic Instructor	An old wizard		10	20	15	3	0
496	Community instructor	This is the last advisor - honest		2	2	4	3	0
497	boatman	An old sailor		20	20	20	20	0
499	controls guide	He's ready for a fight		26	32	33	27	0
500	nurse sarah	She's quite a looker		1	1	1	5	0
501	Tailor	He's ready for a party		26	32	33	27	0
498	skeleton mage	It rattles as it walks		23	24	20	17	5
503	Guard	He tries to keep order around here		31	31	30	22	0
504	Chemist	He looks clever enough		26	32	33	27	0
505	Chancy	He's ready for a bet		26	32	33	27	0
506	Hops	He's drunk		26	32	33	27	0
508	Guidor	He's not that ill		26	32	33	27	0
509	Chancy	He's ready for a bet		26	32	33	27	0
510	Hops	He's drunk		26	32	33	27	0
511	DeVinci	He has a colourful personality		26	32	33	27	0
512	king Lathas	King Lanthas of east ardounge		15	15	60	30	0
513	Head wizard	He runs the wizards guild		10	20	15	3	0
514	Magic store owner	An old wizard		10	20	15	3	0
521	Jungle Spider	A venomous deadly spider		47	45	46	50	5
517	Trufitus	A wise old witch doctor		5	10	5	7	0
523	Jogre	An aggressive humanoid		70	72	33	60	5
531	Ogre chieftan	A slightly bigger uglier ogre		90	92	53	80	5
520	Bartender	I could get a beer off him		0	0	0	3	0
522	Jiminua	She looks very interested in selling some of her wares.		0	0	0	3	0
524	Guard	He tries to keep order around here		31	31	30	22	0
542	UndeadOne	One of Rashaliyas Minions		50	80	59	59	5
527	Guard	He tries to keep order around here		31	31	30	22	0
528	shop keeper	he sells weapons		0	0	0	3	0
529	Bartender	I could get a beer off him		0	0	0	3	0
530	Frenita	runs a cookery shop		0	0	0	3	0
532	rometti	It's a well dressed tree gnome		3	3	3	3	0
533	Rashiliyia	A willowy ethereal being who floats above the ground		80	80	80	80	0
534	Blurberry	It's a red faced tree gnome		3	3	3	3	0
535	Heckel funch	It's another jolly tree gnome		3	3	3	3	0
536	Aluft Gianne	It's a tree gnome chef		3	3	3	3	0
538	Irena	human		0	0	0	0	0
539	Mosol	A jungle warrior		0	0	0	3	0
540	Gnome banker	It's tree gnome banker		3	3	3	3	0
541	King Narnode Shareen	It's a gnome he look's important		3	3	3	3	0
543	Drucas	engraver		20	20	20	20	0
544	tourist	human		26	32	33	27	0
545	King Narnode Shareen	It's a gnome he look's important		3	3	3	3	0
546	Hazelmere	An ancient looking gnome		3	3	3	3	0
547	Glough	An rough looking gnome		3	3	3	3	0
548	Shar	Concerned about the economy	b38c40	0	0	0	3	0
549	Shantay	human		0	0	0	3	0
551	Gnome guard	A tree gnome guard		31	31	31	31	1
552	Gnome pilot	He can fly the glider		3	3	3	3	0
553	Mehman	local	805030	26	32	33	27	0
554	Ana	This lady doesn't look as if she belongs here.		18	17	15	16	0
556	Gnome pilot	He can fly the glider		3	3	3	3	0
557	Shipyard worker	He look's busy		48	48	42	40	1
558	Shipyard worker	He look's busy		48	48	42	40	1
560	Shipyard foreman	He look's busy		69	60	60	59	0
561	Shipyard foreman	He look's busy		69	60	60	59	0
563	Femi	It's a little tree gnome		3	3	3	3	0
502	Mourner	A mourner or plague healer		25	30	20	25	3
565	Anita	It's a little tree gnome		3	3	3	3	0
566	Glough	An rough looking gnome		3	3	3	3	0
569	Gnome pilot	He can fly the glider		3	3	3	3	0
571	Gnome pilot	He can fly the glider		3	3	3	3	0
572	Gnome pilot	He can fly the glider		3	3	3	3	0
573	Sigbert the Adventurer	A Warrior		58	55	60	52	0
567	Salarin the twisted	A crazy evil druid		68	68	72	70	5
575	Tower guard	He stops people going up the tower		41	41	30	22	0
576	Gnome Trainer	He can advise on training		11	11	11	11	0
577	Gnome Trainer	He can advise on training		11	11	11	11	0
578	Gnome Trainer	He can advise on training		11	11	11	11	0
580	Blurberry barman	He serves cocktails	pickpocket	3	3	3	3	0
581	Gnome waiter	He can serve you gnome food	pickpocket	3	3	3	3	0
568	Black Demon	A big scary jet black demon		178	195	168	160	5
584	Earth warrior	A strange inhuman warrior		54	52	51	54	5
594	Moss Giant	his beard seems to have a life of its own		65	62	61	60	5
595	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
587	Gulluck	He sells weapons		11	10	11	11	0
588	Gunnjorn	Not civilised looking		18	18	15	14	0
589	Zadimus	Ghostly Visage of the dead Zadimus		0	0	0	0	0
590	Brimstail	An ancient looking gnome		3	3	3	3	0
591	Gnome child	He's a little fellow	pickpocket	3	3	3	3	0
597	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
596	Goalie	A gnome ball goal catcher		70	70	70	70	0
601	Referee	He controls the game		3	3	3	3	0
609	Gnome Baller	He's on your team	pass to	70	70	70	70	0
610	Gnome Baller	He's on your team	pass to	70	70	70	70	0
611	Cheerleader	It's a little tree gnome		3	3	3	3	0
612	Cheerleader	It's a little tree gnome		3	3	3	3	0
598	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
599	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
616	Fernahei	An enthusiastic fishing shop owner		5	10	5	7	0
617	Jungle Banker	He can look after my money		11	11	8	7	0
618	Cart Driver	He drives the cart		18	15	16	12	0
619	Cart Driver	He drives the cart		18	15	16	12	0
620	Obli	An intelligent looking shop owner		0	0	0	3	0
621	Kaleb	This is Kaleb Paramaya - a warm and friendly inn owner		0	0	0	3	0
622	Yohnus	This is Yohnus - he runs the local blacksmiths		0	0	0	3	0
624	Yanni	Yanni Salika - He buys and sells antiques.		0	0	0	3	0
625	Official	He helps the referee		3	3	3	3	0
626	Koftik	The kings top tracker		18	18	15	14	0
627	Koftik	The kings top tracker		18	18	15	14	0
628	Koftik	The kings top tracker		18	18	15	14	0
629	Koftik	The kings top tracker		18	18	15	14	0
630	Blessed Vermen	A undead servent of iban		7	15	7	30	1
574	Yanille Watchman	He watches out for invading ogres	pickpocket	41	41	30	22	3
645	Othainian	big red and incredibly evil		78	78	78	78	5
646	Doomion	A big scary jet black demon		98	98	98	98	5
634	slave	He seems possessed		18	17	15	16	1
635	slave	He seems possessed		18	17	15	16	1
636	slave	He seems to have been here a while		18	17	15	16	1
637	slave	He seems possessed		18	17	15	16	1
639	slave	He seems possessed		18	17	15	16	1
640	slave	He seems to have been here a while		18	17	15	16	1
641	Kalrag	I think this is one of Ibans pets		78	88	69	78	1
642	Niloof	A short angry guy		20	20	17	16	0
643	Kardia the Witch	She's got warts		30	35	25	10	0
644	Souless	He seems an empty shell		18	17	15	16	1
648	Klank	A short angry guy		20	20	17	16	0
650	Koftik	The kings top tracker		18	18	15	14	0
651	Goblin guard	An imposing green creature		51	48	51	43	1
652	Observatory Professor	He works in the observatory		3	3	3	7	0
653	Ugthanki	A dangerous type of spitting camel that can temporarily blind an opponent.		45	45	45	45	1
654	Observatory assistant	The Professor's assistant		3	3	3	7	0
657	Kamen	A short angry guy		20	20	17	16	0
658	Iban disciple	An evil follower of Iban		18	18	22	20	1
659	Koftik	The kings top tracker		18	18	15	14	0
661	Chadwell	A sturdy looking gent		18	18	15	14	0
662	Professor	The owner of the observatory		3	3	3	7	0
665	Spirit of Scorpius	The undead spirit of the follower of Zamorak		100	100	100	100	0
666	Scorpion	There are nasty scorpions around this grave		22	21	24	17	0
667	Dark Mage	He works in the ways of dark magic		0	0	0	3	0
668	Mercenary	He seems to be guarding an area		32	48	60	60	1
670	Mercenary	He seems to be guarding an area		32	48	30	48	1
671	Mining Slave	A chained slave forced to mine rocks.		18	17	15	16	1
672	Watchtower wizard	A learned man		10	20	15	3	0
673	Ogre Shaman	An intelligent form of ogre		100	100	100	100	0
674	Skavid	Servant race to the ogres		3	3	3	3	0
675	Ogre guard	These ogres protect the city		90	92	53	80	0
676	Ogre guard	These ogres protect the city		90	92	53	80	0
677	Ogre guard	These ogres protect the city		90	92	53	80	0
678	Skavid	Servant race to the ogres		3	3	3	3	0
679	Skavid	Servant race to the ogres		3	3	3	3	0
680	Og	The chieftan of this ogre tribe		90	92	53	80	0
681	Grew	The chieftan of this ogre tribe		90	92	53	80	0
683	Gorad	A high ranking ogre official		90	92	53	80	1
684	Ogre guard	this creature looks very tough		90	98	99	99	1
685	Yanille Watchman	A captured guard of Yanille		41	41	30	22	0
686	Ogre merchant	He sells ogre-inspired items		70	72	33	60	0
687	Ogre trader	He trades in metals		70	72	33	60	0
688	Ogre trader	He trades in food		70	72	33	60	0
689	Ogre trader	He trades in food		70	72	33	60	0
690	Mercenary	He seems to be guarding an area		32	48	30	48	1
691	City Guard	high ranking ogre guards		90	92	53	80	0
693	Lawgof	He guards the mines		20	20	17	16	0
694	Dwarf	A short angry guy		20	20	17	16	1
695	lollk	He looks scared		20	20	17	16	0
696	Skavid	Servant race to the ogres		3	3	3	3	0
647	Holthion	big red and incredibly evil		78	78	78	78	5
698	Nulodion	He's the head of black guard weapon development		20	20	17	16	0
632	Paladin	A paladin of Ardougne		88	85	55	57	3
700	Al Shabim	The leader of a nomadic Bedabin desert people - sometimes referred to as the 'Tenti's'		0	0	0	3	0
701	Bedabin Nomad	A Bedabin nomad - they live in the harshest extremes in the desert		0	0	0	3	0
702	Captain Siad	He's in control of the whole mining camp.		48	48	48	48	1
704	Ogre citizen	A denizen of Gu'Tanoth		70	72	33	60	1
713	kolodion	He runs the mage arena		10	20	15	3	5
706	Ogre	A large dim looking humanoid		70	72	33	60	1
707	Skavid	Servant race to the ogres		3	3	3	3	0
708	Skavid	Servant race to the ogres		3	3	3	3	0
709	Skavid	Servant race to the ogres		3	3	3	3	0
710	Draft Mercenary Guard	He's quickly drafted in to deal with trouble makers		32	48	60	60	1
711	Mining Cart Driver	He drives the mining cart		18	15	16	12	0
712	kolodion	He runs the mage arena		10	20	15	3	0
714	Gertrude	A busy housewife		20	20	20	20	0
715	Shilop	A young boisterous looking lad		0	0	0	3	0
717	Shantay Pass Guard	He seems to be guarding the Shantay Pass		32	32	32	32	1
720	Assistant	He is an assistant to Shantay and helps him to run the pass.		0	0	0	3	0
722	Workman	This person is working on the site	pickpocket	11	11	8	7	0
723	Examiner	As you examine the examiner you examine that she is indeed an examiner!!		1	1	1	5	0
724	Student	A student busily digging!		0	0	0	3	0
725	Student	A student busily digging!		20	20	20	20	0
726	Guide	This person specialises in panning for gold		10	20	15	3	0
727	Student	A student busily digging!		18	20	17	19	0
728	Archaeological expert	An expert on archaeology!		20	20	20	3	0
716	Rowdy Guard	He looks as if he's spoiling for trouble		32	48	60	60	5
730	civillian	She looks aggitated!		0	0	0	3	0
731	civillian	She looks aggitated!		0	0	0	3	0
732	civillian	He looks aggitated!	pickpocket	18	15	16	12	0
734	Murphy	The man in charge of the fishing trawler		18	15	16	12	0
735	Sir Radimus Erkle	A huge muscular man in charge of the Legends Guild		5	10	20	8	0
736	Legends Guild Guard	This guard is protecting the entrance to the Legends Guild.		50	50	50	50	0
737	Escaping Mining Slave	An emancipated slave with cool Desert Clothes.		18	17	15	16	0
738	Workman	This person is working in the mine	pickpocket	11	11	8	7	0
739	Murphy	The man in charge of the fishing trawler		18	15	16	12	0
740	Echned Zekin	An evil spirit of the underworld.		50	50	50	50	0
741	Donovan the Handyman	It's the family odd jobs man		11	11	8	7	0
743	Hobbes the Butler	It's the family butler		11	11	8	7	0
744	Louisa The Cook	It's the family cook		0	0	0	3	0
745	Mary The Maid	The family maid		30	35	25	10	0
746	Stanford The Gardener	It's the family Gardener		8	10	8	7	0
747	Guard	He looks like he's in over his head here		31	31	30	22	0
748	Guard Dog	He doesn't seem pleased to see me		46	45	47	49	0
749	Guard	***EMPTY PLEASE USE OR REPLACE***		8	10	8	7	0
750	Man	A thirsty looking man		11	11	8	7	0
751	Anna Sinclair	The first child of the late Lord Sinclair		11	11	8	7	0
752	Bob Sinclair	The second child of the late Lord Sinclair		11	11	8	7	0
753	Carol Sinclair	The third child of the late Lord Sinclair		11	11	8	7	0
754	David Sinclair	The fourth child of the late Lord Sinclair		11	11	8	7	0
755	Elizabeth Sinclair	The fifth child of the late Lord Sinclair		11	11	8	7	0
756	Frank Sinclair	The sixth child of the late Lord Sinclair		11	11	8	7	0
705	Rock of ages	A huge boulder		150	150	150	150	3
12	Bartender	I could get a beer off him		0	0	0	3	0
33	Apothecary	I wonder if he has any good potions		5	10	5	7	0
54	Aubury	I think he might be a shop keeper		0	0	0	3	0
111	Guildmaster	He's in charge of this place		40	40	40	40	0
123	Lady Keli	An Infamous bandit		20	20	20	20	0
133	Head chef	He looks after the chef's guild		20	20	20	3	0
143	Dwarven shopkeeper	I wonder if he wants to buy any of my junk		20	20	17	16	0
145	Shopkeeper	I wonder what he's got for sale		0	0	0	3	0
157	Grum	Grum the goldsmith		0	0	0	3	0
167	Gerrant	I wonder what he's got for sale		0	0	0	3	0
186	Shop Assistant	Maybe he'd like to buy some of my junk		0	0	0	3	0
199	Dark Warrior	A warrior touched by chaos		23	20	25	17	1
211	Leprechaun	A funny little man who lives in a tree		20	20	17	16	0
215	Monk of entrana	A Peaceful monk		12	12	13	15	0
218	Irksol	Is he invisible or just a set of floating clothes?		2	2	2	3	0
229	Fairy Ladder attendant	A worker in the faerie market		0	0	0	3	0
242	Nora T Hag	She's got warts		30	35	25	10	0
260	Alfonse the waiter	He should get a clean apron		11	11	8	7	0
763	Poison Salesman	Peter Potter - Poison Purveyor		7	9	8	5	0
764	Gujuo	A tall charismatic looking jungle native - he approaches with confidence		60	60	60	60	0
765	Jungle Forester	A woodsman who specialises in large and exotic timber		18	15	16	12	0
766	Ungadulu	An ancient looking Shaman		75	75	75	75	1
769	Nezikchened	An ancient powerful Demon of the Underworld...		178	175	177	160	1
770	Dwarf Cannon engineer	He's the head of black guard weapon development		20	20	17	16	0
771	Dwarf commander	He guards the mines		20	20	17	16	0
772	Viyeldi	The spirit of a dead sorcerer		80	80	80	80	1
773	Nurmof	He sells pickaxes		20	20	17	16	0
774	Fatigue expert	He looks wide awake		8	10	10	13	0
776	Jungle Savage	A savage and fearless Jungle warrior		100	100	60	90	1
22	Lesser Demon	Lesser but still pretty big		80	78	79	79	5
778	Sidney Smith	Sidney Smith - Certification clerk		30	30	30	30	0
779	Siegfried Erkle	An eccentric shop keeper - related to the Grand Vizier of the Legends Guild		30	35	25	10	0
780	Tea seller	He has delicious tea to buy		11	11	8	7	0
781	Wilough	A young son of gertrudes		0	0	0	3	0
782	Philop	Gertrudes youngest		0	0	0	3	0
783	Kanel	Gertrudes youngest's twin brother		0	0	0	3	0
785	Sir Radimus Erkle	A huge muscular man in charge of the Legends Guild		5	10	20	8	0
788	Fionella	She runs the legend's general store		30	35	25	10	0
792	Gundai	He must get lonely out here		18	15	16	12	0
793	Lundail	He sells rune stones		18	15	16	12	0
0	Unicorn	It's a unicorn		23	21	23	19	3
45	skeleton	It rattles as it walks		35	32	30	29	5
74	Giant Spider	I think this spider has been genetically modified		34	30	31	32	5
99	Deadly Red spider	I think this spider has been genetically modified		35	40	36	35	5
179	skeleton	it rattles when it walks		35	32	30	29	5
251	Thug	He likes hitting things		17	19	20	18	5
283	lady	She has a hint of magic about her		0	0	0	3	0
285	lady	She has a hint of magic about her		0	0	0	3	0
306	Bartender	I could get a beer off him		0	0	0	3	0
317	Customs Official	She's here to stop smugglers		14	23	12	15	0
327	Fur trader	A buyer and seller of animal furs		0	0	0	3	0
338	zoo keeper	He looks after Ardougne city zoo		20	20	20	20	1
350	Brother Omad	A Peaceful monk		12	12	13	15	0
354	Joshua	A grumpy fisherman		18	15	16	12	0
375	Guard	It's one of General Khazard's guard's		31	31	30	22	0
419	butler	It's the carnillean family butler		11	11	8	7	0
424	clivet	A strange looking man in black 		20	20	20	20	0
433	4th plague sheep	The sheep has the plague		0	0	0	3	0
445	Mourner	A mourner or plague healer		2	2	2	3	0
454	Bravek	The city warder of West Ardougne		15	15	60	30	0
466	jinno	He doesn't seem to mind his lack of legs		30	30	30	30	0
476	Guide	She gives hints to new adventurers		0	0	0	7	0
488	Guidor's wife	She looks rather concerned		1	1	1	5	0
493	Wilderness guide	He's ready for a fight		26	32	33	27	0
507	DeVinci	He has a colourful personality		26	32	33	27	0
515	Wizard Frumscone	A confused looking wizard		10	20	15	3	0
526	Guard	He tries to keep order around here		31	31	30	22	0
537	Hudo glenfad	It's another jolly tree gnome		3	3	3	3	0
550	charlie	Poor guy?		0	0	0	3	1
559	Shipyard worker	He look's busy		48	48	42	40	1
564	Femi	It's a little tree gnome		3	3	3	3	0
570	Gnome pilot	He can fly the glider		3	3	3	3	0
579	Gnome Trainer	He can advise on training	pickpocket	11	11	11	11	0
623	Serevel	This is Serevel - he sells tickets for the 'Lady of the Waves'		0	0	0	3	0
631	Blessed Spider	One of iban's eight legged friends		34	45	31	32	1
638	slave	He seems to have been here a while		18	17	15	16	1
649	Iban	You feel terror just looking at him		23	24	22	17	0
669	Mercenary Captain	He's in control of the local guards.	watch	48	48	80	80	1
682	Toban	The chieftan of this ogre tribe		90	92	53	80	0
692	Mercenary	He seems to be guarding this area		32	48	30	48	1
699	Dwarf	A short angry guy		20	20	17	16	1
703	Bedabin Nomad Guard	A Bedabin nomad guard - he's protecting something important		70	70	70	70	1
719	Shantay Pass Guard	He seems to be guarding the Shantay Pass		32	32	32	32	0
733	Murphy	The man in charge of the fishing trawler		18	15	16	12	0
742	Pierre the Dog Handler	It's the guy who looks after the family guard dog		11	11	8	7	0
767	Ungadulu	An ancient looking Shaman - he looks very strange with glowing red eyes...		75	75	75	75	1
784	chamber guardian	He hasn't seen much sun latley		18	15	16	12	0
5	Hans	A castle servant		3	3	3	3	3
6	cow	It's a multi purpose cow		9	9	8	8	3
11	Man	One of runescapes many citizens	pickpocket	11	11	8	7	3
29	Rat	A small muddy rat		2	3	4	2	3
34	spider	Incey wincey		1	5	2	2	3
35	Delrith	A freshly summoned demon		37	42	35	7	3
37	Weaponsmaster	The phoenix gang quartermaster		28	35	20	20	3
50	skeleton	It rattles when it walks		21	20	18	18	3
292	Poison Spider	I think this spider has been genetically modified		68	60	62	64	5
367	Dungeon Rat	Overgrown vermin		22	20	10	12	5
386	Khazard Scorpion	A large angry scorpion		40	45	48	30	5
407	Khazard troop	It's one of General Khazard's warrior's		31	31	30	22	5
603	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
62	Goblin	An ugly green creature		9	8	9	5	3
63	farmer	He grows the crops in this area	pickpocket	18	15	16	12	3
64	Thief	He'll take anything that isn't nailed down		23	24	22	17	3
66	Black Knight	A sinister looking knight		48	45	50	42	3
76	Barbarian	Not civilised looking		18	18	15	14	3
78	Gunthor the Brave	The barbarians fearless leader		38	37	40	35	3
80	Ghost	Ooh spooky		23	23	30	25	3
81	Wizard	An old wizard		18	18	15	14	3
93	Monk	A Peaceful monk		12	12	13	15	3
96	Count Draynor	A vicious vampire		35	40	65	35	3
100	Guard	He's here to guard this fortress	pickpocket	31	31	30	22	3
102	White Knight	A chivalrous knight		58	55	60	52	3
108	Black Knight	A sinister looking knight		48	45	50	42	3
109	Greldo	A small green warty creature		9	8	9	5	3
114	Imp	A cheeky little imp		5	4	4	8	3
127	Jailguard	I wonder what he's guarding		36	34	34	32	3
154	Goblin	An ugly green creature		13	16	14	12	3
159	Warrior	A skilled fighter	pickpocket	30	35	25	20	3
192	Wormbrain	Dumb even by goblin standards		9	8	9	5	3
259	Grip	Scar face petes head guard		31	31	60	62	3
276	Sir Mordred	An evil knight		60	57	62	54	3
318	Man	One of Runescape's citizens	pickpocket	11	11	8	7	3
319	farmer	An humble peasant	pickpocket	18	15	16	12	3
320	Warrior	A skilled fighter	pickpocket	30	35	25	20	3
321	Guard	He tries to keep the law and order around here	pickpocket	31	31	30	22	3
322	Knight	A knight of Ardougne	pickpocket	58	55	60	52	3
323	Paladin	A paladin of Ardougne	pickpocket	88	85	55	57	3
324	Hero	A Hero of Ardougne	pickpocket	88	85	80	82	3
342	Rogue	He needs a shave	pickpocket	23	24	22	17	3
348	Forester	He looks after McGrubor's wood		23	24	22	17	3
364	Lucien	He walks with a limp		23	24	22	17	3
383	General Khazard	He look's real nasty		78	75	80	170	3
401	Black Knight titan	He is blocking the way		148	145	150	142	3
402	kalron	he look's lost		3	3	3	3	3
409	gnome troop	It's a tree gnome trooper		3	3	3	3	3
410	khazard warlord	He look's real nasty		78	75	80	170	3
416	unhappy peasant	He looks tired and hungry		28	25	26	22	3
417	happy peasant	He looks well fed and full of energy		28	25	26	22	3
420	carnillean guard	It's a carnillean family guard		31	31	30	22	3
438	Citizen	He look's tired		10	12	11	13	3
439	Citizen	He look's frightened		8	10	10	13	3
440	Citizen	She look's frustrated		14	11	10	13	3
441	Citizen	He look's angry		18	20	20	23	3
442	Citizen	He look's disillusioned		20	18	12	10	3
462	Platform Fisherman	an emotionless fisherman		15	15	60	30	3
463	Platform Fisherman	an emotionless fisherman		15	15	60	30	3
464	Platform Fisherman	an emotionless fisherman		15	15	60	30	3
473	Rat	Overgrown vermin		8	15	2	3	3
516	target practice zombie	The living dead		23	23	28	24	3
518	Colonel Radick	A soldier of the town of Yanille		30	40	70	65	3
519	Soldier	A soldier of the town of Yanille		31	31	30	22	3
525	Ogre	Useful for ranged training		70	72	33	60	3
582	Gnome guard	A tree gnome guard	pickpocket	17	31	31	31	3
583	Gnome child	that's a little gnome	pickpocket	3	3	3	3	3
585	Gnome child	He's a little fellow	pickpocket	3	3	3	3	3
586	Gnome child	hello little gnome	pickpocket	3	3	3	3	3
593	Gnome local	A tree gnome villager	pickpocket	3	3	3	3	3
613	Nazastarool Zombie	One of Rashaliyas Minions		90	95	70	80	3
615	Nazastarool Ghost	One of Rashaliyas Minions		90	95	70	80	3
633	Paladin	A paladin of Ardougne		88	85	55	57	3
697	Ogre guard	These ogres protect the city		90	92	53	80	3
729	civillian	He looks aggitated!		18	20	17	19	3
3	Chicken	Yep definitely a chicken		4	3	4	3	3
65	Guard	He tries to keep order around here	pickpocket	31	31	30	22	3
72	Man	One of Runescapes many citizens	pickpocket	11	11	8	7	3
86	Warrior	A member of Al Kharid's military	pickpocket	18	20	17	19	3
777	Oomlie Bird	A variety of flightless jungle fowl - it has a sharp beak and a bad temper.		20	50	20	40	3
356	Mountain Dwarf	A short angry guy		30	30	27	26	3
399	local gnome	It's a young tree gnome		3	3	3	3	3
495	Mourner	A mourner or plague healer		30	20	20	19	3
592	Gnome local	A tree gnome villager	pickpocket	9	9	9	9	3
614	Nazastarool Skeleton	One of Rashaliyas Minions		90	95	70	80	3
182	Melzar the mad	He looks totally insane		47	47	44	44	5
184	Greater Demon	big red and incredibly evil		88	86	87	87	5
188	Bear	Eek! A bear!		28	27	25	27	5
189	Black Knight	An armoured follower of Zamorak		48	45	50	42	5
190	chaos Dwarf	a dwarf gone bad		60	58	59	61	5
195	skeleton	A Taller than normal skeleton		55	52	50	59	5
196	Dragon	A powerful and ancient dragon		110	110	150	110	5
201	Red Dragon	A big powerful dragon		140	140	140	140	5
202	Blue Dragon	A mother dragon		105	105	105	105	5
203	Baby Blue Dragon	Young but still dangerous		50	50	50	50	5
214	zombie	The living dead		35	32	31	30	5
248	White wolf	A vicious mountain wolf		41	40	42	44	5
249	Pack leader	A vicious mountain wolf		71	70	72	74	5
254	Ice queen	The leader of the ice warriors		104	105	101	104	5
262	Guard Dog	He doesn't seem pleased to see me		46	45	47	49	5
263	Ice spider	I think this spider has been genetically modified		65	60	66	65	5
265	Jailer	Guards prisoners for the black knights		53	50	55	47	5
266	Lord Darquarius	A black knight commander		78	75	80	72	5
270	Chaos Druid	A crazy evil druid		18	18	22	20	5
277	Renegade knight	He isn't very friendly		53	50	55	48	5
351	Thief	A dastardly blanket thief		23	24	22	17	5
352	Head Thief	A dastardly blanket thief		33	34	32	37	5
555	Chaos Druid warrior	A crazy evil druid		48	48	42	40	5
562	Gnome guard	A tree gnome guard		23	23	23	23	5
600	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
602	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
604	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
605	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
606	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
607	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
608	Gnome Baller	A tree gnome ball player	tackle	70	70	70	70	5
655	Souless	A servent to zamorak		23	23	28	24	5
656	Dungeon spider	A nasty poisonous arachnid		10	25	20	35	5
663	San Tojalon	The animated spirit of San Tojalon		120	120	120	120	5
664	Ghost	A doomed victim of zamorak		30	33	33	20	5
718	Rowdy Slave	A slave who's looking for trouble.		18	17	15	16	5
721	Desert Wolf	A vicious Desert wolf		31	30	32	34	5
757	kolodion	He's a shape shifter		70	72	55	65	5
758	kolodion	He's a shape shifter		78	47	69	78	5
760	kolodion	He's a shape shifter		98	105	85	107	5
761	Irvig Senay	The animated spirit of Irvig Senay		125	125	125	125	5
271	Poison Scorpion	It has a very vicious looking tail		27	26	29	23	5
768	Death Wing	A supernatural creature of the underworld		80	80	80	80	5
775	Karamja Wolf	A hungry		61	61	61	61	5
786	Pit Scorpion	Very vicious little scorpions		40	35	45	35	5
787	Shadow Warrior	A sinsistar shadowy figure		61	61	68	67	5
789	Battle mage	He kills in the name of guthix		0	0	90	120	5
790	Battle mage	He kills in the name of zamarok		0	0	90	120	5
791	Battle mage	He kills in the name of Saradomin		0	0	90	120	5
660	Goblin	These goblins have grown strong		18	24	20	16	5
759	kolodion	He's a shape shifter		23	58	28	78	5
762	Ranalph Devere	The animated spirit of Ranalph Devere		130	130	130	130	5
\.


--
-- Data for Name: prayers; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.prayers (id, name, description, required_level, drain_rate) FROM stdin;
1	Thick skin	Increases your defense by 5%	1	15
2	Burst of strength	Increases your strength by 5%	4	15
3	Clarity of thought	Increases your attack by 5%	7	15
4	Rock skin	Increases your defense by 10%	10	30
5	Superhuman strength	Increases your strength by 10%	13	30
6	Improved reflexes	Increases your attack by 10%	16	30
7	Rapid restore	2x restore rate for all stats except hits	19	5
8	Rapid heal	2x restore rate for hitpoints stat	22	10
9	Protect items	Keep 1 extra item if you die	25	10
10	Steel skin	Increases your defense by 15%	28	60
11	Ultimate strength	Increases your strength by 15%	31	60
12	Incredible reflexes	Increases your attack by 15%	34	60
13	Paralyze monster	Stops monsters from fighting back	37	60
14	Protect from missiles	100% protection from ranged attack	40	60
\.


--
-- Data for Name: shop_items; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.shop_items (storeid, itemid, amount) FROM stdin;
0	77	-1
0	78	-1
0	79	-1
0	80	-1
0	89	-1
0	90	-1
0	91	-1
0	92	-1
0	316	-1
0	426	-1
0	429	-1
0	522	-1
1	110	-1
1	111	-1
1	115	-1
1	116	-1
1	119	-1
1	120	-1
1	122	-1
1	123	-1
1	130	-1
1	131	-1
1	196	-1
1	230	-1
1	248	-1
1	431	-1
1	433	-1
1	1006	-1
2	221	-1
2	483	-1
2	486	-1
2	492	-1
2	495	-1
2	498	-1
3	31	-1
3	32	-1
3	33	-1
3	34	-1
3	41	-1
3	101	-1
3	102	-1
3	103	-1
3	197	-1
3	1213	-1
3	1214	-1
3	1215	-1
3	1216	-1
3	1217	-1
3	1218	-1
4	59	-1
4	188	-1
4	189	-1
4	190	-1
4	638	-1
4	640	-1
4	642	-1
4	644	-1
4	646	-1
4	648	-1
4	650	-1
4	652	-1
4	654	-1
4	656	-1
\.


--
-- Data for Name: shops; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.shops (id, name, general) FROM stdin;
0	Weapon shop	f
1	Armour shop	f
2	Potion shop	f
3	Rune Store	f
4	Range Shop	f
\.


--
-- Data for Name: spell_aggressive_level; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.spell_aggressive_level (id, spell) FROM stdin;
0	1
2	2
4	3
6	4
8	5
11	6
14	7
17	8
20	9
23	10
27	11
32	12
33	25
34	25
35	25
37	13
38	13
39	14
40	14
43	15
44	15
45	16
46	16
\.


--
-- Data for Name: spell_runes; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.spell_runes (spellid, itemid, amount) FROM stdin;
0	33	1
0	35	1
1	32	3
1	34	2
1	36	1
2	32	1
2	33	1
2	35	1
3	32	1
3	46	1
4	33	1
4	34	2
4	35	1
5	32	3
5	34	2
5	36	1
6	31	3
6	33	2
6	35	1
7	32	2
7	34	2
7	40	1
8	33	2
8	41	1
9	32	2
9	34	3
9	36	1
10	31	3
10	40	1
11	32	2
11	33	2
11	41	1
12	31	1
12	33	3
12	42	1
13	33	3
13	46	1
14	33	2
14	34	3
14	41	1
15	33	3
15	34	1
15	42	1
16	33	1
16	42	1
17	31	4
17	33	3
17	41	1
18	32	1
18	33	3
18	42	1
19	33	2
19	34	2
19	41	1
20	33	3
20	38	1
21	31	4
21	40	1
22	33	5
22	42	1
23	32	3
23	33	3
23	38	1
24	31	5
24	46	1
25	31	5
25	38	1
26	32	2
26	42	2
27	33	3
27	34	4
27	38	1
28	31	5
28	40	1
29	32	30
29	46	3
29	611	1
30	34	10
30	46	1
31	34	2
31	42	2
32	31	5
32	33	4
32	38	1
33	31	1
33	33	4
33	619	2
34	31	2
34	33	4
34	619	2
35	31	4
35	33	1
35	619	2
36	34	30
36	46	3
36	611	1
37	33	5
37	619	1
38	31	30
38	46	3
38	611	1
39	32	7
39	33	5
39	619	1
40	33	30
40	46	3
40	611	1
41	32	5
41	34	5
41	825	1
42	32	15
42	34	15
42	46	1
43	33	5
43	34	7
43	619	1
44	32	8
44	34	8
44	825	1
45	31	7
45	33	5
45	619	1
46	32	12
46	34	12
46	825	1
47	31	3
47	33	3
47	619	3
\.


--
-- Data for Name: spells; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.spells (id, name, description, required_level, rune_amount, type, experience) FROM stdin;
0	Wind strike	A strength 1 missile attack	1	2	2	22
1	Confuse	Reduces your opponents attack by 5%	3	3	2	26
2	Water strike	A strength 2 missile attack	5	3	2	30
3	Enchant lvl-1 amulet	For use on sapphire amulets	7	2	3	35
4	Earth strike	A strength 3 missile attack	9	3	2	38
5	Weaken	Reduces your opponents strength by 5%	11	3	2	42
6	Fire strike	A strength 4 missile attack	13	3	2	46
7	Bones to bananas	Changes all held bones into bananas!	15	3	6	50
8	Wind bolt	A strength 5 missile attack	17	2	2	54
9	Curse	Reduces your opponents defense by 5%	19	3	2	58
10	Low level alchemy	Converts an item into gold	21	2	3	62
11	Water bolt	A strength 6 missile attack	23	3	2	66
12	Varrock teleport	Teleports you to Varrock	25	3	0	70
13	Enchant lvl-2 amulet	For use on emerald amulets	27	2	3	74
14	Earth bolt	A strength 7 missile attack	29	3	2	76
15	Lumbridge teleport	Teleports you to Lumbridge	31	3	0	82
16	Telekinetic grab	Take an item you can see but can't reach	33	2	3	86
17	Fire bolt	A strength 8 missile attack	35	3	2	90
18	Falador teleport	Teleports you to Falador	37	3	0	94
19	Crumble undead	Hits skeleton, ghosts & zombies hard!	39	3	2	98
20	Wind blast	A strength 9 missile attack	41	2	2	102
21	Superheat item	Smelt 1 ore without a furnace	43	2	3	106
22	Camelot teleport	Teleports you to Camelot	45	2	0	110
23	Water blast	A strength 10 missile attack	47	3	2	114
24	Enchant lvl-3 amulet	For use on ruby amulets	49	2	3	118
25	Iban blast	A strength 25 missile attack!	50	2	2	120
26	Ardougne teleport	Teleports you to Ardougne	51	2	0	122
27	Earth blast	A strength 11 missile attack	53	3	2	126
28	High level alchemy	Convert an item into more gold	55	2	3	130
29	Charge Water Orb	Needs to be cast on a water obelisk	56	3	5	132
30	Enchant lvl-4 amulet	For use on diamond amulets	57	2	3	134
31	Watchtower teleport	Teleports you to the watchtower	58	2	0	138
32	Fire blast	A strength 12 missile attack	59	3	2	140
33	Claws of Guthix	Summons the power of Guthix	60	3	2	140
34	Saradomin strike	Summons the power of Saradomin	60	3	2	140
35	Flames of Zamorak	Summons the power of Zamorak	60	3	2	140
36	Charge earth Orb	Needs to be cast on an earth obelisk	60	3	5	140
37	Wind wave	A strength 13 missile attack	62	2	2	144
38	Charge Fire Orb	Needs to be cast on a fire obelisk	63	3	5	146
39	Water wave	A strength 14 missile attack	65	3	2	150
40	Charge air Orb	Needs to be cast on an air obelisk	66	3	5	152
41	Vulnerability	Reduces your opponents defense by 10%	66	3	2	152
42	Enchant lvl-5 amulet	For use on dragonstone amulets	68	3	3	156
43	Earth wave	A strength 15 missile attack	70	3	2	160
44	Enfeeble	Reduces your opponents strength by 10%	73	3	2	166
45	Fire wave	A strength 16 missile attack	75	3	2	170
46	Stun	Reduces your opponents attack by 10%	80	3	2	180
47	Charge	Increase your mage arena spells damage	80	3	6	180
\.


--
-- Data for Name: tiles; Type: TABLE DATA; Schema: public; Owner: zach
--

COPY public.tiles (colour, unknown, objecttype) FROM stdin;
-16913	1	0
1	3	1
3	2	0
3	4	0
-16913	2	0
-27685	2	0
25	3	1
12345678	5	1
-26426	1	1
-1	5	1
31	3	1
3	4	0
-4534	2	0
32	2	0
-9225	2	0
-3172	2	0
15	2	0
-2	2	0
-1	3	1
-2	4	0
-2	4	1
-2	0	0
-17793	2	0
-14594	1	1
1	3	0
\.


--
-- Name: boundarys idx_16481_doors_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.boundarys
    ADD CONSTRAINT idx_16481_doors_pkey PRIMARY KEY (id);


--
-- Name: game_objects idx_16487_game_objects_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.game_objects
    ADD CONSTRAINT idx_16487_game_objects_pkey PRIMARY KEY (id);


--
-- Name: npcs idx_16493_npcs_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.npcs
    ADD CONSTRAINT idx_16493_npcs_pkey PRIMARY KEY (id);


--
-- Name: prayers idx_16502_prayers_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.prayers
    ADD CONSTRAINT idx_16502_prayers_pkey PRIMARY KEY (id);


--
-- Name: items idx_16508_items_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT idx_16508_items_pkey PRIMARY KEY (id);


--
-- Name: spells idx_16514_spells_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.spells
    ADD CONSTRAINT idx_16514_spells_pkey PRIMARY KEY (id);


--
-- Name: item_wieldable idx_16526_item_wieldable_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.item_wieldable
    ADD CONSTRAINT idx_16526_item_wieldable_pkey PRIMARY KEY (id);


--
-- Name: spell_aggressive_level idx_16532_spell_aggressive_level_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.spell_aggressive_level
    ADD CONSTRAINT idx_16532_spell_aggressive_level_pkey PRIMARY KEY (id);


--
-- Name: shops idx_16535_shops_pkey; Type: CONSTRAINT; Schema: public; Owner: zach
--

ALTER TABLE ONLY public.shops
    ADD CONSTRAINT idx_16535_shops_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

