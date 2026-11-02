package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class OpenCaseFlowTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void afterSaveClientOpenCaseThenSaveCaseDetails() {
    signInToDiary();
    ClientPage client = addClient(CLIENT_NAME);
    Assert.assertEquals(client.heading(), CLIENT_NAME);
    Assert.assertTrue(client.showsOpenACase(), "Open a case form should show after save client");

    CasePage matter = openAndSaveCase(client);
    Assert.assertEquals(matter.heading(), CASE_TITLE);
    Assert.assertEquals(matter.titleValue(), CASE_TITLE);
    Assert.assertEquals(matter.notesValue(), "Listed for admission");
  }
}
