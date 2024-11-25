CREATE TABLE library
(
    id           SERIAL PRIMARY KEY,
    group_name   VARCHAR NOT NULL,
    song         VARCHAR NOT NULL,
    release_date TIMESTAMP NULL,
    text         TEXT NULL,
    link         TEXT NULL
);

INSERT INTO library (group_name, song, release_date, text, link)
VALUES ('The Cosmic Waves', 'Stars Above', '2023-10-05 00:00:00',
        E'test\n\ntest2\n\ntest3\n\ntestverses4',
        'https://example.com/stars_above'),

       ('Echoes of Silence', 'Whispers in the Dark', '2022-05-15 00:00:00',
        E'There are places I remember\nAll my life, though some have changed\nSome forever, not for better\nSome have gone and some remain\n\nAll these places had their moments\nWith lovers and friends, I still can recall\nSome are dead and some are living\nIn my life, I''ve loved them all\n\nBut of all these friends and lovers\nThere is no one compares with you\nAnd these memories lose their meaning\nWhen I think of love as something new\n\nThough I know I''ll never lose affection\nFor people and things that went before\nI know I''ll often stop and think about them\nIn my life, I love you more\n\nSee upcoming rock shows\nGet tickets for your favorite artists\nYou might also like\neuphoria\nKendrick Lamar\nNow and Then\nThe Beatles\nBut Daddy I Love Him\n\nThough I know I''ll never lose affection\nFor people and things that went before\nI know I''ll often stop and think about them\nIn my life, I love you more\n\nIn my life, I love you more',
        'https://example.com/whispers_in_the_dark'),

       ('Neon Dreamers', 'Electric Hearts', '2021-08-22 00:00:00',
        E'Underneath the neon glow,\nWe move together, fast and slow,\nElectric hearts, we feel the spark,\nWe are alive, we leave our mark.\n\nCity lights, they guide the way,\nIn the night, we find our stay,\nThrough the dark, our love will fly,\nLike neon dreams that never die.',
        'https://example.com/electric_hearts'),

       ('Crimson Horizon', 'Burning Sky', '2020-12-01 00:00:00',
        E'The sky is burning, red and wide,\nWe feel the fire deep inside,\nThe storm is raging, wild and free,\nWe rise above the endless sea.\n\nOnward we go, to the unknown,\nWe leave behind what we have known,\nThe fire will guide us to the light,\nThrough the darkness, we’ll take flight.',
        'https://example.com/burning_sky'),

       ('Silver Lining', 'Chasing the Sun', '2019-11-10 00:00:00',
        E'We’re chasing the sun, through the rain,\nA dream so bright, we break the chain,\nThrough the clouds, we’ll find our way,\nTogether we’ll rise, come what may.\n\nThe horizon calls, we run and go,\nWe won’t stop, we won’t slow,\nThe sun awaits, it’s in our sight,\nWe chase the dawn, we own the night.',
        'https://example.com/chasing_the_sun'),

       ('The Lost Horizon', 'Echoes of Tomorrow', '2018-07-19 00:00:00',
        E'We hear the echoes from afar,\nOf tomorrow’s dreams, like distant stars,\nThey whisper secrets in the breeze,\nOf all the things we cannot see.\n\nThe past is gone, but we still hear,\nThe voices calling, drawing near,\nThrough time and space, we’ll find our way,\nTo echoes that will never fade away.',
        'https://example.com/echoes_of_tomorrow');