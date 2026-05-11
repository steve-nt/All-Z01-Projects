function citiesOnly(array){
    return array.map(item => item.city);
}

function upperCasingStates(array){
    return array.map(str => str.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' '));
}

function fahrenheitToCelsius(array){
    return array.map(temp => {
        const f = parseInt(temp); // Extract number from '68°F'
        const c = Math.floor((f - 32) * 5 / 9); // Convert and round down
        return `${c}°C`;
    });
}

function trimTemp(array) {
    return array.map(obj => ({...obj, temperature: obj.temperature.replaceAll(" ", "")}))
}


function tempForecasts(array) {
    return trimTemp(array).map(({ city, state, temperature }) => {
        const tempC = fahrenheitToCelsius([temperature])[0];
        const stateFormatted = upperCasingStates([state])[0];
        return `${tempC}elsius in ${city}, ${stateFormatted}`;
    });
}