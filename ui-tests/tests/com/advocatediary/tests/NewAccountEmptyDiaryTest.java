package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.DiaryPage;
import org.testng.Assert;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.Test;

public class NewAccountEmptyDiaryTest extends BaseTest {

  @Override
  @BeforeMethod
  public void openBrowser() {
    System.setProperty("baseUrl", "http://localhost:8080");
    System.setProperty("headless", "false");
    super.openBrowser();
  }

  @Test
  public void newAccountEmptyDiaryThenSaveCaseDetails() {
    DiaryPage diary = signUpNewAccount();
    Assert.assertTrue(diary.isOpen());
    diary.waitUntilEmpty();
    Assert.assertTrue(diary.isEmpty(), "new account diary must be empty");

    ClientPage client = addClient(CLIENT_NAME);
    CasePage matter = openAndSaveCase(client);
    Assert.assertEquals(matter.heading(), CASE_TITLE);
    Assert.assertEquals(matter.titleValue(), CASE_TITLE);
  }
}
