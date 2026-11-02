package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class NewAccountEmptyDiaryTest extends BaseTest {

  @Test
  public void newAccountEmptyDiaryThenSaveCaseDetails() {
    Assert.assertTrue(signupEmptyDiary().isEmpty(), "new account diary must be empty");
    CasePage matter = openAndSaveCase(addClient(CLIENT_NAME));
    Assert.assertEquals(matter.heading(), CASE_TITLE);
    Assert.assertEquals(matter.titleValue(), CASE_TITLE);
  }
}
