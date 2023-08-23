-- name: GetAuthor :one
select * from authors
where id = $1 LIMIT 1;

-- name: ListAuthors :many
select * from authors
order by name;

-- name: CreateAuthor :one
insert into authors (
	name, bio, age
) values (
	$1, $2, $3
)
returning *;

-- name: UpdateAuthor :one
update authors
	set name = $2,
	bio = $3
where id = $1
returning *;

-- name: DeleteAuthor :exec
delete from authors
where id = $1;
