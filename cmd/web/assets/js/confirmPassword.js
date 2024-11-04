function confirmPassword() {
    const password = document.getElementById("password").value;
    const confirm = document.getElementById("confirm_password").value;
    const errorMessage = document.getElementById("password-error");

    if (password != confirm) {
        errorMessage.classList.remove("hidden"); // Show the error message
        document.getElementById("confirm_password").focus(); // Focus the confirm password field
        return false
    }
    errorMessage.classList.add("hidden"); // Hide the error message
    return true
}