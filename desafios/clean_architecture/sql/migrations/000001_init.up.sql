create table if not exists orders (
        id          varchar(255)    primary key not null
    ,   price       decimal(10,2)   not null
    ,   tax         decimal(10,2)   not null
    ,   final_price decimal(10,2)   not null
);