/* Snippet to load and run this script from browser console:
------------------------------------------------------------------------------
fetch('https://cdn.jsdelivr.net/gh/toszr/dietician@v0.4.0/js/extract-meals.js')
  .then(response => response.text())
  .then(text => eval(text));
--------------------------------- or, a bookmarklet: -------------------------
javascript:(function(){fetch('https://cdn.jsdelivr.net/gh/toszr/dietician@v0.4.0/js/extract-meals.js').then(r=>r.text()).then(t=>eval(t))})();
------------------------------------------------------------------------------
*/

"use strict";

function run($) {
  function getMealsAndIngredients() {
    let meals = [];

    $('[data-cy="MealDropdownOptions_div"]').each(function() {
      // Meal type: first non-empty text node inside the meal node
      var mealType = $(this).find('*').addBack().contents().filter(function() {
        return this.nodeType === 3 && $.trim(this.nodeValue) !== "";
      }).first().text().trim();

      if (!mealType) return;

      let dishes = [];
      $(this).find('[data-cy="dish-tile__wrapper"]').each(function() {
        // Dish name: all text from nodes with data-cy="MenuDishName_div"
        var dishName = $(this).find('[data-cy="MenuDishName_div"]').text().trim();

        // Ingredients: all text from nodes with data-cy="IngredientsAndRecipes_span"
        var ingredients = $(this).find('[data-cy="IngredientsAndRecipes_span"]').text().trim();

        dishes.push({
          dishName: dishName,
          ingredientsList: ingredients
        });
      });

      meals.push({
        mealName: mealType,
        dishes: dishes
      });
    });
    return meals;
  }


  function saveToFile(data, filename) {
    const json = JSON.stringify(data, null, 2);
    const blob = new Blob([json], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  function getBestFilename() {
    let filename = 'meals.json';
    const dateNode = $('[data-cy="DateItemDetails_div"]');
    if (dateNode.length) {
      const dateText = dateNode.text(); // Get all text inside the node
      const dateMatch = dateText.match(/(\d{2})-(\d{2})-(20\d{2})/); // Find a date with a 20xx year
      if (dateMatch) {
        const day = dateMatch[1];
        const month = dateMatch[2];
        const year = dateMatch[3];
        filename = `${day}${month}${year.slice(-2)}.json`;
      }
    }
    return filename;
  }

  const meals = getMealsAndIngredients();
  const filename = getBestFilename();
  saveToFile(meals, filename);
}

async function ensureJQueryLoadedAsync() {
  return new Promise((resolve) => {
    if (typeof jQuery !== "undefined") {
      resolve(jQuery.noConflict(true));
      return;
    }
    if (window._loadingJQuery) {
      // Already loading, poll until available
      const poll = setInterval(function() {
        if (typeof jQuery !== "undefined") {
          clearInterval(poll);
          resolve(jQuery.noConflict(true));
        }
      }, 50);
      return;
    }
    window._loadingJQuery = true;
    var script = document.createElement("script");
    script.src = "https://code.jquery.com/jquery-3.7.1.min.js";
    script.type = "text/javascript";
    script.onload = function() {
      // Poll until jQuery is available, as onload doesn't guarantee it's ready.
      const poll = setInterval(function() {
        if (typeof jQuery !== "undefined") {
          clearInterval(poll);
          window._loadingJQuery = false;
          resolve(jQuery.noConflict(true));
        }
      }, 50);
    };
    document.head.appendChild(script);
  });
}

(async () => {
  const jq = await ensureJQueryLoadedAsync();
  run(jq);
})();
