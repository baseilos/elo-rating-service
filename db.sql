create table "player"
(
	id bigserial not null
		constraint user_pkey
			primary key,
	first_name varchar(255) not null,
	last_name varchar(255) not null,
	nickname varchar(255),
	active boolean not null,
	registered_at timestamp default now()
);

alter table player owner to postgres;

create table game
(
	id bigserial not null
		constraint game_pkey
			primary key,
	player_white bigint not null,
	player_black bigint not null,
	played_at timestamp default now() not null,
	result integer not null
);

comment on column game.result is '> 0 White winner
= 0 Draw
< Black Winner';

alter table game owner to postgres;

