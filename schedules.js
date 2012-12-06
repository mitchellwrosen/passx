function classToString(cls) {
   var str = "";
   if (cls["Lec"]) {
      str += cls["Subject"] + " " + cls["Number"] + "-" +
         cls["Lec"]["Section"] + " " + cls["Lec"]["Days"] + " " +
         cls["Lec"]["From"] + "-" + cls["Lec"]["To"] + "<br>";
   }

   if (cls["Lab"]) {
      str += cls["Subject"] + " " + cls["Number"] + "-" +
         cls["Lab"]["Section"] + " " + cls["Lab"]["Days"] + " " +
         cls["Lab"]["From"] + "-" + cls["Lab"]["To"] + "<br>";
   }

   console.log("class string: " + str);
   return str
}

function scheduleToString(schedule) {
   var str = "";
   for (var i = 0; i < schedule["Classes"].length; i++) {
      str += classToString(schedule["Classes"][i]);
   }
   str += "<br>"

   console.log("schedule string: " + str);
   return str
}

function init() {
   console.log(schedules);
   console.log(schedules.length);
   for (var i = 0; i < schedules.length; i++) {
      document.getElementById("schedules").innerHTML +=
         scheduleToString(schedules[i]);
   }
}

document.addEventListener("DOMContentLoaded", init)
