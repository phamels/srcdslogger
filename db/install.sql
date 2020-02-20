CREATE TABLE IF NOT EXISTS players
(
    id int(11) NOT NULL AUTO_INCREMENT,
    name varchar(250) NOT NULL,
    steamid varchar(100) NOT NULL,
    points int NULL,
    PRIMARY KEY (id)
);

create table if not exists hitgroups
(
	id int(11) auto_increment,
	hitgroup varchar(255) not null,
    primary key (id)
);

create table if not exists weapons
(
	id int auto_increment,
	weapon varchar(250) not null,
	constraint weapons_pk
		primary key (id)
);

create table if not exists maps
(
	id int auto_increment,
	map varchar(50) not null,
	constraint maps_pk
		primary key (id)
);

create table if not exists rounds
(
	id int auto_increment,
	ct_score int default 0 not null,
	t_score int default 0 not null,
	outcome varchar(100) not null,
	winner varchar(10) not null,
	map_id int not null,
	start_time timestamp null,
	end_time timestamp null,
	constraint rounds_pk
		primary key (id)
);

create table if not exists hits
(
	id int auto_increment,
	round_id int not null,
	hitgroup_id int not null,
	weapon_id int not null,
	damage int default 0 not null,
	damage_armor int default 0 not null,
	attacker_id int not null,
	victim_id int null,
	hit_time timestamp not null,
	constraint hits_pk
		primary key (id),
	constraint hits_rounds__fk
		foreign key (round_id) references rounds (id)
			on update cascade on delete cascade,
	constraint hits_attacker__fk
		foreign key (attacker_id) references players (id)
			on update cascade on delete cascade,
	constraint hits_hitgroups_id_fk
		foreign key (hitgroup_id) references hitgroups (id)
			on update cascade on delete cascade,
	constraint hits_victim__fk
		foreign key (victim_id) references players (id)
			on update cascade on delete cascade,
	constraint hits_weapons_id_fk
		foreign key (weapon_id) references weapons (id)
			on update cascade on delete cascade
);

create table if not exists kills
(
	id int auto_increment,
	round_id int null,
	attacker_id int not null,
	victim_id int not null,
	weapon_id int not null,
	kill_time timestamp not null,
	constraint kills_pk
		primary key (id),
	constraint kills_attacker__fk
		foreign key (attacker_id) references players (id)
			on update cascade on delete cascade,
	constraint kills_rounds_id_fk
		foreign key (round_id) references rounds (id)
			on update cascade on delete cascade,
	constraint kills_victim__fk
		foreign key (victim_id) references players (id)
			on update cascade on delete cascade,
	constraint kills_weapons_id_fk
		foreign key (weapon_id) references weapons (id)
			on update cascade on delete cascade
);

create table if not exists player_round
(
	player_id int not null,
	round_id int not null,
	constraint round_player_rounds_id_fk
		foreign key (round_id) references rounds (id)
			on update cascade on delete cascade,
	constraint round_player_players_id_fk
		foreign key (player_id) references players (id)
			on update cascade on delete cascade
);
