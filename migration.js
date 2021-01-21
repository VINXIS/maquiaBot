let fs = require("fs");

let reminders = require("./data/reminders.json");
let profiles = require("./data/osuData/profileCache.json");

// change profiles
for (let i = 0; i < profiles.length; i++) {
    profiles[i].Discord = profiles[i].Discord.id;
    console.log("Updated user # " + (i + 1));
}
fs.writeFile("./data/osuData/profileCache.json", JSON.stringify(profiles), (err) => {
    if (err) throw err;
    console.log("Updated player cache!");
});

// change reminders
for (let i = 0; i < reminders.length; i++) {
    reminders[i].User = reminders[i].User.id;
    console.log("Updated reminder # " + (i + 1));
}
fs.writeFile("./data/reminders.json", JSON.stringify(reminders), (err) => {
    if (err) throw err;
    console.log("Updated reminders!");
});

// change serverData
fs.readdir("./data/serverData", (err, files) => {
    if (err) throw err;

    for (const file of files) {
        if (!/\d+\.json/.test(file))
            continue;
        let serverData = require(`./data/serverData/${file}`);
        // change role automations
        if (serverData.RoleAutomation)
            for (let i = 0; i < serverData.RoleAutomation.length; i++) {
                for (let j = 0; j < serverData.RoleAutomation[i].Roles.length; j++) {
                    serverData.RoleAutomation[i].Roles[j] = serverData.RoleAutomation[i].Roles[j].id;
                }
                console.log("Updated role automation # " + (i + 1));
            }

        if (serverData.Counters)
            for (let i = 0; i < serverData.Counters.length; i++) {
                for (let j = 0; j < serverData.Counters[i].Users.length; j++) {
                    serverData.Counters[i].Users[j].Username = serverData.Counters[i].Users[j].User.username;
                    serverData.Counters[i].Users[j].UserID = serverData.Counters[i].Users[j].User.id;
                    serverData.Counters[i].Users[j].User = undefined;
                }
                console.log("Updated counter # " + (i + 1));
            }
        
        fs.writeFile(`./data/serverData/${file}`, JSON.stringify(serverData), (err) => {
            if (err) throw err;
            console.log("Updated server ID " + (file));
        });
    }
});

// change serverData
fs.readdir("./data/channelData", (err, files) => {
    if (err) throw err;

    for (const file of files) {
        if (!/\d+\.json/.test(file))
            continue;
        let channelData = require(`./data/channelData/${file}`);
        // change role automations
        if (channelData.Tracking) {
            channelData.OsuTracking = channelData.Tracking
            channelData.Tracking = undefined;
            fs.writeFile(`./data/channelData/${file}`, JSON.stringify(channelData), (err) => {
                if (err) throw err;
                console.log("Updated channel ID " + (file));
            });
        }
    }
});

console.log("Ok done");