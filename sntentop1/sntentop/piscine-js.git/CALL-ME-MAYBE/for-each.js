/*
 * CUSTOM FOREACH IMPLEMENTATION - Executes a callback for each array element.
 * 
 * Concepts:
 * - forEach() iterates through each element and executes a callback
 * - Callback receives (element, index, array) as parameters
 * - Optionally accepts a context (this) as third argument
 * - Returns undefined (used for side effects, not transformations)
 * 
 * Real-world examples:
 * - Rendering list items to the DOM
 * - Sending analytics events for each user action
 * - Updating database records one by one
 * - Logging validation errors for each form field
 * 
 * Bonus challenges:
 * - Implement forEach with support for custom context binding
 * - Handle early exit using try/catch with custom exception
 * - Compare performance with for loop vs forEach
 */

const forEach = (arr, callback, context) => {
  for (let i = 0; i < arr.length; i++) {
    callback.call(context, arr[i], i, arr);
  }
};