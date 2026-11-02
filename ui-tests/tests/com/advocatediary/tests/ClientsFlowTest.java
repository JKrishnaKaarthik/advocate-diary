package com.advocatediary.tests;

import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.ClientsPage;
import com.advocatediary.pages.NewClientPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class ClientsFlowTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void clickingClientsOpensClientsPage() {
    ClientsPage clients = signInToDiary().nav().openClients();
    Assert.assertTrue(clients.isClientsPage());
    Assert.assertEquals(clients.heading(), "Clients");
  }

  @Test
  public void clickingAddClientShowsAddClientDetails() {
    NewClientPage add = signInToDiary().nav().openClients().addClient();
    Assert.assertEquals(add.heading(), "Add a client");
    Assert.assertTrue(add.showsAddClientDetails(), "name, phone, email, address, notes, save and cancel should show");
  }

  @Test
  public void saveWithoutNameShowsNameRequired() {
    NewClientPage add = signInToDiary().nav().openClients().addClient();
    add.clickSaveClient();
    Assert.assertTrue(add.showsValidation("Name is required"));
    Assert.assertEquals(add.heading(), "Add a client");
  }

  @Test
  public void invalidEmailShowsAppropriateMessage() {
    NewClientPage add = signInToDiary().nav().openClients().addClient();
    add.fill("Meera Rao", "9876543210", "not-an-email", "Chennai", "Family matter");
    add.clickSaveClient();
    Assert.assertTrue(add.showsValidation("Enter a valid email"));
    Assert.assertEquals(add.heading(), "Add a client");
  }

  @Test
  public void cancelRemovesFilledInformation() {
    NewClientPage add = signInToDiary().nav().openClients().addClient();
    add.fill("Meera Rao", "9876543210", "meera@example.test", "Chennai", "Family matter");
    ClientsPage clients = add.clickCancel();
    Assert.assertTrue(clients.isClientsPage());
    add = clients.addClient();
    Assert.assertTrue(add.allFieldsEmpty(), "cancel should refresh so Add a client is empty");
  }

  @Test
  public void saveClientOpensClientPage() {
    NewClientPage add = signInToDiary().nav().openClients().addClient();
    ClientPage client = add.saveClientWithNotes(
        "Meera Rao",
        "9876543210",
        "meera@example.test",
        "Chennai",
        "Family matter");
    Assert.assertEquals(client.heading(), "Meera Rao");
  }
}
