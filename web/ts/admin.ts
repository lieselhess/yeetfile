import {Endpoints} from "./endpoints.js";
import {AdminFileInfoResponse, AdminUserInfoResponse} from "./interfaces.js";

const init = () => {
    setupUserSearch();
    setupFileSearch();
}

// =============================================================================
// User admin
// =============================================================================

const setupUserSearch = () => {
    let userSearchBtn = document.getElementById("user-search-btn") as HTMLButtonElement;
    let userIDInput = document.getElementById("user-id") as HTMLInputElement;

    userSearchBtn.addEventListener("click", () => {
        let userID = userIDInput.value;
        fetch(Endpoints.format(Endpoints.AdminUserActions, userID)).then(async response => {
            if (!response.ok) {
                alert("Error fetching user: " + await response.text());
                return;
            }

            let responseDiv = document.getElementById("user-response");
            responseDiv.innerHTML = "";

            let userInfo = new AdminUserInfoResponse(await response.json());
            let userDiv = generateUserActionsHTML(userInfo);
            responseDiv.appendChild(userDiv);
        }).catch((error: Error) => {
            alert("Error fetching user");
            console.error(error);
        })
    });
}

const generateUserActionsHTML = (userInfo: AdminUserInfoResponse): HTMLDivElement => {
    let userResponseDiv = document.createElement("div") as HTMLDivElement;

    userResponseDiv.className = "bordered-box visible";

    let userInfoElement = document.createElement("code");
    userInfoElement.innerText = `ID: ${userInfo.id}
Email: ${userInfo.email}
Storage Used: ${userInfo.storageUsed}
Send Used: ${userInfo.sendUsed}`;

    userResponseDiv.appendChild(userInfoElement);
    userResponseDiv.appendChild(document.createElement("br"));

    let deleteBtnID = `delete-user-${userInfo.id}`;
    let deleteButton = document.createElement("button");
    deleteButton.id = deleteBtnID;
    deleteButton.className = "red-button";
    deleteButton.innerText = "Delete User and Uploads";

    userResponseDiv.appendChild(deleteButton);

    deleteButton.addEventListener("click", () => {
        if (!confirm("Deleting this user will also delete all files they have " +
            "uploaded. Do you wish to proceed?")) {
            return;
        }

        fetch(Endpoints.format(Endpoints.AdminUserActions, userInfo.id), {
            method: "DELETE"
        }).then(async response => {
            if (!response.ok) {
                alert("Failed to delete user! " + await response.text());
            } else {
                alert("User and their content has been deleted!");
                userResponseDiv.innerHTML = "";
                userResponseDiv.className = "hidden";
            }
        }).catch(error => {
            alert("Failed to delete user");
            console.error(error);
        });
    });

    if (userInfo.files.length > 0) {
        let header = document.createElement("span");
        header.className = "span-header";
        header.innerText = "Files:"
        userResponseDiv.appendChild(header);
    }

    for (let i = 0; i < userInfo.files.length; i++) {
        let fileDiv = generateFileActionsHTML(userInfo.files[i]);
        userResponseDiv.appendChild(fileDiv);
    }

    return userResponseDiv;
}

// =============================================================================
// File admin
// =============================================================================

const setupFileSearch = () => {
    let fileSearchBtn = document.getElementById("file-search-btn") as HTMLButtonElement;
    let fileIDInput = document.getElementById("file-id") as HTMLInputElement;

    fileSearchBtn.addEventListener("click", () => {
        let fileID = fileIDInput.value;
        fetch(Endpoints.format(Endpoints.AdminFileActions, fileID)).then(async response => {
            if (!response.ok) {
                alert("Error fetching file: " + await response.text());
                return
            }

            let responseDiv = document.getElementById("file-response");
            responseDiv.innerHTML = "";

            let fileInfo = new AdminFileInfoResponse(await response.json());
            let fileDiv = generateFileActionsHTML(fileInfo);
            responseDiv.appendChild(fileDiv);
        }).catch(error => {
            console.log(error);
        })
    });
}

const generateFileActionsHTML = (fileInfo: AdminFileInfoResponse): HTMLDivElement => {
    let fileResponseDiv = document.createElement("div") as HTMLDivElement;

    fileResponseDiv.className = "bordered-box visible";

    let fileInfoElement = document.createElement("code");
    fileInfoElement.innerText = `ID: ${fileInfo.id}
Stored Name (encrypted): ${fileInfo.bucketName}
Size: ${fileInfo.size}
Owner ID: ${fileInfo.ownerID}
Modified: ${fileInfo.modified}`;

    fileResponseDiv.appendChild(fileInfoElement);
    fileResponseDiv.appendChild(document.createElement("br"));

    let deleteBtnID = `delete-file-${fileInfo.id}`;
    let deleteButton = document.createElement("button");
    deleteButton.id = deleteBtnID;
    deleteButton.className = "red-button";
    deleteButton.innerText = "Delete File";

    fileResponseDiv.appendChild(deleteButton);

    deleteButton.addEventListener("click", () => {
        if (!confirm("Deleting this file is irreversible. Proceed?")) {
            return;
        }

        fetch(Endpoints.format(Endpoints.AdminFileActions, fileInfo.id), {
            method: "DELETE"
        }).then(async response => {
            if (!response.ok) {
                alert("Failed to delete file! " + await response.text());
            } else {
                alert("The file has been deleted!");
                fileResponseDiv.innerHTML = "";
                fileResponseDiv.className = "hidden";
            }
        }).catch(error => {
            alert("Failed to delete file");
            console.error(error);
        });
    });

    return fileResponseDiv;
}

if (document.readyState !== "loading") {
    init();
} else {
    document.addEventListener("DOMContentLoaded", () => {
        init();
    });
}