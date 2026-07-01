"use strict";

const form = document.getElementById("dank-auth-form");
const fetchButton = document.getElementById("fetch-requests-btn");

function getCookie(key) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${key}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
}

function checkAuth() {
  const cookieToken = getCookie("admin-token");
  if (cookieToken) {
    form.style.display = "none";
    fetchButton.style.display = "block";
  }
}

function init() {
  checkAuth();
  console.log("helloooo 👋");
}

init();
