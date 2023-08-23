create table authors (
	id BIGSERIAL PRIMARY KEY,
	name text not null,
	age int not null default 0,
	bio text
);
