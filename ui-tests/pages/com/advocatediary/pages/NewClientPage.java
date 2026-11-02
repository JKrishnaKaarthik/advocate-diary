package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class NewClientPage extends BasePage {
  private final By heading = testId("new-client-heading");
  private final By fullName = testId("client-name");
  private final By phone = testId("client-phone");
  private final By email = testId("client-email");
  private final By address = testId("client-address");
  private final By notes = testId("client-notes");
  private final By save = testId("save-client");
  private final By cancel = testId("cancel-client");

  public NewClientPage(WebDriver driver) {
    super(driver);
  }

  public NewClientPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsAddClientDetails() {
    return shown(heading)
        && shown(fullName)
        && shown(phone)
        && shown(email)
        && shown(address)
        && shown(notes)
        && shown(save)
        && shown(cancel);
  }

  public void fill(String name, String phoneNumber, String emailAddress, String street, String note) {
    type(fullName, name);
    type(phone, phoneNumber);
    type(email, emailAddress);
    type(address, street);
    type(notes, note);
  }

  public void clickSaveClient() {
    pause(400);
    click(save);
  }

  public ClientsPage clickCancel() {
    pause(400);
    click(cancel);
    return new ClientsPage(driver).waitUntilLoaded();
  }

  public ClientPage saveClient(String name, String phoneNumber, String emailAddress, String street) {
    fill(name, phoneNumber, emailAddress, street, "");
    clickSaveClient();
    return new ClientPage(driver).waitUntilLoaded();
  }

  public ClientPage saveClientWithNotes(
      String name, String phoneNumber, String emailAddress, String street, String note) {
    fill(name, phoneNumber, emailAddress, street, note);
    clickSaveClient();
    return new ClientPage(driver).waitUntilLoaded();
  }

  public boolean allFieldsEmpty() {
    return isBlank(visible(fullName).getAttribute("value"))
        && isBlank(visible(phone).getAttribute("value"))
        && isBlank(visible(email).getAttribute("value"))
        && isBlank(visible(address).getAttribute("value"))
        && isBlank(visible(notes).getAttribute("value"));
  }

  public boolean showsNameRequired() {
    return hintShown("full_name");
  }

  public boolean showsEmailInvalid() {
    return hintShown("email");
  }

  private boolean isBlank(String value) {
    return value == null || value.isBlank();
  }
}
