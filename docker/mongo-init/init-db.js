// Création des bases si non présentes
const createIfNotExists = (dbName) => {
  const currentDB = db.getSiblingDB(dbName);

  if (!currentDB.getCollectionNames().includes("init")) {
    currentDB.createCollection("init");
    currentDB.init.insertOne({
      createdAt: new Date(),
      msg: `Base ${dbName} initialisée ✅`
    });
    print(`✅ La base ${dbName} a été créée et initialisée.`);
  } else {
    print(`ℹ️ La base ${dbName} existe déjà, aucune action effectuée.`);
  }
};

createIfNotExists("todof");
createIfNotExists("todof_test");
