// MATCH (locations:Location)
// RETURN locations;

// MATCH (item:Item)--(t)
// RETURN item.name, t.name;

// MATCH (player:Player)--(loc:Location)
// RETURN player.username, loc.name;

MATCH (p:Player)
RETURN p;
