"use strict";

function ensureJQueryLoadedAsync() {
  return new Promise((resolve) => {
    if (typeof jQuery !== "undefined") {
      resolve();
      return;
    }
    if (window._loadingJQuery) {
      // Already loading, poll until available
      const poll = setInterval(function() {
        if (typeof jQuery !== "undefined") {
          clearInterval(poll);
          resolve();
        }
      }, 50);
      return;
    }
    window._loadingJQuery = true;
    var script = document.createElement("script");
    script.src = "https://code.jquery.com/jquery-3.7.1.min.js";
    script.type = "text/javascript";
    script.onload = function() {
      window._loadingJQuery = false;
      resolve();
    };
    document.head.appendChild(script);
  });
}

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
      // Dish name: first child node with a data-cy attribute that is not a wrapper or ingredients
      var dishName = $(this).find('[data-cy]').filter(function() {
        var val = $(this).attr('data-cy');
        return val !== "dish-tile__wrapper" && val !== "IngredientsAndRecipes_span";
      }).first().text().trim();

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

(async function() {
  await ensureJQueryLoadedAsync();
  const meals = getMealsAndIngredients();
  const filename = getBestFilename();
  saveToFile(meals, filename);
})();
