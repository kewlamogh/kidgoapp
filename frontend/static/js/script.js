let person;
let classID = localStorage.getItem('classID');
let sender = localStorage.getItem('sender') || null
let people = [
    "John",
    "Joe",
    "Jill",
    "Alice",
    "Bob",
    "Tanvi",
    "Max",
    "Inaaya",
    "Rahul",
    "Olivia",
]

if (localStorage.getItem("person")) {
    person = localStorage.getItem("person")
    document.getElementById("characterpick").style.display = "none"
    document.getElementById("greeting").innerText = "Hi " + person + "!"
    
    loadChats()
} else {
    document.getElementById("logout").style.display = "none"
    document.getElementById("help").style.display = "none"
    document.getElementById("greeting").innerText = "Click the green button with your name on it."

    // document.getElementById("classID").addEventListener("input", function () {
    //     let val = document.getElementById("classID").value

    //     if (val.includes(" ") || val.includes("/")) {
    //         alert("inavlid classroom id, cannot include spaces or /")
    //     } else {
    //         localStorage.setItem("classID", val)
    //     }
    // })
}

if (sender) {
    setChat(sender)
}

function passBlobUp(blob) {
    const formData = new FormData();
    formData.append('from', person);

    console.log(classID)

    formData.append('class', classID);
    formData.append('to', sender);
    formData.append('audioBlob', blob, 'audio.wav');

    fetch('sendmessage', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            console.log('Audio Blob uploaded successfully!');
        } else {
            console.error('Failed to upload Audio Blob:', response.status);
        }
    })
    .catch(error => console.error('Error uploading Audio Blob:', error));
}

function setPerson(x) {
    person = x
    localStorage.setItem("person", person)

    // if (!localStorage.getItem("classID")) {
    //     alert("Please enter your classroom ID!")
    //     return
    // }

    window.location.href = window.location.href
}

function loadChats() {
    let menu = document.getElementById("chatselect")
    let heading = document.createElement("h1")
    heading.innerText = "Select a friend!"

    menu.appendChild(heading)

    for (i of people) {
        if (i.toLowerCase() == person.toLowerCase()) {
            continue
        }

        let div = document.createElement("button")
        
        div.innerText = i
        div.onclick = () => {
            setChat(div.innerText)
        }

        menu.appendChild(div)
    }
}

async function getInbox() {
    if (person) {
        let resp = await fetch(`/inbox?user=${person}&class=${classID}`, {
            method: "GET",
        })

            
        let inbox = await resp.json()
        return inbox
    } else {
        alert("Sign in first!")
    }
}

async function getOutbox() {
    if (person) {
        let resp = await fetch(`/outbox?user=${person}&class=${classID}`, {
            method: "GET",
        })

    console.log("i hate this")
        let outbox = await resp.json()
        console.log(outbox)
        return outbox
    } else {
        alert("Sign in first!")
    }
}

function setChat(x) {
    sender = x
    localStorage.setItem("sender", sender)
    populateSendMenu()
}

async function playChat() {
    let inbox = await getInbox()
    let recentInbox = inbox.filter(item => item.sender == sender.toLowerCase())
    let outbox = await getOutbox()
    console.log(outbox)
    let recentOutbox = outbox.filter(item => item.reciever == sender.toLowerCase())

    let recent = recentInbox.concat(recentOutbox);
    recent.sort((a, b) => new Date(a.time) - new Date(b.time)) // Sort the array based on the time attribute
   
    // Sample list
    const numElements = Math.min(5, recent.length);
    const selectedElements = recent.slice(-numElements);

    console.log(selectedElements)
    playAudioFromList(selectedElements)
}

function refresh() {
    window.location.href = window.location.href
}

function populateSendMenu() {
    document.getElementById("sendmessageheading").innerText = "Send a message to " + sender

    let menu = document.getElementById("sendmenu")
    menu.innerHTML = "Your recording will be 10 seconds long<br/>"
    
    let recordButton = document.createElement("button")
    recordButton.innerText = "Start recording"
    recordButton.style.backgroundColor = "blue"
    recordButton.onclick = startRecording

    let refreshButton = document.createElement("button")
    refreshButton.innerText = "Refresh"
    refreshButton.style.backgroundColor = "yellow"
    refreshButton.onclick = () => {
        localStorage.setItem("sender", sender).href
        refresh()
    }

    let play = document.createElement("button")
    play.innerText = "Play conversation"
    play.style.backgroundColor = "yellow"
    play.style.color = "black"
    play.onclick = playChat

    menu.appendChild(recordButton)
    menu.appendChild(play)
}

let mediaRecorder;
let recordedChunks = [];


// Start recording audio for 10 seconds
function startRecording() {
    let menu = document.getElementById("sendmenu")
    let oldMenu = menu.innerHTML
    menu.innerHTML += "<br/>Recording will stop in 10 seconds..."

    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(stream => {
            mediaRecorder = new MediaRecorder(stream);
            console.log("recording started")

            mediaRecorder.ondataavailable = event => {
                recordedChunks.push(event.data);
            };

            mediaRecorder.onstop = () => {
                const audioBlob = new Blob(recordedChunks, { type: 'audio/wav' });
                // const audioUrl = URL.createObjectURL(audioBlob);
                // console.log('Recording complete. Audio URL:', audioUrl);
                menu.innerHTML = oldMenu

                passBlobUp(audioBlob)
            };

            mediaRecorder.start();

            setTimeout(() => {
                mediaRecorder.stop();
                refresh()
            }, 10000); // Stop recording after 10 seconds
        })
        .catch(error => {
            console.error('Error accessing microphone:', error);
        });
}



function playAudioFromList(listOfObjects, index = 0) {
    if (index >= listOfObjects.length) {
        return; // Exit if all objects have been processed
    }

    const obj = listOfObjects[index];
    const audioLink = obj.s3_resource_link;

    if (audioLink) {
        const audio = new Audio(audioLink);
        audio.onended = () => {
            // Play next audio when the current one ends
            playAudioFromList(listOfObjects, index + 1);
        };
        audio.play();
    }
}