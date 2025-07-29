function showDialog(id) {
    document.getElementById(id)?.showModal();
}
function closeDialog(id) {
    document.getElementById(id)?.close();
}
function displayLocations(id) {
    fetch('locations/' + (Number(id) + 1))
        .then(res => res.json())
        .then(data => {
            const content = document.getElementById("LocationsPlace" + id);
            console.log(content, id);
            if (data.locations.length > 0) {
                content.textContent = data.locations.join(', ');
            }
        })
        .catch(err => console.error(err));
}

function displayMore(id) {
    // relations

    fetch("/relation/" + (Number(id) + 1))
        .then(response => response.json())
        .then(data => {
            const content = document.getElementById("RelationPlace" + id);

            const lines = [];

            for (const location in data.datesLocations) {
                const dates = data.datesLocations[location].join(", ");
                lines.push(location + ": " + dates);
            }

            content.textContent = lines.join("\n");
        })
        .catch(error => console.error(error));

    // dates
    fetch("/dates/" + (Number(id) + 1))
        .then(res => res.json())
        .then(data => {
            const content = document.getElementById("DatesPlace" + id);
            if (data.dates && data.dates.length) {
                content.textContent = data.dates.join(", ");
            }
        })
        .catch(error => console.error(error));
}


