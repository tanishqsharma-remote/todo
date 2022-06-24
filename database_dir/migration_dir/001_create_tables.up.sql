create table if not exists users(id serial primary key , username varchar, password varchar);
create table if not exists todolist(user_id int , task varchar primary key , completed bool, archived bool, foreign key (user_id) references users(id));
