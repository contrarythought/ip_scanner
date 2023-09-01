create table ip_addresses (
    id serial primary key,
    ip_addr text unique not null,
    ports_avail text[],
    geo_loc text
);