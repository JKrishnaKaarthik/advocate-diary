package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class ClientCaseTest extends BaseTest {

  @Test
  public void lawyerCanAddClientAndOpenCase() {
    Assert.assertTrue(signInToDiary().isOpen());
    CasePage matter = openCasePage();
    Assert.assertEquals(matter.heading(), CASE_OPEN_TITLE);
    Assert.assertTrue(matter.nav().openInbox().isOpen());
  }
}
