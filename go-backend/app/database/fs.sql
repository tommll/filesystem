drop sequence if exists items_id_seq;
CREATE SEQUENCE items_id_seq;

drop table if exists items;
create table items(
    id int not null primary key,
    name varchar(64) unique,
    item_type int,
    order_num int,
    parent_id int,
    data varchar(200),

    foreign key (parent_id) references items(id)
);

CREATE INDEX idx_item_order ON items (order_num);

-- insert root node
INSERT INTO items (id, item_type, order_num, name) VALUES (0, 0, 0, '~')
