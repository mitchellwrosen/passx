var classes = [];

var selectedClass;

function updateClasses() {
   console.log($("#classes").children().length + " classes");
}

function onClassRbClicked(classNum) {
   // Hide old class sections.
   $(".class" + selectedClass + "-sec").css("display", "none");

   // Show new class sections.
   $(".class" + classNum + "-sec").css("display", "inline-block");

   selectedClass = classNum;
}

function onClassLiClicked() {
   console.log("clicked");
}

function onAddClass() {
   $.get('static/templates.html', function(templates) {
      classes.push({});

      var template = $(templates).filter('#new-class').html();
      $("#classes").append(
         Mustache.render(template, { "classNum": classes.length - 1 }));
   });
}

function onAddLecture() {
   $.get('static/templates.html', function(templates) {
      var template = $(templates).filter('#new-lec').html();
      $("#sections").append(
         Mustache.render(template, { "classNum": selectedClass }));
   });
}
