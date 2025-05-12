create table if not exists t_coin
(
    pk_serial TEXT,
    symbol    TEXT,
    num       REAL
);


DELETE
FROM t_coin
where pk_serial = '6f9d25e43a6b43e2bf500f5d4c7f7a63';

insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'BTC', 0);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'ADA', 23671);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'BNB', 3.1467);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'WIF', 2322);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'XRP', 595.88);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'SUI', 315.82);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'DOGE', 3263);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'ETH', 0.2175);
insert into t_coin(pk_serial, symbol, num)
values ('6f9d25e43a6b43e2bf500f5d4c7f7a63', 'USDC', 266);
