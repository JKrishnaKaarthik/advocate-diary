package com.advocatediary.pages;

import java.util.List;
import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;

public class NewClientPage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='new-client-heading']");
  private final By fullName = By.name("full_name");
  private final By phone = By.name("phone");
  private final By email = By.name("email");
  private final By address = By.name("address");
  private final By notes = By.name("notes");
  private final By save = By.cssSelector("[data-testid='save-client']");
  private final By cancel = By.cssSelector("[data-testid='cancel-client'], a.button.secondary");
  private final By hints = By.cssSelector("p.info");

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
    return "Add a client".equals(heading())
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

  public String fieldValue(String fieldName) {
    return visible(By.name(fieldName)).getAttribute("value");
  }

  public boolean allFieldsEmpty() {
    return isBlank(fieldValue("full_name"))
        && isBlank(fieldValue("phone"))
        && isBlank(fieldValue("email"))
        && isBlank(fieldValue("address"))
        && isBlank(fieldValue("notes"));
  }

  public String validationMessage() {
    List<WebElement> messages = driver.findElements(hints);
    return messages.stream()
        .filter(WebElement::isDisplayed)
        .map(WebElement::getText)
        .findFirst()
        .orElse("");
  }

  public boolean showsValidation(String expected) {
    return driver.findElements(hints).stream()
        .anyMatch(el -> el.isDisplayed() && el.getText().contains(expected));
  }

  private boolean isBlank(String value) {
    return value == null || value.isBlank();
  }
}
