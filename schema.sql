create table unfinished_commands (
	id serial primary key,
	command varchar(32) not null,
	chat_id int not null,
	created_by int not null,
	created_at timestamp default current_timestamp not null,
	unique (chat_id, created_by)
);

create table shopping_items (
	id serial primary key,
	name varchar(255) not null,
	chat_id int not null,
	created_by int not null,
	created_at timestamp default current_timestamp not null
);
