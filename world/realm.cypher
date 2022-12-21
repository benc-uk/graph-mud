MATCH (n)
DETACH DELETE n;

// ------------------------------------------------------------------

CREATE (lobby:Location {name: 'lobby', description: 'The game lobby, a magical space that resides outside the game realm. Use the $play command to enter the realm and start playing!'})
RETURN lobby;

CREATE (l1:Location {gameEntry: true, name: 'courtyard', description: 'A courtyard with a moss covered fountain in the middle, it is surrounded by a high wall with a gate at the far end. Ivy grows up the walls and over the gate.'})
CREATE (l2:Location {name: 'hallway', description: 'A long hallway with a staircase at the end'})
CREATE (l3:Location {name: 'kitchens', description: 'A large kitchen with a large oven'})
CREATE (l4:Location {name: 'training', description: 'Training grounds with a large wooden dummy and an archery range'})
CREATE (l5:Location {name: 'gates', description: 'Castle gates, there is a small door leading into the gatehouse, and a huge open gate leading beyond the castle walls.'})
CREATE (l6:Location {name: 'gatehouse', description: 'A dusty gatehouse, there are steep stairs leading up into the castle walls. There are tables and chairs in the corner.'})
CREATE (l7:Location {name: 'outside-castle-gates', description: 'Some grass and the rest of the world, which is inaccessable, until the game is finished.'})
CREATE (l8:Location {name: 'castle-walls', description: 'The castle walls. It\'s nice up here.'})

MERGE (l1)-[:NORTH]->(l2)
MERGE (l1)-[:SOUTH]->(l5)
MERGE (l2)-[:SOUTH]->(l1)
MERGE (l2)-[:EAST]->(l3)
MERGE (l3)-[:WEST]->(l2)

MERGE (l1)-[:WEST]->(l4)
MERGE (l4)-[:EAST]->(l1)

MERGE (l5)-[:NORTH]->(l1)
MERGE (l5)-[:EAST]->(l6)
MERGE (l5)-[:SOUTH]->(l7)

MERGE (l6)-[:WEST]->(l5)
MERGE (l6)-[:UP]->(l8)

MERGE (l7)-[:NORTH]->(l5)

MERGE (l8)-[:DOWN]->(l6)

CREATE (i1:Item {name: 'key', description: 'A small key'})
CREATE (i2:Item {name: 'sword', description: 'A long sword'})
CREATE (i3:Item {name: 'potion', description: 'A bubbling potion'})
CREATE (i4:Item {name: '50 gold', description: '50 gold coins', type: 'gold', value: 50})

CREATE (l1)-[:STORES]->(i1)
CREATE (l1)-[:STORES]->(i2)
CREATE (l3)-[:STORES]->(i3)
CREATE (l6)-[:STORES]->(i4);
