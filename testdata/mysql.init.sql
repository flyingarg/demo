create database demo
truncate demo.user_details;
truncate demo.tokens;
truncate demo.requests;
create table user_details(
                             id mediumint not null auto_increment,
                             username varchar(50),
                             plan varchar(10),
                             PRIMARY KEY (id),
                             UNIQUE KEY 'username' (username)
);

create table tokens(
                       user_id int not null,
                       user_secret varchar(64),
                       last_updated timestamp,
                       PRIMARY KEY (user_id)
);
create table requests(
                         id varchar(100),
                         user_id int,
                         data_format varchar(4),
                         data blob,
                         status bool,
                         request_time timestamp,
                         PRIMARY KEY (id)
);
insert into user_details(username, plan) values('mohit', 'BRONZE');
insert into tokens(user_id, user_secret, last_updated) values(1, 'asdas234-9nhj324brf9f834', now());
