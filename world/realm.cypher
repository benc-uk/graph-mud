MATCH (n)
DETACH DELETE n;

// ------------------------------------------------------------------

CREATE (lobby:Location {name: 'lobby', description: 'This is the game lobby, it exists outside the game realm. You can chat to other players but can not engage in combat etc. Use the $play command to enter the relams and start playing!'})
RETURN lobby;

CREATE (l1:Location {name: 'courtyard', description: 'A courtyard with a moss covered fountain in the middle, it is surrounded by a high wall with a gate at the far end. Ivy grows up the walls and over the gate.'})
CREATE (l2:Location {name: 'hallway', description: 'A long hallway with a staircase at the end'})
CREATE (l3:Location {name: 'kitchens', description: 'A large kitchen with a large oven'})
MERGE (l1)-[:NORTH]->(l2)
MERGE (l2)-[:SOUTH]->(l1)
MERGE (l2)-[:EAST]->(l3)
MERGE (l3)-[:WEST]->(l2)

CREATE (i1:Item {name: 'key', description: 'A small key'})
CREATE (i2:Item {name: 'sword', description: 'A long sword'})
CREATE (i3:Item {name: 'potion', description: 'A bubbling potion'})

CREATE (l1)-[:HOLDS]->(i1)
CREATE (l1)-[:HOLDS]->(i2)
CREATE (l3)-[:HOLDS]->(i3);
