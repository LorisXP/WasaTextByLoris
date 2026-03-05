import axios from "./services/axios.js";
import auth from "./services/auth.js";


/**
 * Effettua il login
 * @returns {Promise<string>} Un messaggio di errore, o una stringa vuota se il login è riuscito
 */
async function doLogin(userName) {
    let errormsg = "Errore durante il login";
    try {
        if (userName || typeof userName == "string") {
            console.log("userName is valid");

            //Effettua la richiesta di login al server
            const response = await axios.post("/api/auth", {
                userName: userName,
            });
            console.debug("Login response:", response);

            //Se la risposta è positiva, ottieni l'userID e salvalo nell'autenticazione
            if (response.status === 200 || response.status === 201) {
                console.log("Login successful, processing response");

                //Ottieni l'userID dalla risposta e salvalo nell'autenticazione
                let userID = response.data.userID;
                if (userID > 0) {
                    auth.setUserID(userID);
                    // Resetta il messaggio di errore
                    errormsg = "";
                    console.log("Authenticated, userID set", userID);
                } else {
                    console.warn("Login succeeded but no userID returned");
                }
            } else {
                errormsg = "Impossibile accedere";
            }
                
        } else {
            errormsg = "Inserisci uno username valido";
        }

    } catch (e) {
        errormsg = e.toString();
        console.error("Login error:", e);
    } 

    return errormsg;
}