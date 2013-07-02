SET CLIENT_ENCODING TO UTF8;
SET STANDARD_CONFORMING_STRINGS TO ON;
BEGIN;

DROP TABLE IF EXISTS tube;
create table tube (
  lat float8,
  lng float8,
  geom_4326 geometry(Point,4326),
  name varchar,
  id serial primary key
);

\copy tube(name,lat,lng) FROM 'Java_Code_Test8\tube.csv' DELIMITERS ',' CSV;

UPDATE tube SET geom_4326 = ST_SetSRID(ST_MakePoint(lng, lat),4326);
create index tube_gix on tube using gist ( geom_4326 );


DROP TABLE IF EXISTS t_6043;
create table t_6043 (
  lat float8,
  lng float8,
  geom_4326 geometry(Point,4326),
  name varchar,
  ts timestamp,
  id serial primary key
);

\copy t_6043(name,lat,lng,ts) FROM 'Java_Code_Test8\6043.csv' DELIMITERS ',' CSV;

UPDATE t_6043 SET geom_4326 = ST_SetSRID(ST_MakePoint(lng, lat),4326);
create index t_6043_gix on t_6043 using gist ( geom_4326 );


DROP TABLE IF EXISTS t_5937;
create table t_5937 (
  lat float8,
  lng float8,
  geom_4326 geometry(Point,4326),
  name varchar,
  ts timestamp,
  id serial primary key
);

\copy t_5937(name,lat,lng,ts) FROM 'Java_Code_Test8\5937.csv' DELIMITERS ',' CSV;

UPDATE t_5937 SET geom_4326 = ST_SetSRID(ST_MakePoint(lng, lat),4326);
create index t_5937_gix on t_5937 using gist ( geom_4326 );


select st_distance(a.geom_4326, b.geom_4326), b.name
from tube a, tube b
where a.name = 'New Cross';

COMMIT;

--VACUUM ANALYZE tube;
